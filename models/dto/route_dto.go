package dto

type RoutePointRequestDTO struct {
	PointName string  `json:"point_name" validate:"required"`
	Order     int     `json:"point_order" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type RouteRequestDTO struct {
	RouteName        string                 `json:"route_name" validate:"required"`
	RouteDescription string                 `json:"route_description" validate:"required"`
	Points           []RoutePointRequestDTO `json:"points" validate:"required"`
	RouteStatus      string                 `json:"route_status" validate:"required"`
}
