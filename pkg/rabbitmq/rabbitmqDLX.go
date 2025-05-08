package rabbitmq

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

// DLXConfig 用于统一描述业务队列及其死信配置
type DLXConfig struct {
	BusinessExchange     string
	BusinessExchangeType string
	BusinessQueue        string
	BusinessRoutingKey   string

	DLXExchange   string
	DLXRoutingKey string

	TTL int32 // 单位：毫秒
}

// =====================
// 工具函数：安全关闭通道
// =====================

func SafeClose(ch *amqp091.Channel) {
	if ch != nil {
		_ = ch.Close()
	}
}

// =====================
// 1. 队列/交换机声明（每次使用新 channel）
// =====================

func DeclareDLXExchange(ch *amqp091.Channel, name, kind string) error {
	return ch.ExchangeDeclare(
		name, kind,
		true, false, false, false,
		nil,
	)
}

func DeclareDLXQueue(ch *amqp091.Channel, name string, args amqp091.Table) error {
	_, err := ch.QueueDeclare(
		name, true, false, false, false, args,
	)
	if err != nil {
		return fmt.Errorf("declare queue [%s] failed: %w", name, err)
	}
	return nil
}

func BindDLXQueueToExchange(ch *amqp091.Channel, queue, exchange, routingKey string) error {
	return ch.QueueBind(queue, routingKey, exchange, false, nil)
}

// =====================
// 2. 配置封装方法
// =====================

func BuildBusinessQueueArgs(config DLXConfig) amqp091.Table {
	args := amqp091.Table{}
	if config.TTL > 0 {
		args["x-message-ttl"] = config.TTL
	}
	args["x-dead-letter-exchange"] = config.DLXExchange
	if config.DLXRoutingKey != "" {
		args["x-dead-letter-routing-key"] = config.DLXRoutingKey
	}
	return args
}

// =====================
// 3. Setup 封装（每次用新 channel）
// =====================

// SetupDLX 初始化业务队列及其死信队列（DLX）结构，包括交换机、队列和绑定关系
func SetupDLX(config DLXConfig) error {
	// 创建一个新的 channel，每次使用都临时打开，保证并发安全
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer SafeClose(ch) // 使用封装的安全关闭，防止 nil 或已关闭的 channel 造成 panic

	// =============================
	// 1. 声明死信交换机（DLX Exchange）
	// =============================
	if err := DeclareDLXExchange(ch, config.DLXExchange, "direct"); err != nil {
		return fmt.Errorf("declare DLX exchange failed: %w", err)
	}

	// =============================
	// 2. 声明死信队列（DLX Queue）
	//    用于接收所有被拒绝、过期或无法路由的消息
	// =============================
	if err := DeclareDLXQueue(ch, config.DLXExchange+"_queue", nil); err != nil {
		return fmt.Errorf("declare DLX queue failed: %w", err)
	}

	// =============================
	// 3. 绑定死信队列到死信交换机
	//    绑定 routingKey，用于匹配转入死信的消息
	// =============================
	if err := BindDLXQueueToExchange(ch, config.DLXExchange+"_queue", config.DLXExchange, config.DLXRoutingKey); err != nil {
		return fmt.Errorf("bind DLX queue failed: %w", err)
	}

	// =============================
	// 4. 声明业务交换机（Business Exchange）
	//    用于正常业务流程的消息投递
	// =============================
	if err := DeclareDLXExchange(ch, config.BusinessExchange, config.BusinessExchangeType); err != nil {
		return fmt.Errorf("declare business exchange failed: %w", err)
	}

	// =============================
	// 5. 构建业务队列的参数
	//    包括 TTL（延迟时间）和死信交换机配置
	//	  这一步将普通队列和死信队列绑定
	// =============================
	bizArgs := BuildBusinessQueueArgs(config)

	// =============================
	// 6. 声明业务队列
	//    消息将暂存在此队列中直到过期或被消费
	// =============================
	if err := DeclareDLXQueue(ch, config.BusinessQueue, bizArgs); err != nil {
		return fmt.Errorf("declare business queue failed: %w", err)
	}

	// =============================
	// 7. 绑定业务队列到业务交换机
	//    routingKey 控制业务消息投递路径
	// =============================
	if err := BindDLXQueueToExchange(ch, config.BusinessQueue, config.BusinessExchange, config.BusinessRoutingKey); err != nil {
		return fmt.Errorf("bind business queue failed: %w", err)
	}

	// =============================
	// 8. 输出成功日志
	// =============================
	log.Printf("DLX setup success: Biz [%s → %s], DLX [%s → %s]",
		config.BusinessExchange, config.BusinessQueue,
		config.DLXExchange, config.DLXExchange+"_queue")

	return nil
}

// =====================
// 4. 消费与发布
// =====================

// PublishDLXMessage 发布消息到交换机，内部自动关闭 channel
func PublishDLXMessage(exchange, routingKey, body string) error {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer SafeClose(ch)

	return ch.Publish(
		exchange, routingKey, false, false,
		amqp091.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp091.Persistent,
			Body:         []byte(body),
		},
	)
}

// ConsumeDLXMessages 返回消息 channel 与底层 amqp channel，调用方必须在消费完成后手动关闭 ch
func ConsumeDLXMessages(queue string, autoAck bool) (<-chan amqp091.Delivery, *amqp091.Channel, error) {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create channel: %w", err)
	}

	msgs, err := ch.Consume(
		queue, "", autoAck, false, false, false, nil,
	)
	if err != nil {
		SafeClose(ch)
		return nil, nil, fmt.Errorf("failed to consume messages: %w", err)
	}
	return msgs, ch, nil
}

// GetOneDLXMessage 获取队列中一条消息，返回 msg 及 channel，调用方负责关闭 ch
func GetOneDLXMessage(queue string, autoAck bool) (*amqp091.Delivery, *amqp091.Channel, error) {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create channel: %w", err)
	}

	msg, ok, err := ch.Get(queue, autoAck)
	if err != nil {
		SafeClose(ch)
		return nil, nil, fmt.Errorf("get message error: %w", err)
	}
	if !ok {
		SafeClose(ch)
		return nil, nil, nil
	}
	return &msg, ch, nil
}
