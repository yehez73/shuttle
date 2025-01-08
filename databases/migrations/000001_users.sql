-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
	user_id BIGINT PRIMARY KEY,
	user_uuid UUID UNIQUE NOT NULL,
	user_username VARCHAR(255) NOT NULL,
	user_email VARCHAR(255) NOT NULL,
	user_password VARCHAR(255) NOT NULL,
	user_role VARCHAR(20) NOT NULL,
	user_role_code VARCHAR(5) NULL DEFAULT NULL,
	user_status VARCHAR(20) DEFAULT 'offline',
	user_last_active TIMESTAMPTZ NULL DEFAULT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by VARCHAR(255),
	updated_at TIMESTAMPTZ NULL DEFAULT NULL,
	updated_by VARCHAR(255) NULL DEFAULT NULL,
	deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
	deleted_by VARCHAR(255) NULL DEFAULT NULL
);

CREATE INDEX idx_user_uuid ON users(user_uuid);
CREATE INDEX idx_user_username ON users(user_username);
CREATE INDEX idx_user_email ON users(user_email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd