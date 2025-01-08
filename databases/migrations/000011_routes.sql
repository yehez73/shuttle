-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "route_assignment" (
	"route_id" BIGINT PRIMARY KEY,
	"route_uuid" UUID UNIQUE NOT NULL,
	"driver_uuid" UUID NOT NULL,
	"student_uuid" UUID NOT NULL,
	"route_name" VARCHAR(100) NOT NULL,
	"route_description" TEXT NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"created_by" VARCHAR(255) NOT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"updated_by" VARCHAR(255) NULL DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"deleted_by" VARCHAR(255) NULL DEFAULT NULL,
	"school_uuid" UUID NOT NULL,
	FOREIGN KEY ("driver_uuid") REFERENCES "driver_details" ("user_uuid") ON UPDATE NO ACTION ON DELETE NO ACTION,
	FOREIGN KEY ("driver_uuid") REFERENCES "users" ("user_uuid") ON UPDATE NO ACTION ON DELETE NO ACTION,
	FOREIGN KEY ("driver_uuid") REFERENCES "driver_details" ("user_uuid") ON UPDATE NO ACTION ON DELETE SET NULL,
	FOREIGN KEY ("driver_uuid") REFERENCES "users" ("user_uuid") ON UPDATE NO ACTION ON DELETE SET NULL,
	FOREIGN KEY ("school_uuid") REFERENCES "schools" ("school_uuid") ON UPDATE NO ACTION ON DELETE NO ACTION,
	FOREIGN KEY ("student_uuid") REFERENCES "students" ("student_uuid") ON UPDATE NO ACTION ON DELETE NO ACTION
);

CREATE INDEX idx_route_uuid ON route_assignment(route_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS route_assignment;
-- +goose StatementEnd
