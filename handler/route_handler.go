package handler

import (
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
)

type RouteHandlerInterface interface {
	AddRoute(c *fiber.Ctx) error
}

type routeHandler struct {
	routeService services.RouteService
}

func NewRouteHttpHandler(routeService services.RouteService) RouteHandlerInterface {
	return &routeHandler{
		routeService: routeService,
	}
}

func (handler *routeHandler) AddRoute(c *fiber.Ctx) error {
	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		return utils.InternalServerErrorResponse(c, "Token does not contain school uuid", nil)
	}

	username, ok := c.Locals("user_name").(string)
	if !ok {
		return utils.InternalServerErrorResponse(c, "Token does not contain username", nil)
	}

	route := new(dto.RouteRequestDTO)
	if err := c.BodyParser(route); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", nil)
	}

	if err := utils.ValidateStruct(c, route); err != nil {
		return utils.BadRequestResponse(c, err.Error(), nil)
	}

	if err := handler.routeService.AddRoute(*route, schoolUUIDStr, username); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Route added successfully", nil)
}