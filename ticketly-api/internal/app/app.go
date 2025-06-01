package app

import (
	"context"
	"errors"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/ticketly/ticketly_api/internal/config"
	http2 "github.com/tclutin/ticketly/ticketly_api/internal/delivery/http"
	rabbitmq2 "github.com/tclutin/ticketly/ticketly_api/internal/delivery/rabbitmq"
	messageRepository "github.com/tclutin/ticketly/ticketly_api/internal/repository/message"
	ticketRepository "github.com/tclutin/ticketly/ticketly_api/internal/repository/ticket"
	userRepository "github.com/tclutin/ticketly/ticketly_api/internal/repository/user"
	ticketService "github.com/tclutin/ticketly/ticketly_api/internal/service/ticket"
	userService "github.com/tclutin/ticketly/ticketly_api/internal/service/user"
	"github.com/tclutin/ticketly/ticketly_api/pkg/client/centrifugo"
	"github.com/tclutin/ticketly/ticketly_api/pkg/client/postgres"
	"github.com/tclutin/ticketly/ticketly_api/pkg/client/rabbitmq"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		cfg.RabbitMQ.MLAnalysisQueue,
		cfg.RabbitMQ.MLResultsQueue,
	)

	casdoorClient := casdoorsdk.NewClient(
		cfg.Casdoor.Endpoint,
		cfg.Casdoor.ClientID,
		cfg.Casdoor.ClientSecret,
		cfg.Casdoor.Certificate,
		cfg.Casdoor.Organization,
		cfg.Casdoor.Application)

	centrifugoClient := centrifugo.New(cfg.Centrifugo.URL, cfg.Centrifugo.APIKey, cfg.Centrifugo.Secret)

	publisher := rabbitmq2.NewPublisher(rabbitmqClient.Ch(), cfg.RabbitMQ.Exchange)

	consumer := rabbitmq2.NewConsumer(rabbitmqClient.Ch())

	//users stuff
	userRepo := userRepository.NewRepository(postgresClient)

	userSrv := userService.NewService(userRepo)

	//messages stuff
	messageRepo := messageRepository.NewRepository(postgresClient)

	//tickets stuff
	ticketRepo := ticketRepository.NewRepository(postgresClient)

	ticketSrv := ticketService.NewService(ticketRepo, userRepo, messageRepo, publisher, consumer, centrifugoClient)
	if err := ticketSrv.ConsumeClients(context.Background()); err != nil {
		panic(err)
	}

	if err := ticketSrv.ConsumeMLResults(context.Background()); err != nil {
		panic(err)
	}

	router := http2.InitRouter(userSrv, ticketSrv, casdoorClient)

	return &App{
		pgClient:       postgresClient,
		rabbitmqClient: rabbitmqClient,
		httpServer: &http.Server{
			Addr:    net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
			Handler: router,
		},
	}
}

func (a *App) Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			//app.logger.Error("Server stopped with error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	a.Stop()
}

func (a *App) Stop() {
	a.rabbitmqClient.Close()
	a.pgClient.Close()

	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		//a.logger.Error("Error during server shutdown", "error", err)
		os.Exit(1)
	}
}
