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

type AssignTicketDTO struct {
	TicketID   uint64 `json:"ticket_id"`
	OperatorID uint64 `json:"operator_id"`
}

type AssignedTicketDTO struct {
	Channel           string `json:"channel"`
	ConnectionToken   string `json:"connection_token"`
	SubscriptionToken string `json:"subscription_token"`
}

type SendMessageDTO struct {
	TicketID   uint64 `json:"ticket_id"`
	OperatorID uint64 `json:"operator_id"`
	Message    string `json:"message"`
}

type ChannelInfo struct {
	Name              string `json:"name"`
	SubscriptionToken string `json:"subscription_token"`
}

type ConnectionsDTO struct {
	ConnectionToken string        `json:"connection_token"`
	Channels        []ChannelInfo `json:"channels"`
}
