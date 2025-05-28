package service

import (
	"github.com/tclutin/ticketly/telegram_bot/internal/storage"
	"github.com/tclutin/ticketly/telegram_bot/pkg/client/ticketly"
)

type Ticketly interface {
	CreateUser(externalId, username string) (uint64, error)
}

type TicketService struct {
	client  ticketly.Client
	storage storage.Storage
}

func NewTicketService(client ticketly.Client) *TicketService {
	return &TicketService{client: client}
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

func (t *TicketService) CreateTicket(userId, chatId uint64, ticketType, content string) error {
	//ticketId, err := t.client.CreateTicket(ticketly.CreateTicketRequest{
	//	Type:    ticketly.TicketType(ticketType),
	//	UserID:  userId,
	//	Content: content,
	//})
	//
	//if err != nil {
	//	return err
	//}
	return nil
}
