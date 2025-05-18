package bot

import (
	"context"
	"github.com/tclutin/ticketly/telegram_bot/internal/config"
	"github.com/tclutin/ticketly/telegram_bot/internal/fsm"
	"github.com/tclutin/ticketly/telegram_bot/pkg/client/redis"
	"gopkg.in/telebot.v4"
	"log/slog"
	"time"
)

type Bot struct {
	bot *telebot.Bot
	fsm *fsm.FSM
}

func New() *Bot {
	cfg := config.MustLoad()

	redisClient := redis.NewClientRedis(cfg.Redis.Host, cfg.Redis.Port)
	fsmStore := fsm.NewRedisStore(redisClient)
	fsm := fsm.New(fsmStore)

	redisClient.Set(context.Background(), "хуй", "иди нахуй", 1*time.Hour)

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
		fsm: fsm,
	}
}

func (b *Bot) Run() {
	b.bot.Handle("/start1", func(c telebot.Context) error {
		return c.Send("Нажми кнопку, чтобы создать заявку:", &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			ReplyKeyboard: [][]telebot.ReplyButton{{
				{Text: "📨 Создать заявку"},
				{Text: "📨 Создатвь заявку"},
			}},
		})
	})

	b.bot.Handle(telebot.OnText, b.fsm.Middleware())

	b.bot.Start()
}
