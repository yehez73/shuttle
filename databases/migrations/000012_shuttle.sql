-- +goose Up
-- +goose StatementBegin
CREATE TYPE shuttle_status AS ENUM (
    'home',
    'waiting_to_be_taken_to_school',
    'going_to_school',
    'at_school',
    'waiting_to_be_taken_to_home',
    'going_to_home'
);

ALTER TYPE shuttle_status
    OWNER TO postgres;

-- Create the shuttle table with shuttle_status ENUM type for the status column
CREATE TABLE IF NOT EXISTS shuttle (
    shuttle_id BIGINT PRIMARY KEY,
    shuttle_uuid UUID UNIQUE NOT NULL,
    student_uuid UUID NOT NULL,
    driver_uuid UUID NOT NULL,
    status shuttle_status NOT NULL DEFAULT 'home',
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NULL DEFAULT NULL,
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
    FOREIGN KEY (driver_uuid) REFERENCES users (user_uuid) ON UPDATE NO ACTION ON DELETE SET NULL,
    FOREIGN KEY (student_uuid) REFERENCES students (student_uuid) ON UPDATE NO ACTION ON DELETE SET NULL
);

CREATE INDEX idx_shuttle_uuid ON shuttle(shuttle_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shuttle;
DROP TYPE IF EXISTS shuttle_status;
-- +goose StatementEnd
