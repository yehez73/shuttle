-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS students (
	student_id BIGINT PRIMARY KEY,
	student_uuid UUID UNIQUE NOT NULL,
	parent_uuid UUID NOT NULL REFERENCES users (user_uuid) ON UPDATE NO ACTION ON DELETE SET NULL,
	school_uuid UUID NOT NULL REFERENCES schools (school_uuid) ON UPDATE NO ACTION ON DELETE SET NULL,
	student_first_name VARCHAR(255) NOT NULL,
	student_last_name VARCHAR(255) NOT NULL,
	student_gender VARCHAR(20) NOT NULL,
	student_grade VARCHAR(10) NOT NULL,
	student_address TEXT NULL DEFAULT NULL,
	student_pickup_point JSON NULL DEFAULT NULL,
	created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
	created_by VARCHAR(255) NULL DEFAULT NULL,
	updated_at TIMESTAMPTZ NULL DEFAULT NULL,
	updated_by VARCHAR(255) NULL DEFAULT NULL,
	deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
	deleted_by VARCHAR(255) NULL DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS students;
-- +goose StatementEnd
