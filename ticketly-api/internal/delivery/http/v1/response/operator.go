package response

import (
	"github.com/google/uuid"
	"time"
)

type Operator struct {
	OperatorID uint64    `json:"operator_id"`
	CasdoorID  uuid.UUID `json:"casdoor_id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
