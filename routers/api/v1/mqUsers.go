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
	msgs, err := rabbitmq.ConsumeMessage(
		"user_register",
		"direct",
		"user_register_q",
		"user.register",
		true,
	)
	if err != nil {
		log.Printf("消费消息失败: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	select {
	case msg := <-msgs:
		appG.Response(200, 200, gin.H{"message": string(msg.Body)})
	default:
		c.JSON(200, gin.H{"message": "暂无消息"})
	}
}

func ConsumeAckMessage(c *gin.Context) {
	appG := app.Gin{C: c}
	msgs, err := rabbitmq.ConsumeMessageWithAck(
		"user_register",
		"direct",
		"user_register_q",
		"user.register",
	)
	if err != nil {
		appG.Response(500, 500, "注册消费者失败")
	}

	go func() {
		for msg := range msgs {
			log.Printf("收到消息: %s", msg.Body)
			// 处理业务逻辑
			if len(msg.Body) > 0 {
				// 处理成功，手动 ACK
				if err := msg.Ack(false); err != nil {
					log.Printf("ACK失败: %v", err)
				}
			} else {
				// 处理失败，拒绝消息，不重回队列
				if err := msg.Nack(false, false); err != nil {
					log.Printf("NACK失败: %v", err)
				}
			}
		}
	}()

}
