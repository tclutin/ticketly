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

type TicketsDTO struct {
	Tickets []PreviewTicketDTO `json:"tickets"`
}

type PreviewTicketDTO struct {
	TicketID  uint64  `json:"ticket_id"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	Preview   string  `json:"preview"`
	Sentiment *string `json:"sentiment"`
}
