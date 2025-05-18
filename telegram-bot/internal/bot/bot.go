package bot

import (
	"github.com/tclutin/ticketly/telegram_bot/internal/config"
	"gopkg.in/telebot.v4"
	"log/slog"
)

type Bot struct {
	bot *telebot.Bot
}

func New() *Bot {
	cfg := config.MustLoad()

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.Bot.Token,
		Poller: &telebot.LongPoller{Timeout: cfg.Bot.Timeout},
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
	b.bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Нажми кнопку, чтобы создать заявку:", &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			ReplyKeyboard: [][]telebot.ReplyButton{{
				{Text: "📨 Создать заявку"},
				{Text: "📨 Создатвь заявку"},
			}},
		})
	})

	b.bot.Start()
}
