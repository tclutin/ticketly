package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Bot   BotConfig
	Redis RedisConfig
}

type BotConfig struct {
	Token   string
	Timeout time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
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

	return &Config{
		Bot: BotConfig{
			Token:   token,
			Timeout: timeout,
		},
		Redis: RedisConfig{
			Host: redisHost,
			Port: redisPort,
		},
	}
}
