package operator

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Create(ctx context.Context) (uint64, error) {
	panic("implement me")
}

func (r *Repository) GetById(ctx context.Context) (uint64, error) {
	panic("implement me")
}

func (r *Repository) GetByCasdoorId(ctx context.Context) (uint64, error) {
	panic("implement me")
}
