package v1

import (
	"fmt"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/rabbitmq"
	"github.com/EDDYCJY/go-gin-example/service/rabbitmq_service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func AddMqUsers(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form rabbitmq_service.AddRabbitMQUserForm
	)

	httpCode, errCode := app.BindJsonAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	user := rabbitmq_service.ConvertAddFormToUser(form)
	userID, err := user.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ORDER_FAIL, nil)
		return
	}

	//添加成功之后给 RabbitMQ 发送注册成功的消息
	message := fmt.Sprintf(`{"user_id": %d}`, userID)
	pushErr := rabbitmq.PublishMessage(
		"user_register",   // 交换机名
		"direct",          // 类型
		"user_register_q", // 队列名
		"user.register",   // 路由键
		message,           // 消息体
	)
	if pushErr != nil {
		log.Printf("发送注册消息失败: %v", err)
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func ConsumeMessage(c *gin.Context) {
	appG := app.Gin{C: c}
	msgs, ch, err := rabbitmq.ConsumeMessage(
		"user_register",
		"direct",
		"user_register_q",
		"user.register",
		true, // autoAck
	)
	if err != nil {
		log.Printf("消费消息失败: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer ch.Close() // ⚠️ 记得释放资源

	select {
	case msg := <-msgs:
		appG.Response(200, 200, gin.H{"message": string(msg.Body)})
	default:
		c.JSON(200, gin.H{"message": "暂无消息"})
	}
}

func ConsumeAckMessage(c *gin.Context) {
	appG := app.Gin{C: c}
	msgs, ch, err := rabbitmq.ConsumeMessage(
		"user_register",
		"direct",
		"user_register_q",
		"user.register",
		false, // autoAck = false，手动ACK
	)
	if err != nil {
		appG.Response(500, 500, "注册消费者失败")
		return
	}

	// 异步处理消息
	go func() {
		defer ch.Close() // ⚠️ 释放资源（不要忘了）
		for msg := range msgs {
			log.Printf("收到消息: %s", msg.Body)

			if len(msg.Body) > 0 {
				// 处理成功，手动 ACK
				if err := msg.Ack(false); err != nil {
					log.Printf("ACK失败: %v", err)
				}
			} else {
				// 处理失败，NACK 不重回队列
				if err := msg.Nack(false, false); err != nil {
					log.Printf("NACK失败: %v", err)
				}
			}
		}
	}()

	appG.Response(200, 200, "消费者注册成功")
}

var dlxConfig = rabbitmq.DLXConfig{
	BusinessExchange:     "my.dlx.exchange",
	BusinessExchangeType: "direct",
	BusinessQueue:        "my.dlx.queue",
	BusinessRoutingKey:   "my.routing.key",
	DLXExchange:          "my.dlx.exchange",
	DLXRoutingKey:        "my.dlx.routing",
	TTL:                  10000, // 10秒
}

func SendDeadlineMessage(c *gin.Context) {
	appG := app.Gin{C: c}

	// 1. 初始化队列/交换机（只需一次，建议放到系统启动时）
	if err := rabbitmq.SetupDLX(dlxConfig); err != nil {
		log.Fatal("初始化 DLX 失败:", err)
	}

	// 2. 发送消息到业务队列
	err := rabbitmq.PublishDLXMessage(
		dlxConfig.BusinessExchange,
		dlxConfig.BusinessRoutingKey,
		"hello, world",
	)
	if err != nil {
		log.Fatal("发送消息失败:", err)
	}

	appG.Response(200, 200, "消息发送成功，等待进入死信队列")
}

func ConsumeDeadlineMessageOne(c *gin.Context) {
	log.Println("启动业务队列消费者：ConsumeDeadlineMessageOne")

	// 从业务队列消费消息，autoAck 设置为 false，确保需要手动确认消息
	msgs, ch, err := rabbitmq.ConsumeDLXMessages(dlxConfig.BusinessQueue, true)
	if err != nil {
		log.Printf("消费消息失败: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rabbitmq.SafeClose(ch) // 避免泄露

	select {
	case msg := <-msgs:
		log.Println("收到业务消息:", string(msg.Body))

		time.Sleep(11 * time.Second) // 模拟处理超时

		log.Println("消息已超时，拒绝消息并转入死信队列")
		msg.Nack(false, false) // 拒绝消息，不重回队列

		c.JSON(200, gin.H{"message": "消息已超时，转入死信队列"})

	default:
		log.Println("业务队列暂无消息")
		c.JSON(200, gin.H{"message": "暂无消息"})
	}
}

func ConsumeDeadlineMessageTwo(c *gin.Context) {
	log.Println("启动死信队列消费者：ConsumeDeadlineMessageTwo")

	msgs, ch, err := rabbitmq.ConsumeDLXMessages(dlxConfig.DLXExchange+"_queue", false)
	if err != nil {
		log.Printf("消费死信消息失败: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rabbitmq.SafeClose(ch)

	select {
	case msg := <-msgs:
		log.Println("收到死信消息:", string(msg.Body))
		msg.Ack(true) // 手动确认
		c.JSON(200, gin.H{"dead-letter-message": string(msg.Body)})
	default:
		log.Println("死信队列暂无消息")
		c.JSON(200, gin.H{"message": "暂无死信消息"})
	}
}
