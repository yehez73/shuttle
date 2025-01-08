-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS schools (
	school_id BIGINT PRIMARY KEY,
	school_uuid UUID UNIQUE NOT NULL,
	school_name VARCHAR(255) NOT NULL,
	school_address TEXT NOT NULL,
	school_contact VARCHAR(20) NOT NULL,
	school_email VARCHAR(255) NOT NULL,
	school_description TEXT NULL DEFAULT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by VARCHAR(255) NULL DEFAULT NULL,
	updated_at TIMESTAMPTZ NULL DEFAULT NULL,
	updated_by VARCHAR(255) NULL DEFAULT NULL,
	deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
	deleted_by VARCHAR(255) NULL DEFAULT NULL,
	school_point JSON NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_school_uuid ON schools(school_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS schools;
-- +goose StatementEnd
