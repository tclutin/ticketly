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
	// Стартовая команда
	b.bot.Handle("/start", func(c telebot.Context) error {
		userID := c.Sender().ID
		_ = b.fsm.Set(userID, "select_type")
		return c.Send("Выбери тип тикета: чат или одно сообщение?")
	})

	// FSM: выбор типа тикета
	b.fsm.Register("select_type", func(c telebot.Context, userID int64) error {
		// можно сохранить тип тикета отдельно в Redis
		_ = b.fsm.Set(userID, "select_category")
		return c.Send("Выбери категорию тикета:")
	})

	// FSM: выбор категории
	b.fsm.Register("select_category", func(c telebot.Context, userID int64) error {
		_ = b.fsm.Set(userID, "write_message")
		return c.Send("Введите текст обращения:")
	})

	// FSM: текст обращения
	b.fsm.Register("write_message", func(c telebot.Context, userID int64) error {
		_ = b.fsm.Clear(userID)
		return c.Send("✅ Заявка принята!")
	})

	b.bot.Handle("/start1", func(c telebot.Context) error {
		return c.Send("Нажми кнопку, чтобы создать заявку:", &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			ReplyKeyboard: [][]telebot.ReplyButton{{
				{Text: "📨 Создать заявку"},
				{Text: "📨 Создатвь заявку"},
			}},
		})
	})
	// FSM Middleware — на каждое текстовое сообщение
	b.bot.Handle(telebot.OnText, b.fsm.Middleware())

	b.bot.Start()
}
