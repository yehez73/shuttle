-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS routes (
    route_id BIGINT PRIMARY KEY,
    route_name_uuid UUID UNIQUE NOT NULL,
    school_uuid UUID NOT NULL,
    route_name VARCHAR(100) NOT NULL,
    route_description TEXT NULL DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMPTZ NULL DEFAULT NULL,
    updated_by VARCHAR(255) NULL DEFAULT NULL,
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
    deleted_by VARCHAR(255) NULL DEFAULT NULL
);

CREATE INDEX idx_route_name_uuid ON routes(route_name_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS routes CASCADE;
-- +goose StatementEnd
