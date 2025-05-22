package repository

import (
	"context"
	"github.com/tclutin/ticketly/ticketly_api/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	GeByExternalId(ctx context.Context, externalId string) (models.User, error)
}
