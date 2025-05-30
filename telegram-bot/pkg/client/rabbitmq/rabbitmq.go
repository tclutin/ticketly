package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type Client struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func NewClient(url string, exchange string) *Client {
	conn, err := amqp091.Dial(url)
	if err != nil {
		log.Fatalln("failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()

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
