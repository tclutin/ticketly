package operator

import "github.com/google/uuid"

type CreateOperatorDTO struct {
	CasdooID uuid.UUID
	Email    string
	Name     string
}
