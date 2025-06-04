package models

import "time"

type Ticket struct {
	TicketID   uint64     `db:"ticket_id"`
	UserID     uint64     `db:"user_id"`
	OperatorID *uint64    `db:"operator_id"`
	Status     string     `db:"status"`
	Type       string     `db:"type"`
	Sentiment  *string    `db:"sentiment"`
	CreatedAt  time.Time  `db:"created_at"`
	ClosedAt   *time.Time `db:"closed_at"`
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
