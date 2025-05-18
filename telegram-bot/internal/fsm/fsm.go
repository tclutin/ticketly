package fsm

import (
	"gopkg.in/telebot.v4"
	"log/slog"
)

type HandlerFunc func(c telebot.Context) error

type Handlers map[string]HandlerFunc

type Store interface {
	Set(userID int64, state string) error
	Get(userID int64) (string, error)
	Clear(userID int64) error
}

type FSM struct {
	handlers Handlers
	store    Store
}

func New(store Store) *FSM {
	return &FSM{
		handlers: make(Handlers),
		store:    store,
	}
}

func (f *FSM) Register(state string, handler HandlerFunc) {
	f.handlers[state] = handler
}

func (f *FSM) Set(userID int64, state string) error {
	return f.store.Set(userID, state)
}

func (f *FSM) Clear(userID int64) error {
	return f.store.Clear(userID)
}

func (f *FSM) Middleware() telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userID := c.Sender().ID

		state, err := f.store.Get(userID)
		if err != nil {
			username := c.Sender().Username
			chatID := c.Message().Chat.ID

			slog.Error("failed to get user state",
				slog.Int64("telegram_user_id", userID),
				slog.String("telegram_username", username),
				slog.Int64("telegram_chat_id", chatID),
				slog.Any("error", err),
			)
			return c.Send("Упс! Что-то пошло не так... Попробуйте ещё раз 🙈")
		}

		if handler, ok := f.handlers[state]; ok {
			return handler(c)
		}

		return c.Send("Напиши /start чтобы начать.")
	}
}
