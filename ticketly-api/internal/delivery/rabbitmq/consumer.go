package rabbitmq

import "github.com/rabbitmq/amqp091-go"

type EventConsumer interface {
	Consume(queue string) (<-chan []amqp091.Delivery, error)
}

type Consumer struct {
	channel *amqp091.Channel
}

func NewConsumer(ch *amqp091.Channel) *Consumer {
	return &Consumer{channel: ch}
}

func (c *Consumer) Consume(queue string) (<-chan amqp091.Delivery, error) {
	return c.channel.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}
