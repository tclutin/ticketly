package models

import "time"

type Ticket struct {
	TicketID   uint64
	UserID     uint64
	OperatorID *uint64
	Status     string
	Type       string
	Sentiment  *string
	CreatedAt  time.Time
	ClosedAt   *time.Time
}
