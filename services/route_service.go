package services

import (
	"shuttle/errors"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"time"

	"github.com/google/uuid"
)

type RouteServiceInterface interface {
	AddRoute(route dto.RouteRequestDTO, schoolUUID, username string) error
}

type RouteService struct {
	routeRepository repositories.RouteRepositoryInterface
}

func NewRouteService(routeRepository repositories.RouteRepositoryInterface) RouteServiceInterface {
	return &RouteService{
		routeRepository: routeRepository,
	}
}

func (service *RouteService) AddRoute(route dto.RouteRequestDTO, schoolUUID, username string) error {
	routeEntity := entity.Route{
		RouteID:        time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		RouteUUID:      uuid.New(),
	}

	if err := service.routeRepository.AddRoute(routeEntity); err != nil {
		return errors.NewInternalServerErr("Failed to add route")
	}

	return nil
}