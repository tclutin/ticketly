package ticket

type CreateTicketDTO struct {
	UserID  uint64 `json:"user_id"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type CloseTicketDTO struct {
	TicketID   uint64 `json:"ticket_id"`
	OperatorID uint64 `json:"operator_id"`
	Message    string `json:"message"`
}
