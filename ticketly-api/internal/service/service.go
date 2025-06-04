package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/tclutin/ticketly/ticketly_api/internal/models"
	"github.com/tclutin/ticketly/ticketly_api/internal/service/operator"
	"github.com/tclutin/ticketly/ticketly_api/internal/service/ticket"
	"github.com/tclutin/ticketly/ticketly_api/internal/service/user"
)

type UserService interface {
	GetByExternalId(ctx context.Context, externalId string) (user models.User, err error)
	GetById(ctx context.Context, userId uint64) (user models.User, err error)
	Create(ctx context.Context, dto user.RegisterUserDTO) (uint64, error)
}

type TicketService interface {
	Create(ctx context.Context, model ticket.CreateTicketDTO) (uint64, error)
	GetAll(ctx context.Context) ([]models.PreviewTicket, error)
	Close(ctx context.Context, dto ticket.CloseTicketDTO) error
	Assign(ctx context.Context, dto ticket.AssignTicketDTO) (ticket.AssignedTicketDTO, error)
	GetHistory(ctx context.Context, ticketId uint64) ([]models.MessagePreview, error)
	SendMessage(ctx context.Context, dto ticket.SendMessageDTO) error
	GetActiveConnections(ctx context.Context, operatorId uint64) (ticket.ConnectionsDTO, error)
}

type OperatorService interface {
	GetById(ctx context.Context, operatorId uint64) (models.Operator, error)
	GetByCasdoorId(ctx context.Context, casdoorId uuid.UUID) (models.Operator, error)
	SyncOperator(ctx context.Context, dto operator.CreateOperatorDTO) (uint64, error)
}
