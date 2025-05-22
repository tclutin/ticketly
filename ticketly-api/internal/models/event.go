package models

type MessageEvent struct {
	TicketID uint64 `json:"ticket_id"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}
