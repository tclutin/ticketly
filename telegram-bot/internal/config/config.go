package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Bot BotConfig
}

type BotConfig struct {
	Token   string
	Timeout time.Duration
}

func MustLoad() *Config {
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

	return &Config{
		Bot: BotConfig{
			Token:   token,
			Timeout: timeout,
		},
	}
}
