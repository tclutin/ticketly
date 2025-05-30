package ticket

import (
	"context"
	"errors"
	"fmt"
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

func (r *Repository) Create(ctx context.Context, ticket models.Ticket) (uint64, error) {
	sql := `INSERT INTO public.tickets (
		user_id,
		operator_id,
		status,
		type,
		sentiment,
		created_at,
		closed_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ticket_id`

	row := r.pool.QueryRow(
		ctx,
		sql,
		ticket.UserID,
		ticket.OperatorID,
		ticket.Status,
		ticket.Type,
		ticket.Sentiment,
		ticket.CreatedAt,
		ticket.ClosedAt,
	)

	var ticketId uint64
	if err := row.Scan(&ticketId); err != nil {
		return 0, err
	}

	return ticketId, nil
}

func (r *Repository) GetUserByTicketId(ctx context.Context, ticketId uint64) (models.User, error) {
	query := `SELECT u.user_id, u.external_id, u.username, u.source, u.is_banned, u.created_at FROM users AS u INNER JOIN tickets AS t ON u.user_id = t.user_id WHERE t.ticket_id = $1`

	row := r.pool.QueryRow(ctx, query, ticketId)

	var user models.User
	err := row.Scan(&user.UserID, &user.ExternalID, &user.Username, &user.Source, &user.IsBanned, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, coreerrors.ErrUserNotFound
		}
	}

	return user, nil
}

func (r *Repository) GetTicketById(ctx context.Context, ticketId uint64) (models.Ticket, error) {
	sql := `SELECT * FROM public.tickets WHERE ticket_id = $1`

	row := r.pool.QueryRow(ctx, sql, ticketId)

	var ticket models.Ticket
	err := row.Scan(
		&ticket.TicketID,
		&ticket.UserID,
		&ticket.OperatorID,
		&ticket.Status,
		&ticket.Type,
		&ticket.Sentiment,
		&ticket.CreatedAt,
		&ticket.ClosedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Ticket{}, coreerrors.ErrTicketNotFound
		}
		return ticket, err
	}

	return ticket, nil
}

// GetActiveTickets TODO: сделать потом универсальный метод
func (r *Repository) GetInProgressRealtimeTickets(ctx context.Context, operatorId uint64) ([]models.Ticket, error) {
	sql := `SELECT * FROM public.tickets WHERE status = 'in_progress' AND type = 'realtime-chat' AND operator_id = $1`

	rows, err := r.pool.Query(ctx, sql, operatorId)
	if err != nil {
		return nil, err
	}

	tickets, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Ticket])
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *Repository) GetAllWithFirstMessage(ctx context.Context) ([]models.PreviewTicket, error) {
	sql := `
			SELECT
				ticket_id,
				type,
				status,
				content as preview,
				sentiment,
				created_at,
				closed_at
			FROM
				public.tickets
			LEFT JOIN LATERAL (
			  	SELECT content
			  	FROM messages
			  	WHERE messages.ticket_id = tickets.ticket_id
			  	ORDER BY created_at ASC
			  	LIMIT 1
			) m ON true`

	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	tickets, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.PreviewTicket])
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

// Update TODO: лучше сделать обновление по конкретным полям, но времени не хватает
func (r *Repository) Update(ctx context.Context, ticketId uint64, model models.Ticket) error {
	sql := `UPDATE public.tickets SET
			operator_id = $1,
			status = $2,
			sentiment = $3,
			closed_at = $4
			WHERE ticket_id = $5`

	_, err := r.pool.Exec(
		ctx,
		sql,
		model.OperatorID,
		model.Status,
		model.Sentiment,
		model.ClosedAt,
		ticketId,
	)

	return err
}

func (r *Repository) HasActiveRealtimeTicket(ctx context.Context, userId uint64) (bool, error) {
	sql := `SELECT COUNT(*) FROM public.tickets WHERE user_id = $1 AND status IN ('open', 'in_progress') AND type = 'realtime-chat'`

	var count int
	err := r.pool.QueryRow(ctx, sql, userId).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to count active tickets: %w", err)
	}

	return count > 0, nil
}
