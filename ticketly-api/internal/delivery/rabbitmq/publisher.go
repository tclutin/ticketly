package rabbitmq

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
)

type EventPublisher interface {
	Publisher(routingKey string) (<-chan amqp091.Delivery, error)
}

type Publisher struct {
	channel  *amqp091.Channel
	exchange string
}

func NewPublisher(ch *amqp091.Channel, exchange string) *Publisher {
	return &Publisher{
		channel:  ch,
		exchange: exchange,
	}
}

func (p *Publisher) Publish(routingKey string, event any) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.channel.Publish(p.exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        bytes,
	})
}
