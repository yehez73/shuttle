-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vehicles (
	vehicle_id BIGINT PRIMARY KEY,
	vehicle_uuid UUID UNIQUE NOT NULL,
	school_uuid UUID NULL DEFAULT NULL,
	driver_uuid UUID NULL DEFAULT NULL,
	vehicle_name VARCHAR(50) NOT NULL,
	vehicle_number VARCHAR(20) NOT NULL,
	vehicle_type VARCHAR(20) NOT NULL,
	vehicle_color VARCHAR(20) NOT NULL,
	vehicle_seats INTEGER NOT NULL,
	vehicle_status VARCHAR(20) NULL DEFAULT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by VARCHAR(255) NULL DEFAULT NULL,
	updated_at TIMESTAMPTZ NULL DEFAULT NULL,
	updated_by VARCHAR(255) NULL DEFAULT NULL,
	deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
	deleted_by VARCHAR(255) NULL DEFAULT NULL
);

CREATE INDEX idx_vehicle_uuid ON vehicles(vehicle_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vehicles;
-- +goose StatementEnd
