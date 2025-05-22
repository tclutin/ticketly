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
