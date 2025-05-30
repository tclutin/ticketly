package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tclutin/ticketly/telegram_bot/internal/broker"
	"github.com/tclutin/ticketly/telegram_bot/internal/models"
	"github.com/tclutin/ticketly/telegram_bot/internal/storage"
	"github.com/tclutin/ticketly/telegram_bot/pkg/client/ticketly"
	"log/slog"
	"strconv"
)

var (
	ErrAlreadyTicketExists = errors.New("ticket already exists")
)

type sendToTelegram func(chatID int64, msg string) error

type Ticketly interface {
	CreateUser(externalId, username string) (uint64, error)
	CreateTicket(userId uint64, chatId int64, ticketType, content string) error
}

type TicketService struct {
	client         ticketly.Client
	storage        storage.Storage
	consumer       broker.EventConsumer
	sendToTelegram sendToTelegram
}

func NewTicketService(client ticketly.Client, storage storage.Storage, consumer broker.EventConsumer, sendToTelegram sendToTelegram) *TicketService {
	return &TicketService{client: client, storage: storage, consumer: consumer, sendToTelegram: sendToTelegram}
}

func (t *TicketService) CreateUser(externalId, username string) (uint64, error) {
	user, err := t.client.GetUserByExternalId(externalId)
	if err != nil {
		userId, err := t.client.Register(ticketly.RegisterUserRequest{
			ExternalID: externalId,
			Username:   username,
			Source:     ticketly.Telegram,
		})

		if err != nil {
			return 0, err
		}

		return userId, nil
	}

	return user.UserID, nil
}

func (t *TicketService) CreateTicket(userId uint64, chatId int64, ticketType, content string) error {
	tType := ticketly.TicketType(ticketType)

	ticketId, err := t.client.CreateTicket(ticketly.CreateTicketRequest{
		Type:    tType,
		UserID:  userId,
		Content: content,
	})

	if err != nil {
		return err
	}

	if tType == ticketly.RealtimeChat {
		if err = t.storage.SetChatMeta(context.Background(), chatId, models.RealtimeChatMeta{
			TicketID: ticketId,
			Type:     "realtime-chat",
			Status:   "open",
		}); err != nil {
			return err
		}
	}

	return nil
}

func (t *TicketService) ListenerOutgoing(ctx context.Context, queue string) error {
	msgs, err := t.consumer.Consume(queue)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var event models.TicketMessageEvent
			if err = json.Unmarshal(msg.Body, &event); err != nil {
				slog.Error("unmarshal failed:", slog.Any("error", err))
				continue
			}

			externalId, err := strconv.ParseInt(event.ExternalID, 10, 64)
			if err != nil {
				slog.Error("failed to parse external_id", slog.Any("error", err))
				continue
			}

			if event.Type == "realtime-chat" {
				if err = t.storage.SetChatMeta(context.Background(), externalId, models.RealtimeChatMeta{
					TicketID: event.TicketID,
					Type:     event.Type,
					Status:   event.Status,
				}); err != nil {
					slog.Error("failed to store", slog.Any("error", err))
					continue
				}
			}

			if event.Type == "realtime-chat" && event.Status == "closed" {
				if err := t.storage.DeleteChatMeta(context.Background(), externalId); err != nil {
					slog.Error("failed to delete", slog.Any("error", err))
					continue
				}
			}

			if err = t.sendToTelegram(externalId, event.Content); err != nil {
				//send to dead letter queue? tg user can block us
				slog.Error("telegram send failed",
					"external_id", externalId,
					"message", event.Content,
					"error", err)
				continue
			}
		}
	}()

	return nil
}
