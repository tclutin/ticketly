package ticket

type CreateTicketDTO struct {
	UserID  uint64
	Type    string
	Content string
}

type CloseTicketDTO struct {
	TicketID   uint64
	OperatorID uint64
	Message    string
}
