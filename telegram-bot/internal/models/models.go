package models

type TicketMessageEvent struct {
	TicketID   uint64 `json:"ticket_id"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	Content    string `json:"content"`
}

type RealtimeChatMeta struct {
	TicketID uint64 `json:"ticket_id"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}
