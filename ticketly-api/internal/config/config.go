package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type Config struct {
	HTTP       HTTPServer
	Postgres   Postgres
	RabbitMQ   RabbitMQ
	Centrifugo Centrifugo
}

type HTTPServer struct {
	Host string `env:"HTTP_HOST"`
	Port string `env:"HTTP_PORT"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DATABASE"`
}

type RabbitMQ struct {
	URL             string `env:"RABBITMQ_URL"`
	Exchange        string `env:"RABBITMQ_EXCHANGE"`
	ToOperatorQueue string `env:"RABBITMQ_TO_OPERATOR_QUEUE"`
	ToClientQueue   string `env:"RABBITMQ_TO_CLIENT_QUEUE"`
}

type Centrifugo struct {
	URL    string        `env:"CENTRIFUGO_API_URL"`
	APIKey string        `env:"CENTRIFUGO_API_KEY"`
	Secret string        `env:"CENTRIFUGO_JWT_SECRET"`
	TTL    time.Duration `env:"CENTRIFUGO_JWT_TTL"`
}

func MustLoad() *Config {
	var config Config

	if err := godotenv.Load(); err != nil {
		log.Fatalln("failed to load .env file:", err)
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalln("failed to load .env file:", err)
	}
	return &config
}

func (p Postgres) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.Database)
}
