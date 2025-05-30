package bot

import (
	"context"
	"github.com/tclutin/ticketly/telegram_bot/internal/broker"
	"github.com/tclutin/ticketly/telegram_bot/internal/config"
	"github.com/tclutin/ticketly/telegram_bot/internal/handler"
	"github.com/tclutin/ticketly/telegram_bot/internal/middleware"
	"github.com/tclutin/ticketly/telegram_bot/internal/service"
	"github.com/tclutin/ticketly/telegram_bot/internal/storage"
	rabbitmq2 "github.com/tclutin/ticketly/telegram_bot/pkg/client/rabbitmq"
	"github.com/tclutin/ticketly/telegram_bot/pkg/client/redis"
	"github.com/tclutin/ticketly/telegram_bot/pkg/client/ticketly"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/storage/memory"
	"github.com/vitaliy-ukiru/telebot-filter/v2/dispatcher"
	"gopkg.in/telebot.v4"
	"log/slog"
)

type Bot struct {
	bot      *telebot.Bot
	rabbitmq *rabbitmq2.Client
}

func New() *Bot {
	cfg := config.MustLoad()

	rabbitmqClient := rabbitmq2.NewClient(cfg.RabbitMQ.URL, cfg.RabbitMQ.Exchange)

	redisClient := redis.NewClient(cfg.Redis.Host, cfg.Redis.Port)

	redisStorage := storage.NewStorage(redisClient)

	publisher := broker.NewPublisher(rabbitmqClient.Ch(), cfg.RabbitMQ.Exchange)

	consumer := broker.NewConsumer(rabbitmqClient.Ch())

	client := ticketly.NewClient()

	bot, err := telebot.NewBot(telebot.Settings{
		Token:     cfg.Bot.Token,
		Poller:    &telebot.LongPoller{Timeout: cfg.Bot.Timeout},
		ParseMode: telebot.ModeMarkdown,
	})

	if err != nil {
		slog.Error("failed to initialize telegram bot", slog.Any("error", err))
		return nil
	}

	sendToTelegram := func(chatId int64, msg string) error {
		_, err := bot.Send(&telebot.Chat{ID: chatId}, msg)
		return err
	}

	srv := service.NewTicketService(client, redisStorage, consumer, sendToTelegram)

	if err = srv.ListenerOutgoing(context.Background(), "chat.outgoing"); err != nil {
		slog.Error("failed to initialize listener", slog.Any("error", err))
		return nil
	}

	g := bot.Group()

	dp := dispatcher.NewDispatcher(g)

	mn := fsm.New(memory.NewStorage())

	bot.Use(middleware.RedirectMiddleware(redisStorage, publisher), middleware.EnsureRegisteredMiddleware(srv))

	handler.NewHandler(dp, mn, bot, srv).Register()

	return &Bot{
		bot:      bot,
		rabbitmq: rabbitmqClient,
	}
}

func (b *Bot) Run() {
	b.bot.Start()
}

func (b *Bot) Stop() {
	//TODO
}
