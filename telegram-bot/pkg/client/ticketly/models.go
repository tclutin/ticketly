package ticketly

import (
	"fmt"
	"unicode/utf8"
)

type TicketType string

const (
	OnlyMessage  TicketType = "only-message"
	RealtimeChat TicketType = "realtime-chat"
)

type Source string

const (
	Telegram Source = "telegram"
)

type CreateTicketRequest struct {
	Type    TicketType `json:"type"`
	UserID  uint64     `json:"user_id"`
	Content string     `json:"content"`
}

func (c CreateTicketRequest) Validate() error {
	if c.UserID == 0 {
		return fmt.Errorf("invalid user id")
	}

	if c.Type != OnlyMessage && c.Type != RealtimeChat {
		return fmt.Errorf("invalid type")
	}

	if utf8.RuneCountInString(c.Content) > 500 || utf8.RuneCountInString(c.Content) < 6 {
		return fmt.Errorf("invalid content")
	}

	return nil
}

type RegisterUserRequest struct {
	ExternalID string
	Username   string
	Source     Source
}

func (c RegisterUserRequest) Validate() error {
	if utf8.RuneCountInString(c.ExternalID) == 0 {
		return fmt.Errorf("invalid external id")
	}

	if utf8.RuneCountInString(c.Username) < 3 || utf8.RuneCountInString(c.Username) > 40 {
		return fmt.Errorf("invalid username")
	}

	if c.Source != Telegram {
		return fmt.Errorf("invalid source")
	}

	return nil
}
