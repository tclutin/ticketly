package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/tclutin/ticketly/ticketly_api/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (uint64, error)
	GeByExternalId(ctx context.Context, externalId string) (models.User, error)
	GetById(ctx context.Context, userId uint64) (models.User, error)
}

type TicketRepository interface {
	Create(ctx context.Context, ticket models.Ticket) (uint64, error)
	Update(ctx context.Context, ticketId uint64, model models.Ticket) error
	GetTicketById(ctx context.Context, ticketId uint64) (models.Ticket, error)
	GetUserByTicketId(ctx context.Context, ticketId uint64) (models.User, error)
	GetInProgressRealtimeTickets(ctx context.Context, operatorId uint64) ([]models.Ticket, error)
	GetAllWithFirstMessage(ctx context.Context) ([]models.PreviewTicket, error)
	HasActiveRealtimeTicket(ctx context.Context, userID uint64) (bool, error)
}

type MessageRepository interface {
	Create(ctx context.Context, msg models.Message) (uint64, error)
	Update(ctx context.Context, messageId uint64, msg models.Message) error
	GetById(ctx context.Context, messageId uint64) (models.Message, error)
	GetAll(ctx context.Context, ticketId uint64) ([]models.MessagePreview, error)
}

type OperatorRepository interface {
	Create(ctx context.Context, operator models.Operator) (uint64, error)
	GetByCasdoorId(ctx context.Context, casdoorId uuid.UUID) (models.Operator, error)
	Update(ctx context.Context, operator models.Operator) error
	GetById(ctx context.Context, operatorId uint64) (models.Operator, error)
}
