package request

type RegisterUserRequest struct {
	ExternalID string `json:"external_id" binding:"required"`
	Username   string `json:"username" binding:"required,min=4,max=40"`
	Source     string `json:"source" binding:"required,min=3,max=40,oneof=telegram"`
}
