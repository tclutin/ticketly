package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Bot      Bot
	Redis    Redis
	RabbitMQ RabbitMQ
}

type Bot struct {
	Token   string
	Timeout time.Duration
}

type Redis struct {
	Host     string
	Port     string
	Password string
}

type RabbitMQ struct {
	URL      string `env:"RABBITMQ_URL"`
	Exchange string `env:"RABBITMQ_EXCHANGE"`
	//queue?
}

func MustLoad() *Config {
	// BotConfig
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("BOT_TOKEN is not set")
	}

	timeoutStr := os.Getenv("BOT_TIMEOUT")
	if timeoutStr == "" {
		log.Fatalln("BOT_TIMEOUT is not set")
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatalln("Failed to parse BOT_TIMEOUT:", err)
	}

	// RedisConfig
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Fatalln("REDIS_HOST is not set")
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		log.Fatalln("REDIS_PORT is not set")
	}

	// RabbitMQConfig
	rabbitmqHost := os.Getenv("RABBITMQ_URL")
	if rabbitmqHost == "" {
		log.Fatalln("RABBITMQ_URL is not set")
	}

	rabbitmqExchange := os.Getenv("RABBITMQ_EXCHANGE")
	if rabbitmqExchange == "" {
		log.Fatalln("RABBITMQ_EXCHANGE is not set")
	}

	return &Config{
		Bot: Bot{
			Token:   token,
			Timeout: timeout,
		},
		Redis: Redis{
			Host: redisHost,
			Port: redisPort,
		},
		RabbitMQ: RabbitMQ{
			URL:      rabbitmqHost,
			Exchange: rabbitmqExchange,
		},
	}
}
