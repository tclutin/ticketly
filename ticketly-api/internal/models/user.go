package models

import "time"

type User struct {
	UserID     uint64
	ExternalID string
	Username   string
	Source     string
	IsBanned   bool
	CreatedAt  time.Time
}
