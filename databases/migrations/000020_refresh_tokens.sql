-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS refresh_tokens (
	id BIGINT PRIMARY KEY,
	user_uuid UUID UNIQUE NOT NULL,
	refresh_token TEXT NOT NULL,
	issued_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
	expired_at TIMESTAMPTZ NOT NULL,
	is_revoked BOOLEAN NULL DEFAULT 'false',
	last_used_at TIMESTAMPTZ NULL DEFAULT NULL,
	FOREIGN KEY (user_uuid) REFERENCES users (user_uuid) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens CASCADE;
-- +goose StatementEnd