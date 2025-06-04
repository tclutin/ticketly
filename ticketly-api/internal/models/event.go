package models

type TicketMessageEvent struct {
	TicketID   uint64 `json:"ticket_id"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	Content    string `json:"content"`
}

type MlAnalysisEvent struct {
	MessageID uint64 `json:"message_id"`
	TicketID  uint64 `json:"ticket_id"`
	Content   string `json:"content"`
}

type MlResultEvent struct {
	MessageID  uint64  `json:"message_id"`
	TicketID   uint64  `json:"ticket_id"`
	Sentiment  string  `json:"sentiment"`
	Confidence float64 `json:"confidence"`
}
