package operator

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/ticketly/ticketly_api/internal/models"
	coreerrors "github.com/tclutin/ticketly/ticketly_api/internal/service/errors"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Create(ctx context.Context, operator models.Operator) (uint64, error) {
	sql := `INSERT INTO public.operators (casdoor_id, email, name, last_sync, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING operator_id;`

	row := r.pool.QueryRow(
		ctx,
		sql,
		operator.CasdoorID,
		operator.Email,
		operator.Name,
		operator.LastSync,
		operator.CreatedAt,
		operator.UpdatedAt,
	)

	var operatorId uint64
	if err := row.Scan(&operatorId); err != nil {
		return 0, err
	}

	return operatorId, nil
}

func (r *Repository) Update(ctx context.Context, operator models.Operator) error {
	sql := `UPDATE public.operators SET email = $1, name = $2, last_sync = $3, updated_at = $4 WHERE casdoor_id = $5;`

	_, err := r.pool.Exec(
		ctx,
		sql,
		operator.Email,
		operator.Name,
		operator.LastSync,
		operator.UpdatedAt,
		operator.CasdoorID,
	)

	return err
}

func (r *Repository) GetById(ctx context.Context, operatorId uint64) (models.Operator, error) {
	sql := `SELECT * FROM operators WHERE operator_id = $1`

	row := r.pool.QueryRow(ctx, sql, operatorId)

	var operator models.Operator
	err := row.Scan(
		&operator.OperatorID,
		&operator.CasdoorID,
		&operator.Email,
		&operator.Name,
		&operator.LastSync,
		&operator.CreatedAt,
		&operator.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return operator, coreerrors.ErrOperatorNotFound
		}
		return operator, err
	}

	return operator, nil
}

func (r *Repository) GetByCasdoorId(ctx context.Context, casdoorId uuid.UUID) (models.Operator, error) {
	sql := `SELECT * FROM operators WHERE casdoor_id = $1`

	row := r.pool.QueryRow(ctx, sql, casdoorId)

	var operator models.Operator
	err := row.Scan(
		&operator.OperatorID,
		&operator.CasdoorID,
		&operator.Email,
		&operator.Name,
		&operator.LastSync,
		&operator.CreatedAt,
		&operator.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return operator, coreerrors.ErrOperatorNotFound
		}
		return operator, err
	}

	return operator, nil
}
