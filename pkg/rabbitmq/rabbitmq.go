package rabbitmq

import (
	"fmt"
	"log"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/rabbitmq/amqp091-go"
)

var RabbitMQConn *amqp091.Connection
var RabbitMQChannel *amqp091.Channel

// Setup 初始化 RabbitMQ 连接和通道
func Setup() error {
	// 连接 RabbitMQ
	conn, err := amqp091.Dial(setting.RabbitMQSetting.Host)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	RabbitMQConn = conn

	// 创建通道
	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	RabbitMQChannel = channel

	log.Println("RabbitMQ channel created successfully!")
	return nil
}

// DeclareQueue 声明队列（支持多个队列动态注册）
func DeclareQueue(queueName string) error {
	_, err := RabbitMQChannel.QueueDeclare(
		queueName, // 队列名
		false,     // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	log.Printf("Queue declared: %s", queueName)
	return nil
}

// DeclareExchange 声明交换机
func DeclareExchange(exchangeName, exchangeType string) error {
	err := RabbitMQChannel.ExchangeDeclare(
		exchangeName, // 交换机名
		exchangeType, // 交换机类型 (direct, fanout, topic)
		true,         // durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}
	log.Printf("Exchange declared: %s with type: %s", exchangeName, exchangeType)
	return nil
}

// BindQueueToExchange 将队列绑定到交换机
func BindQueueToExchange(queueName, exchangeName, routingKey string) error {
	err := RabbitMQChannel.QueueBind(
		queueName,    // 队列名
		routingKey,   // 路由键
		exchangeName, // 交换机名
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %w", err)
	}
	log.Printf("Queue %s bound to exchange %s with routing key %s", queueName, exchangeName, routingKey)
	return nil
}

// Publish 使用指定交换机发布消息
func Publish(exchangeName, exchangeType, routingKey, message string) error {
	// 声明交换机
	if err := DeclareExchange(exchangeName, exchangeType); err != nil {
		return fmt.Errorf("failed to declare exchange before publishing: %w", err)
	}

	err := RabbitMQChannel.Publish(
		exchangeName, // 交换机名
		routingKey,   // 路由键
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	log.Printf("Message sent to exchange [%s] with routing key [%s]: %s", exchangeName, routingKey, message)
	return nil
}

// Consume 消费消息
func Consume(exchangeName, exchangeType, queueName, routingKey string) (<-chan amqp091.Delivery, error) {
	// 声明交换机
	if err := DeclareExchange(exchangeName, exchangeType); err != nil {
		return nil, fmt.Errorf("failed to declare exchange before consuming: %w", err)
	}

	// 声明队列
	if err := DeclareQueue(queueName); err != nil {
		return nil, fmt.Errorf("failed to declare queue before consuming: %w", err)
	}

	// 绑定队列到交换机
	if err := BindQueueToExchange(queueName, exchangeName, routingKey); err != nil {
		return nil, fmt.Errorf("failed to bind queue to exchange: %w", err)
	}

	// 注册消费者
	msgs, err := RabbitMQChannel.Consume(
		queueName, // 队列名
		"",        // consumer tag
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	log.Printf("Consuming messages from queue %s on exchange %s with routing key %s", queueName, exchangeName, routingKey)
	return msgs, nil
}

// Close 关闭连接和通道
func Close() {
	if RabbitMQChannel != nil {
		_ = RabbitMQChannel.Close()
	}
	if RabbitMQConn != nil {
		_ = RabbitMQConn.Close()
	}
}

// PublishMessage 是统一的高级封装：自动声明交换机、队列、绑定并发布消息
func PublishMessage(exchangeName, exchangeType, queueName, routingKey, message string) error {
	// 声明交换机
	if err := DeclareExchange(exchangeName, exchangeType); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// 声明队列（如果有传）
	if queueName != "" {
		if err := DeclareQueue(queueName); err != nil {
			return fmt.Errorf("failed to declare queue: %w", err)
		}

		// 绑定队列（只有非 fanout 类型才需要 routingKey）
		if err := BindQueueToExchange(queueName, exchangeName, routingKey); err != nil {
			return fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	// 发布消息
	if err := Publish(exchangeName, exchangeType, routingKey, message); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func ConsumeMessageWithAck(exchangeName, exchangeType, queueName, routingKey string) (<-chan amqp091.Delivery, error) {
	return ConsumeMessage(exchangeName, exchangeType, queueName, routingKey, false)
}

// ConsumeMessage 是高级封装：自动声明交换机、队列、绑定并注册消费者
func ConsumeMessage(exchangeName, exchangeType, queueName, routingKey string, autoAck bool) (<-chan amqp091.Delivery, error) {
	// fanout 模式不使用 routingKey
	if exchangeType == "fanout" {
		routingKey = ""
	}

	if err := DeclareExchange(exchangeName, exchangeType); err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	if err := DeclareQueue(queueName); err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := BindQueueToExchange(queueName, exchangeName, routingKey); err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	// 注册消费者，autoAck 由参数控制
	msgs, err := RabbitMQChannel.Consume(
		queueName,
		"",
		autoAck, // 这里使用传入参数
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Consuming from [%s] on exchange [%s] with routing key [%s], autoAck=%v", queueName, exchangeName, routingKey, autoAck)
	return msgs, nil
}
