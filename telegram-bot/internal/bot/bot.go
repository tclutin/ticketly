package bot

import (
	"github.com/tclutin/ticketly/telegram_bot/internal/config"
	"github.com/tclutin/ticketly/telegram_bot/internal/handler"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
	"github.com/vitaliy-ukiru/fsm-telebot/v2/pkg/storage/memory"
	"github.com/vitaliy-ukiru/telebot-filter/v2/dispatcher"
	"gopkg.in/telebot.v4"
	"log/slog"
)

type Bot struct {
	bot *telebot.Bot
}

func New() *Bot {
	cfg := config.MustLoad()

	bot, err := telebot.NewBot(telebot.Settings{
		Token:     cfg.Bot.Token,
		Poller:    &telebot.LongPoller{Timeout: cfg.Bot.Timeout},
		ParseMode: telebot.ModeMarkdown,
	})

	if err != nil {
		slog.Error("failed to initialize telegram bot", slog.Any("error", err))
		return nil
	}

	return &Bot{
		bot: bot,
	}
}

func (b *Bot) Run() {
	g := b.bot.Group()

	dp := dispatcher.NewDispatcher(g)

	mn := fsm.New(memory.NewStorage())

	handler.NewHandler(dp, mn, b.bot).Register()

	b.bot.Start()
}
