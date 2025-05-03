package rabbitmq

import (
	"fmt"
	"log"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/rabbitmq/amqp091-go"
)

var RabbitMQConn *amqp091.Connection

// Setup 初始化 RabbitMQ 连接
func Setup() error {
	conn, err := amqp091.Dial(setting.RabbitMQSetting.Host)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	RabbitMQConn = conn
	log.Println("RabbitMQ connection established.")
	return nil
}

// DeclareQueue 声明队列
func DeclareQueue(queueName string, durable bool, otherArgs amqp091.Table) error {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, durable,
		false, false, false, otherArgs,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	log.Printf("Queue declared: %s", queueName)
	return nil
}

// DeclareExchange 声明交换机
func DeclareExchange(exchangeName, exchangeType string) error {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, exchangeType,
		true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}
	log.Printf("Exchange declared: %s (%s)", exchangeName, exchangeType)
	return nil
}

// BindQueueToExchange 绑定队列到交换机
func BindQueueToExchange(queueName, exchangeName, routingKey string) error {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind queue to exchange: %w", err)
	}
	log.Printf("Bound queue [%s] to exchange [%s] with routing key [%s]", queueName, exchangeName, routingKey)
	return nil
}

// Publish 发布消息
func Publish(exchangeName, exchangeType, routingKey, message string) error {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	err = ch.Publish(
		exchangeName, routingKey,
		false, false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}
	log.Printf("Published to exchange [%s] routingKey [%s]: %s", exchangeName, routingKey, message)
	return nil
}

// PublishMessage 高级封装：声明+绑定+发布
func PublishMessage(exchangeName, exchangeType, queueName, routingKey, message string) error {
	if err := DeclareExchange(exchangeName, exchangeType); err != nil {
		return err
	}

	if queueName != "" {
		if err := DeclareQueue(queueName, true, nil); err != nil {
			return err
		}
		if err := BindQueueToExchange(queueName, exchangeName, routingKey); err != nil {
			return err
		}
	}

	return Publish(exchangeName, exchangeType, routingKey, message)
}

// ConsumeMessage 消费者注册（自动声明交换机/队列/绑定）
func ConsumeMessage(exchangeName, exchangeType, queueName, routingKey string, autoAck bool) (<-chan amqp091.Delivery, *amqp091.Channel, error) {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		return nil, nil, err
	}

	if exchangeType == "fanout" {
		routingKey = ""
	}

	if err := ch.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil); err != nil {
		ch.Close()
		return nil, nil, fmt.Errorf("declare exchange: %w", err)
	}

	if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		ch.Close()
		return nil, nil, fmt.Errorf("declare queue: %w", err)
	}

	if err := ch.QueueBind(queueName, routingKey, exchangeName, false, nil); err != nil {
		ch.Close()
		return nil, nil, fmt.Errorf("bind queue: %w", err)
	}

	msgs, err := ch.Consume(
		queueName, "", autoAck, false, false, false, nil,
	)
	if err != nil {
		ch.Close()
		return nil, nil, fmt.Errorf("consume: %w", err)
	}

	log.Printf("Consuming from queue [%s] (exchange: %s, routingKey: %s)", queueName, exchangeName, routingKey)
	return msgs, ch, nil // 返回 channel 供调用者控制生命周期
}

func ConsumeMessageWithAck(exchangeName, exchangeType, queueName, routingKey string) (<-chan amqp091.Delivery, *amqp091.Channel, error) {
	return ConsumeMessage(exchangeName, exchangeType, queueName, routingKey, false)
}

// Close 关闭连接
func Close() {
	if RabbitMQConn != nil {
		_ = RabbitMQConn.Close()
	}
}
