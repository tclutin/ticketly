package models

import (
	"github.com/google/uuid"
	"time"
)

type Operator struct {
	OperatorID uint64
	CasdoorID  uuid.UUID
	Email      string
	Name       string
	LastSync   time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
