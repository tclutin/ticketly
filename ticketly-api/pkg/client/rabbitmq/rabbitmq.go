package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type Client struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewClient(url string, exchange string, queue ...string) *Client {
	conn, err := amqp091.Dial(url)
	if err != nil {
		log.Fatalln("failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()

	err = ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalln("failed to declare exchange:", err)
	}

	for _, item := range queue {
		q, err := ch.QueueDeclare(
			item,
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			log.Fatalln("Failed to declare RabbitMQ queue:", err)
		}

		err = ch.QueueBind(
			q.Name,
			item,
			exchange,
			false,
			nil,
		)

		if err != nil {
			log.Fatalln("Failed to bind RabbitMQ queue:", err)
		}
	}
	return &Client{conn: conn, ch: ch}
}

func (c *Client) Conn() *amqp091.Connection {
	return c.conn
}

func (c *Client) Ch() *amqp091.Channel {
	return c.ch
}

func (c *Client) Close() {
	c.ch.Close()
	c.conn.Close()
}
