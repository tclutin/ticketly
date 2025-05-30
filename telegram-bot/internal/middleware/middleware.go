package middleware

import (
	"context"
	"github.com/tclutin/ticketly/telegram_bot/internal/broker"
	"github.com/tclutin/ticketly/telegram_bot/internal/keyboard"
	"github.com/tclutin/ticketly/telegram_bot/internal/models"
	"github.com/tclutin/ticketly/telegram_bot/internal/service"
	"github.com/tclutin/ticketly/telegram_bot/internal/storage"
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

			externalId := strconv.FormatInt(c.Sender().ID, 10)
			username := c.Sender().Username

			if len(username) == 0 {
				username = "JazzLord"
			}

			userId, err := srv.CreateUser(externalId, username)
			if err != nil {
				slog.Error("failed to ensure user registration",
					slog.String("externalId", externalId),
					slog.String("username", username),
					slog.Any("error", err))
				return c.Send("❗ Произошла ошибка при регистрации. Попробуйте позже.")
			}

			c.Set("user_id", userId)

			return next(c)
		}
	}
}

func RedirectMiddleware(storage storage.Storage, publisher broker.EventPublisher) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if c.Sender() == nil {
				return next(c)
			}

			switch c.Text() {
			case "/start",
				keyboard.BtnCreateTicket,
				keyboard.BtnBack,
				keyboard.BtnCancel,
				keyboard.BtnConfirmYes,
				keyboard.BtnConfirmNo,
				keyboard.BtnTicketTypeSingle,
				keyboard.BtnTicketTypeRealtime:
				return next(c)
			}

			externalId := c.Sender().ID

			meta, err := storage.GetChatMeta(context.Background(), externalId)
			if err != nil {
				return err
			}

			if meta.Type == "realtime-chat" && meta.Status == "in_progress" {
				err = publisher.Publish("chat.incoming", models.TicketMessageEvent{
					TicketID:   meta.TicketID,
					ExternalID: strconv.FormatInt(externalId, 10),
					Status:     "in_progress",
					Type:       meta.Type,
					Content:    c.Text(),
				})

				if err != nil {
					return err
				}
			}

			return next(c)
		}
	}
}
