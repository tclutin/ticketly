package response

import "time"

type User struct {
	UserID     uint64    `json:"user_id"`
	ExternalID string    `json:"external_id"`
	Username   string    `json:"username"`
	Source     string    `json:"source"`
	IsBanned   bool      `json:"is_banned"`
	CreatedAt  time.Time `json:"created_at"`
}
