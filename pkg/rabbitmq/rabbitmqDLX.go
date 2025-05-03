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

func SetupDLX(config DLXConfig) error {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer SafeClose(ch)

	if err := DeclareDLXExchange(ch, config.DLXExchange, "direct"); err != nil {
		return fmt.Errorf("declare DLX exchange failed: %w", err)
	}
	if err := DeclareDLXQueue(ch, config.DLXExchange+"_queue", nil); err != nil {
		return fmt.Errorf("declare DLX queue failed: %w", err)
	}
	if err := BindDLXQueueToExchange(ch, config.DLXExchange+"_queue", config.DLXExchange, config.DLXRoutingKey); err != nil {
		return fmt.Errorf("bind DLX queue failed: %w", err)
	}

	if err := DeclareDLXExchange(ch, config.BusinessExchange, config.BusinessExchangeType); err != nil {
		return fmt.Errorf("declare business exchange failed: %w", err)
	}
	bizArgs := BuildBusinessQueueArgs(config)
	if err := DeclareDLXQueue(ch, config.BusinessQueue, bizArgs); err != nil {
		return fmt.Errorf("declare business queue failed: %w", err)
	}
	if err := BindDLXQueueToExchange(ch, config.BusinessQueue, config.BusinessExchange, config.BusinessRoutingKey); err != nil {
		return fmt.Errorf("bind business queue failed: %w", err)
	}

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
