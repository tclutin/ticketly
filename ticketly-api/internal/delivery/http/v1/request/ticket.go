package request

type CreateTicketRequest struct {
	UserID  uint64 `json:"user_id" binding:"required"`
	Type    string `json:"type" binding:"required,oneof=only-message realtime-chat"`
	Content string `json:"content" binding:"required,min=6,max=500"`
}

type CloseTicketRequest struct {
	Content string `json:"content" binding:"required,min=6,max=500"`
}
