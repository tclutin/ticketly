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

type PreviewTicket struct {
	TicketID  uint64     `json:"ticket_id" db:"ticket_id"`
	Type      string     `json:"type" db:"type"`
	Status    string     `json:"status" db:"status"`
	Preview   string     `json:"preview" db:"preview"`
	Sentiment *string    `json:"sentiment" db:"sentiment"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	ClosedAt  *time.Time `json:"closed_at" db:"closed_at"`
}
