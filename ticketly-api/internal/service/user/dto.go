package user

type RegisterUserDTO struct {
	ExternalID string `json:"external_id"`
	Username   string `json:"username"`
	Source     string `json:"source"`
}
