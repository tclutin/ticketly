package user

import (
	"context"
	"github.com/tclutin/ticketly/ticketly_api/internal/models"
	"github.com/tclutin/ticketly/ticketly_api/internal/repository"
	coreerrors "github.com/tclutin/ticketly/ticketly_api/internal/service/errors"
	"time"
)

type Service struct {
	repo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) *Service {
	return &Service{
		repo: userRepo,
	}
}

func (s *Service) Create(ctx context.Context, dto RegisterUserDTO) (uint64, error) {
	if _, err := s.GetByExternalId(ctx, dto.ExternalID); err == nil {
		return 0, coreerrors.ErrUserAlreadyExists
	}

	user := models.User{
		ExternalID: dto.ExternalID,
		Username:   dto.Username,
		Source:     dto.Source,
		IsBanned:   false,
		CreatedAt:  time.Now().UTC(),
	}

	return s.repo.Create(ctx, user)
}

func (s *Service) GetByExternalId(ctx context.Context, externalId string) (user models.User, err error) {
	return s.repo.GeByExternalId(ctx, externalId)
}

func (s *Service) GetById(ctx context.Context, userId uint64) (user models.User, err error) {
	return s.repo.GetById(ctx, userId)
}
