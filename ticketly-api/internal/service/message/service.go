package message

import (
	"context"
	"github.com/tclutin/ticketly/ticketly_api/internal/models"
	"github.com/tclutin/ticketly/ticketly_api/internal/repository"
	"time"
)

type Service struct {
	repo repository.MessageRepository
}

func NewService(messageRepo repository.MessageRepository) *Service {
	return &Service{
		repo: messageRepo,
	}
}

func (s *Service) SendMessage(ctx context.Context, message models.Message) (uint64, error) {
	msg := models.Message{
		TicketID:   message.TicketID,
		SenderType: message.SenderType,
		Content:    message.Content,
		Sentiment:  message.Sentiment,
		CreatedAt:  time.Now().UTC(),
	}

	messageId, err := s.repo.Create(ctx, msg)
	if err != nil {
		return 0, err
	}

	return messageId, nil
}
