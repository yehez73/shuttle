package repositories

import (
	"github.com/jmoiron/sqlx"
)

type RouteRepositoryInterface interface {
}

type RouteRepository struct {
	db *sqlx.DB
}

func NewRouteRepository(db *sqlx.DB) RouteRepositoryInterface {
	return &RouteRepository{
		db: db,
	}
}