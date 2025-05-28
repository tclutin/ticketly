package bot

import (
	"github.com/tclutin/ticketly/telegram_bot/internal/config"
	"github.com/tclutin/ticketly/telegram_bot/internal/handler"
	"github.com/tclutin/ticketly/telegram_bot/internal/middleware"
	"github.com/tclutin/ticketly/telegram_bot/internal/service"
	"github.com/tclutin/ticketly/telegram_bot/pkg/client/ticketly"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/storage/memory"
	"github.com/vitaliy-ukiru/telebot-filter/v2/dispatcher"
	"gopkg.in/telebot.v4"
	"log/slog"
)

type Bot struct {
	bot    *telebot.Bot
	client ticketly.Client
}

func New() *Bot {
	cfg := config.MustLoad()

	client := ticketly.NewClient()

	srv := service.NewTicketService(client)

	bot, err := telebot.NewBot(telebot.Settings{
		Token:     cfg.Bot.Token,
		Poller:    &telebot.LongPoller{Timeout: cfg.Bot.Timeout},
		ParseMode: telebot.ModeMarkdown,
	})

	if err != nil {
		slog.Error("failed to initialize telegram bot", slog.Any("error", err))
		return nil
	}

	g := bot.Group()

	dp := dispatcher.NewDispatcher(g)

	mn := fsm.New(memory.NewStorage())

	bot.Use(middleware.EnsureRegisteredMiddleware(srv))

	handler.NewHandler(dp, mn, bot, srv).Register()

	return &Bot{
		bot:    bot,
		client: client,
	}
}

func (b *Bot) Run() {
	b.bot.Start()
}
