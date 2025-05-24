package message

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/ticketly/ticketly_api/internal/models"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Create(ctx context.Context, msg models.Message) (uint64, error) {
	query := `
        INSERT INTO public.messages (ticket_id, sender_type, content, sentiment, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING message_id
    `

	row := r.pool.QueryRow(
		ctx,
		query,
		msg.TicketID,
		msg.SenderType,
		msg.Content,
		msg.Sentiment,
		msg.CreatedAt,
	)

	var messageId uint64
	if err := row.Scan(&messageId); err != nil {
		return 0, err
	}

	return messageId, nil
}

func (r *Repository) Get(ctx context.Context, msg models.Message) (uint64, error) {
	query := `
        INSERT INTO public.messages (ticket_id, sender_type, content, sentiment, created_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING message_id
    `

	row := r.pool.QueryRow(
		ctx,
		query,
		msg.TicketID,
		msg.SenderType,
		msg.Content,
		msg.Sentiment,
		msg.CreatedAt,
	)

	var messageId uint64
	if err := row.Scan(&messageId); err != nil {
		return 0, err
	}

	return messageId, nil
}
