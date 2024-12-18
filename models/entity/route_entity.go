package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type Route struct {
	RouteID          int64          `db:"route_id"`
	RouteUUID        uuid.UUID      `db:"route_uuid"`
	SchoolUUID       uuid.UUID      `db:"school_uuid"`
	RouteName        string         `db:"route_name"`
	RouteDescription string         `db:"route_description"`
	RouteStatus      string         `db:"route_status"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	CreatedBy        sql.NullString `db:"created_by"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
	UpdatedBy        sql.NullString `db:"updated_by"`
	DeletedAt        sql.NullTime   `db:"deleted_at"`
	DeletedBy        sql.NullString `db:"deleted_by"`
}

type RoutePoint struct {
	RoutePointID        int64          `db:"route_point_id"`
	RouteUUID           uuid.UUID      `db:"route_point_uuid"`
	RoutePointName      string         `db:"route_point_name"`
	RoutePointOrder     int            `db:"route_point_order"`
	RotePointLatitude   float64        `db:"route_point_latitude"`
	RoutePointLongitude float64        `db:"route_point_longitude"`
	CreatedAt           sql.NullTime   `db:"created_at"`
	CreatedBy           sql.NullString `db:"created_by"`
	UpdatedAt           sql.NullTime   `db:"updated_at"`
	UpdatedBy           sql.NullString `db:"updated_by"`
	DeletedAt           sql.NullTime   `db:"deleted_at"`
	DeletedBy           sql.NullString `db:"deleted_by"`
}