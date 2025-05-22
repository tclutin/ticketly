package user

import (
	"context"
	"errors"
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

func (r *Repository) Create(ctx context.Context, user models.User) (uint64, error) {
	sql := `INSERT INTO public.users (
                          external_id,
                          username,
                          source,
                          is_banned,
                          created_at
            ) VALUES ($1, $2, $3, $4, $5) RETURNING user_id`

	row := r.pool.QueryRow(
		ctx,
		sql,
		user.ExternalID,
		user.Username,
		user.Source,
		user.IsBanned,
		user.CreatedAt,
	)

	var userId uint64
	if err := row.Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}

func (r *Repository) GeByExternalId(ctx context.Context, externalId string) (models.User, error) {
	sql := `SELECT * FROM public.users WHERE external_id = $1`

	row := r.pool.QueryRow(ctx, sql, externalId)

	var usr models.User
	err := row.Scan(
		&usr.UserID,
		&usr.ExternalID,
		&usr.Username,
		&usr.Source,
		&usr.IsBanned,
		&usr.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return usr, coreerrors.ErrUserNotFound
		}
		return usr, err
	}

	return usr, nil
}

func (r *Repository) GetById(ctx context.Context, userId uint64) (models.User, error) {
	sql := `SELECT * FROM public.users WHERE user_id = $1`

	row := r.pool.QueryRow(ctx, sql, userId)

	var usr models.User
	err := row.Scan(
		&usr.UserID,
		&usr.ExternalID,
		&usr.Username,
		&usr.Source,
		&usr.IsBanned,
		&usr.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return usr, coreerrors.ErrUserNotFound
		}
		return usr, err
	}

	return usr, nil
}
