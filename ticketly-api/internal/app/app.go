package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/ticketly/ticketly_api/internal/config"
	"github.com/tclutin/ticketly/ticketly_api/pkg/client/postgres"
	"net/http"
)

type App struct {
	pgClient   *pgxpool.Pool
	httpServer *http.Server
}

func New() *App {
	cfg := config.MustLoad()
	fmt.Println(cfg.Postgres.DSN())

	postgresClient := postgres.NewClient(context.Background(), cfg.Postgres.DSN())

	return &App{
		pgClient: postgresClient,
	}
}

func (a *App) Run() {

}

func (a *App) Stop() {
	a.pgClient.Close()
}
