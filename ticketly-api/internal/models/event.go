package models

//
//const (
//	EventTypeMessage      = "new_message"
//	EventTypeStatusUpdate = "status_update"
//	EventTypeSentiment    = "sentiment_update"
//)
//
//type Event struct {
//	TicketID  uint64          `json:"ticket_id"`
//	EventType string          `json:"event_type"`
//	Payload   json.RawMessage `json:"payload"`
//}
//
//type TicketEventPayload struct {
//	Status  string `json:"status"`
//	Type    string `json:"type"`
//	Content string `json:"content"`
//}
//
//type MessageEventPayload struct {
//	MessageID  uint64    `json:"message_id"`
//	SenderType string    `json:"sender_type"`
//	Content    string    `json:"text"`
//	CreatedAt  time.Time `json:"created_at"`
//}

type TicketMessageEvent struct {
	TicketID   uint64 `json:"ticket_id"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	Content    string `json:"content"`
}
