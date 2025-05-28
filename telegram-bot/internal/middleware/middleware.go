package middleware

import (
	"github.com/tclutin/ticketly/telegram_bot/internal/service"
	"gopkg.in/telebot.v4"
	"log/slog"
	"strconv"
)

func EnsureRegisteredMiddleware(srv service.Ticketly) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if c.Sender() == nil {
				return next(c)
			}

			userId := strconv.FormatInt(c.Sender().ID, 10)
			username := c.Sender().Username

			_, err := srv.CreateUser(userId, username)
			if err != nil {
				slog.Error("failed to ensure user registration",
					slog.String("user_id", userId),
					slog.String("username", username),
					slog.Any("error", err))
				return c.Send("❗ Произошла ошибка при регистрации. Попробуйте позже.")
			}

			return next(c)
		}
	}
}
