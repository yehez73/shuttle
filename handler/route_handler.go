package handler

import (
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RouteHandlerInterface interface {
	GetAllRoutesByAS(c *fiber.Ctx) error
	GetSpecRouteByAS(c *fiber.Ctx) error
	GetAllRoutesByDriver(c *fiber.Ctx) error
	GetSpecRouteByDriver(c *fiber.Ctx) error 
	AddRoute(c *fiber.Ctx) error
	UpdateRoute(c *fiber.Ctx) error
	DeleteRoute(c *fiber.Ctx) error
}

type routeHandler struct {
	routeService services.RouteServiceInterface
}

func NewRouteHttpHandler(routeService services.RouteServiceInterface) RouteHandlerInterface {
	return &routeHandler{
		routeService: routeService,
	}
}

func (handler *routeHandler) GetAllRoutesByAS(c *fiber.Ctx) error {
	schoolUUID, ok := c.Locals("schoolUUID").(string)
	if !ok {
		return utils.BadRequestResponse(c, "Invalid token or schoolUUID", nil)
	}

	routes, err := handler.routeService.GetAllRoutesByAS(schoolUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch routes",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"routes": routes,
	})
}

func (handler *routeHandler) GetSpecRouteByAS(c *fiber.Ctx) error {
	routeNameUUID := c.Params("id")

	driverUUID, err := handler.routeService.GetDriverUUIDByRouteName(routeNameUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get driver UUID",
		})
	}

	routeResponse, err := handler.routeService.GetSpecRouteByAS(routeNameUUID, driverUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(routeResponse)
}

func (handler *routeHandler) GetAllRoutesByDriver(c *fiber.Ctx) error {
	driverUUID, ok := c.Locals("userUUID").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token does not contain driver UUID",
		})
	}

	_, err := uuid.Parse(driverUUID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	routes, err := handler.routeService.GetAllRoutesByDriver(driverUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch routes",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"routes": routes,
	})
}

func (handler *routeHandler) GetSpecRouteByDriver(c *fiber.Ctx) error {
	driverUUID, ok := c.Locals("userUUID").(string)
	if !ok || driverUUID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	studentUUID := c.Params("id")
	if studentUUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student UUID is required",
		})
	}

	route, err := handler.routeService.GetSpecRouteByDriver(driverUUID, studentUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch specific route",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"route": route,
	})
}

func (handler *routeHandler) AddRoute(c *fiber.Ctx) error {
	schoolUUID, ok := c.Locals("schoolUUID").(string)
	if !ok {
		return utils.InternalServerErrorResponse(c, "Token does not contain schoolUUID", nil)
	}

	username, ok := c.Locals("user_name").(string)
	if !ok {
		return utils.InternalServerErrorResponse(c, "Token does not contain username", nil)
	}

	route := new(dto.RoutesRequestDTO)
	if err := c.BodyParser(route); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", nil)
	}

	if err := utils.ValidateStruct(c, route); err != nil {
		return utils.BadRequestResponse(c, err.Error(), nil)
	}

	err := handler.routeService.AddRoute(*route, schoolUUID, username)
	if err != nil {
		if err.Error() == "student not found" {
			return utils.BadRequestResponse(c, "Student not found", nil)
		}
		if err.Error() == "driver not found" {
			return utils.BadRequestResponse(c, "Driver not found", nil)
		}
		if err.Error() == "driver already assigned to another route" {
			return utils.BadRequestResponse(c, "Driver already assigned to another route", nil)
		}
		return utils.InternalServerErrorResponse(c, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Route added successfully", nil)
}

func (handler *routeHandler) UpdateRoute(c *fiber.Ctx) error {
	routenameUUID := c.Params("id")

	schoolUUID, ok := c.Locals("schoolUUID").(string)
	if !ok {
		return utils.InternalServerErrorResponse(c, "Token does not contain schoolUUID", nil)
	}

	username, ok := c.Locals("user_name").(string)
	if !ok {
		return utils.InternalServerErrorResponse(c, "Token does not contain username", nil)
	}

	route := new(dto.RoutesRequestDTO)
	if err := c.BodyParser(route); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", nil)
	}

	if err := utils.ValidateStruct(c, route); err != nil {
		return utils.BadRequestResponse(c, err.Error(), nil)
	}

	err := handler.routeService.UpdateRoute(*route, routenameUUID, schoolUUID, username)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Route updated successfully", nil)
}

func (handler *routeHandler) DeleteRoute(c *fiber.Ctx) error {
	routenameUUID := c.Params("id")

	schoolUUID, ok := c.Locals("schoolUUID").(string)
	if !ok {
		return utils.InternalServerErrorResponse(c, "Token does not contain schoolUUID", nil)
	}

	username, ok := c.Locals("user_name").(string)
	if !ok {
		return utils.InternalServerErrorResponse(c, "Token does not contain username", nil)
	}

	err := handler.routeService.DeleteRoute(routenameUUID, schoolUUID, username)
	if err != nil {
		if err.Error() == "route not found" {
			return utils.NotFoundResponse(c, "Route not found", nil)
		}
		return utils.InternalServerErrorResponse(c, err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Route deleted successfully", nil)
}