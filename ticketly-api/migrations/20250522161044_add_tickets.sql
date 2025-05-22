-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.tickets (
    ticket_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    operator_id BIGINT,
    status VARCHAR(100) NOT NULL CHECK (status IN ('open', 'in_progress', 'closed')) DEFAULT 'open',
    type VARCHAR(100) NOT NULL CHECK (type IN ('only-message', 'realtime-chat')),
    sentiment VARCHAR(100) CHECK (sentiment IS NULL OR sentiment IN ('positive', 'neutral', 'negative')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES public.users(user_id)
);

CREATE TABLE IF NOT EXISTS public.messages (
    message_id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL,
    sender_type VARCHAR(20) CHECK (sender_type IN ('user', 'operator', 'system')),
    content TEXT NOT NULL,
    sentiment VARCHAR(100) CHECK (sentiment IN ('positive', 'neutral', 'negative')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.messages;
DROP TABLE IF EXISTS public.tickets;
-- +goose StatementEnd
