package models

import "time"

type Message struct {
	MessageID  uint64
	TicketID   uint64
	SenderType string
	Content    string
	Sentiment  *string
	CreatedAt  time.Time
}

type MessagePreview struct {
	MessageID  uint64    `json:"message_id" db:"message_id"`
	TicketID   uint64    `json:"ticket_id" db:"ticket_id"`
	SenderType string    `json:"sender_type" db:"sender_type"`
	Content    string    `json:"text" db:"content"`
	Sentiment  *string   `json:"sentiment" db:"sentiment"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
