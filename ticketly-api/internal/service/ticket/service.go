package ticket

import (
	"context"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/rabbitmq"
	"github.com/tclutin/ticketly/ticketly_api/internal/models"
	"github.com/tclutin/ticketly/ticketly_api/internal/repository"
	coreerrors "github.com/tclutin/ticketly/ticketly_api/internal/service/errors"
	"time"
)

type Service struct {
	repo        repository.TicketRepository
	userRepo    repository.UserRepository
	messageRepo repository.MessageRepository
	publisher   rabbitmq.EventPublisher
}

func NewService(
	ticketRepo repository.TicketRepository,
	userRepo repository.UserRepository,
	messageRepo repository.MessageRepository,
	publisher rabbitmq.EventPublisher,
) *Service {
	return &Service{
		repo:        ticketRepo,
		userRepo:    userRepo,
		messageRepo: messageRepo,
		publisher:   publisher,
	}
}

// Create TODO refactor and tx
func (s *Service) Create(ctx context.Context, dto CreateTicketDTO) (uint64, error) {
	_, err := s.userRepo.GetById(ctx, dto.UserID)
	if err != nil {
		return 0, err
	}

	ticket := models.Ticket{
		UserID:    dto.UserID,
		Type:      dto.Type,
		Status:    "open",
		CreatedAt: time.Now().UTC(),
	}

	ticketId, err := s.repo.Create(ctx, ticket)
	if err != nil {
		return 0, err
	}

	msg := models.Message{
		TicketID:   ticketId,
		SenderType: "user",
		Content:    dto.Content,
		CreatedAt:  time.Now().UTC(),
	}

	_, err = s.messageRepo.Create(ctx, msg)
	if err != nil {
		return 0, err
	}

	return ticketId, nil
}

func (s *Service) GetAll(ctx context.Context) ([]models.PreviewTicket, error) {
	return s.repo.GetAllWithFirstMessage(ctx)
}

// Close TODO refactor and tx
func (s *Service) Close(ctx context.Context, dto CloseTicketDTO) error {
	ticket, err := s.repo.GetTicketById(ctx, dto.TicketID)
	if err != nil {
		return err
	}

	if ticket.Status == "closed" {
		return coreerrors.ErrTicketAlreadyClosed
	}

	//operatorid, сделать проверку, чтобы другой человек не смог закрыть тикет
	if ticket.Type == "only-message" {
		ticket.OperatorID = &dto.OperatorID
	}

	now := time.Now().UTC()

	ticket.Status = "closed"
	ticket.ClosedAt = &now

	if err = s.repo.Update(ctx, ticket.TicketID, ticket); err != nil {
		return err
	}

	msg := models.Message{
		TicketID:   ticket.TicketID,
		SenderType: "operator",
		Content:    dto.Message,
		CreatedAt:  time.Now().UTC(),
	}

	_, err = s.messageRepo.Create(ctx, msg)
	if err != nil {
		return err
	}

	event := models.MessageEvent{
		TicketID: ticket.TicketID,
		Status:   ticket.Status,
		Type:     ticket.Type,
		Content:  dto.Message,
	}

	if err = s.publisher.Publish("chat.outgoing", event); err != nil {
		return err
	}

	//centrifuga disconnect?

	return nil
}
