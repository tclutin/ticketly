package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/ticketly/ticketly_api/internal/config"
	rabbitmq2 "github.com/tclutin/ticketly/ticketly_api/internal/delivery/rabbitmq"
	"github.com/tclutin/ticketly/ticketly_api/pkg/client/postgres"
	"github.com/tclutin/ticketly/ticketly_api/pkg/client/rabbitmq"
	"net/http"
)

type App struct {
	pgClient       *pgxpool.Pool
	httpServer     *http.Server
	rabbitmqClient *rabbitmq.Client
}

func New() *App {
	cfg := config.MustLoad()

	postgresClient := postgres.NewClient(context.Background(), cfg.Postgres.DSN())

	rabbitmqClient := rabbitmq.NewClient(
		cfg.RabbitMQ.URL,
		cfg.RabbitMQ.Exchange,
		cfg.RabbitMQ.ToClientQueue,
		cfg.RabbitMQ.ToOperatorQueue,
	)

	_ = rabbitmq2.NewPublisher(rabbitmqClient.Ch(), cfg.RabbitMQ.Exchange)

	return &App{
		pgClient:       postgresClient,
		rabbitmqClient: rabbitmqClient,
	}
}

func (a *App) Run() {

}

func (a *App) Stop() {
	a.rabbitmqClient.Close()
	a.pgClient.Close()
}
