package message

import (
	"context"
	"github.com/jackc/pgx/v5"
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

// нужно ли передавать id сообщения
func (r *Repository) Update(ctx context.Context, messageId uint64, msg models.Message) error {
	sql := `
        UPDATE public.messages
        SET 
            sender_type = $1,
            content = $2,
            sentiment = $3,
            created_at = $4
        WHERE message_id = $5
    `

	_, err := r.pool.Exec(
		ctx,
		sql,
		msg.SenderType,
		msg.Content,
		msg.Sentiment,
		msg.CreatedAt,
		messageId,
	)

	return err
}

func (r *Repository) GetAll(ctx context.Context, ticketId uint64) ([]models.MessagePreview, error) {
	query := `
		   SELECT
				message_id,
				ticket_id,
				sender_type,
				content,
				sentiment,
				created_at
		   FROM
		   		public.messages
		   WHERE
		       ticket_id = $1
    `

	rows, err := r.pool.Query(ctx, query, ticketId)
	if err != nil {
		return nil, err
	}

	messages, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.MessagePreview])
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *Repository) GetById(ctx context.Context, messageId uint64) (models.Message, error) {
	query := `
		   SELECT
				message_id,
				ticket_id,
				sender_type,
				content,
				sentiment,
				created_at
		   FROM
		   		public.messages
		   WHERE
		       message_id = $1
    `

	row := r.pool.QueryRow(ctx, query, messageId)

	var message models.Message

	err := row.Scan(
		&message.MessageID,
		&message.TicketID,
		&message.SenderType,
		&message.Content,
		&message.Sentiment,
		&message.CreatedAt)

	if err != nil {
		return message, err
	}

	return message, nil
}
