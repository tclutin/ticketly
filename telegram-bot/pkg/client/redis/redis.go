package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"net"
)

func NewClient(host, port string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(host, port),
		Password: "",
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalln("fail to connect to redis:", err)
	}

	return client
}
