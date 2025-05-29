package ticket

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/rabbitmq"
	"github.com/tclutin/ticketly/ticketly_api/internal/models"
	"github.com/tclutin/ticketly/ticketly_api/internal/repository"
	coreerrors "github.com/tclutin/ticketly/ticketly_api/internal/service/errors"
	"github.com/tclutin/ticketly/ticketly_api/pkg/client/centrifugo"
	"log/slog"
	"strconv"
	"time"
)

type Service struct {
	repo        repository.TicketRepository
	userRepo    repository.UserRepository
	messageRepo repository.MessageRepository
	publisher   rabbitmq.EventPublisher
	consumer    rabbitmq.EventConsumer
	centrifugo  *centrifugo.Client
}

func NewService(
	ticketRepo repository.TicketRepository,
	userRepo repository.UserRepository,
	messageRepo repository.MessageRepository,
	publisher rabbitmq.EventPublisher,
	consumer rabbitmq.EventConsumer,
	centrifugo *centrifugo.Client,

) *Service {
	return &Service{
		repo:        ticketRepo,
		userRepo:    userRepo,
		messageRepo: messageRepo,
		publisher:   publisher,
		consumer:    consumer,
		centrifugo:  centrifugo,
	}
}

// Create TODO refactor and tx
func (s *Service) Create(ctx context.Context, dto CreateTicketDTO) (uint64, error) {
	_, err := s.userRepo.GetById(ctx, dto.UserID)
	if err != nil {
		return 0, err
	}

	if dto.Type == "realtime-chat" {
		exists, err := s.repo.HasActiveRealtimeTicket(ctx, dto.UserID)
		if err != nil {
			return 0, err
		}
		if exists {
			return 0, coreerrors.ErrActiveTicketAlreadyExists
		}
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

func (s *Service) Assign(ctx context.Context, dto AssignTicketDTO) (AssignedTicketDTO, error) {
	ticket, err := s.repo.GetTicketById(ctx, dto.TicketID)
	if err != nil {
		return AssignedTicketDTO{}, err
	}

	if ticket.Status == "closed" || ticket.Status == "in_progress" {
		return AssignedTicketDTO{}, coreerrors.ErrTicketAlreadyClosedOrInProgress
	}

	if ticket.Type == "only-message" {
		return AssignedTicketDTO{}, coreerrors.ErrTicketWrongType
	}

	ticket.OperatorID = &dto.OperatorID
	ticket.Status = "in_progress"

	if err = s.repo.Update(ctx, dto.TicketID, ticket); err != nil {
		return AssignedTicketDTO{}, err
	}

	ticketEvent := models.TicketEvent{
		TicketID: ticket.TicketID,
		Status:   ticket.Status,
		Type:     ticket.Type,
		Content:  fmt.Sprintf("Оператор подключился к чату. Ваш тикет #%d.", ticket.TicketID),
	}

	if err = s.publisher.Publish("chat.outgoing", ticketEvent); err != nil {
		return AssignedTicketDTO{}, err
	}

	channel := s.ToChannel(ticket.TicketID)

	connToken, err := s.NewAccessToken(dto.OperatorID, 1*time.Hour)
	if err != nil {
		return AssignedTicketDTO{}, err
	}

	subToken, err := s.NewSubscriptionToken(dto.OperatorID, channel, 1*time.Hour)
	if err != nil {
		return AssignedTicketDTO{}, err
	}

	return AssignedTicketDTO{
		Channel:           channel,
		ConnectionToken:   connToken,
		SubscriptionToken: subToken,
	}, nil
}

func (s *Service) GetHistory(ctx context.Context, ticketId uint64) ([]models.MessagePreview, error) {
	return s.messageRepo.GetAll(ctx, ticketId)
}

func (s *Service) SendMessage(ctx context.Context, dto SendMessageDTO) error {
	ticket, err := s.repo.GetTicketById(ctx, dto.TicketID)
	if err != nil {
		return err
	}

	if ticket.Status == "closed" || ticket.Status == "open" {
		return coreerrors.ErrTicketWrongStatus
	}

	if *ticket.OperatorID != dto.OperatorID {
		return coreerrors.ErrOperatorNotAssigned
	}

	if ticket.Type == "only-message" {
		return coreerrors.ErrTicketWrongType
	}

	msg := models.Message{
		TicketID:   ticket.TicketID,
		SenderType: "operator",
		Content:    dto.Message,
		CreatedAt:  time.Now().UTC(),
	}

	messageId, err := s.messageRepo.Create(ctx, msg)
	if err != nil {
		return err
	}

	ticketEvent := models.TicketEvent{
		TicketID: ticket.TicketID,
		Status:   ticket.Status,
		Type:     ticket.Type,
		Content:  dto.Message,
	}

	if err = s.publisher.Publish("chat.outgoing", ticketEvent); err != nil {
		return err
	}

	messageEvent := models.MessagePreview{
		MessageID:  messageId,
		TicketID:   ticket.TicketID,
		Content:    dto.Message,
		SenderType: "operator",
		CreatedAt:  time.Now().UTC(),
	}

	if err = s.centrifugo.Publish(s.ToChannel(ticket.TicketID), messageEvent); err != nil {
		return err
	}

	return nil
}

func (s *Service) ConsumeClients(ctx context.Context) error {
	msgs, err := s.consumer.Consume("chat.incoming")
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var event models.TicketEvent
			if err = json.Unmarshal(msg.Body, &event); err != nil {
				slog.Error("unmarshal failed:", slog.Any("error", err))
				continue
			}

			_, err = s.repo.GetTicketById(ctx, event.TicketID)
			if err != nil {
				slog.Error("ticket not found:", slog.Any("error", err))
				continue
			}

			model := models.Message{
				TicketID:   event.TicketID,
				SenderType: "client",
				Content:    event.Content,
				CreatedAt:  time.Now().UTC(),
			}

			messageId, err := s.messageRepo.Create(ctx, model)
			if err != nil {
				slog.Error("failed to create message", slog.Any("error", err))
				continue
			}

			messageEvent := models.MessagePreview{
				MessageID:  messageId,
				TicketID:   event.TicketID,
				Content:    event.Content,
				SenderType: "client",
				CreatedAt:  time.Now().UTC(),
			}

			if err = s.centrifugo.Publish(s.ToChannel(event.TicketID), messageEvent); err != nil {
				slog.Error("failed to publish to centrifugo", slog.Any("error", err))
				continue
			}
		}
	}()

	return nil
}

func (s *Service) GetActiveConnections(ctx context.Context, operatorId uint64) (ConnectionsDTO, error) {
	tickets, err := s.repo.GetInProgressRealtimeTickets(ctx, operatorId)
	if err != nil {
		return ConnectionsDTO{}, err
	}

	var channels []ChannelInfo
	for _, ticket := range tickets {
		channel := s.ToChannel(ticket.TicketID)

		token, err := s.NewSubscriptionToken(operatorId, channel, 1*time.Hour)
		if err != nil {
			return ConnectionsDTO{}, err
		}

		channels = append(channels, ChannelInfo{
			Name:              channel,
			SubscriptionToken: token,
		})
	}

	token, err := s.NewAccessToken(operatorId, 1*time.Hour)
	if err != nil {
		return ConnectionsDTO{}, err
	}

	connections := ConnectionsDTO{
		ConnectionToken: token,
		Channels:        channels,
	}

	return connections, nil
}

func (s *Service) NewAccessToken(userID uint64, ttl time.Duration) (string, error) {
	claim := jwt.MapClaims{
		"exp": time.Now().UTC().Add(ttl).Unix(),
		"sub": strconv.FormatUint(userID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte("app.go"))
}

func (s *Service) NewSubscriptionToken(userID uint64, channel string, ttl time.Duration) (string, error) {
	claim := jwt.MapClaims{
		"exp":     time.Now().UTC().Add(ttl).Unix(),
		"sub":     strconv.FormatUint(userID, 10),
		"channel": channel,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte("app.go"))
}

// Close TODO refactor and tx
func (s *Service) Close(ctx context.Context, dto CloseTicketDTO) error {
	ticket, err := s.repo.GetTicketById(ctx, dto.TicketID)
	if err != nil {
		return err
	}

	if ticket.Status == "closed" {
		return coreerrors.ErrTicketAlreadyClosedOrInProgress
	}

	//operatorid, сделать проверку, чтобы другой человек не смог закрыть тикет

	now := time.Now().UTC()

	ticket.Status = "closed"
	ticket.ClosedAt = &now
	ticket.OperatorID = &dto.OperatorID

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

	event := models.TicketEvent{
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

func (s *Service) ToChannel(ticketId uint64) string {
	return fmt.Sprintf("ticket:%v", ticketId)
}
