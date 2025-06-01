package operator

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/tclutin/ticketly/ticketly_api/internal/models"
	"github.com/tclutin/ticketly/ticketly_api/internal/repository"
	coreerrors "github.com/tclutin/ticketly/ticketly_api/internal/service/errors"
	"time"
)

const syncThreshold = 24 * time.Hour

type Service struct {
	repo repository.OperatorRepository
}

func NewService(repo repository.OperatorRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) SyncOperator(ctx context.Context, dto CreateOperatorDTO) (uint64, error) {
	operator, err := s.repo.GetByCasdoorId(ctx, dto.CasdooID)
	if err != nil {
		if errors.Is(err, coreerrors.ErrOperatorNotFound) {
			return s.repo.Create(ctx, models.Operator{
				CasdoorID: dto.CasdooID,
				Email:     dto.Email,
				Name:      dto.Name,
				LastSync:  time.Now().UTC(),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			})
		}
		return 0, fmt.Errorf("failed to get operator: %w", err)
	}

	//но если пользователь будет удалён в casdoor

	if time.Since(operator.LastSync) > syncThreshold {
		operator.Email = dto.Email
		operator.Name = dto.Name
		operator.LastSync = time.Now().UTC()
		if err = s.repo.Update(ctx, operator); err != nil {
			return 0, fmt.Errorf("failed to update operator: %w", err)
		}
	}

	return operator.OperatorID, nil
}

func (s *Service) GetByCasdoorId(ctx context.Context, casdoorId uuid.UUID) (models.Operator, error) {
	return s.repo.GetByCasdoorId(ctx, casdoorId)
}

func (s *Service) GetById(ctx context.Context, operatorId uint64) (models.Operator, error) {
	return s.repo.GetById(ctx, operatorId)
}
