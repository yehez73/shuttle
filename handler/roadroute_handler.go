package handler

import (
	"shuttle/errors"
	"shuttle/logger"
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllRoutes(c *fiber.Ctx) error {
	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
	if err != nil {
		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
			"school_id": c.Locals("schoolId").(string),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	routes, err := services.GetAllRoutes(SchoolObjID)
	if err != nil {
		logger.LogError(err, "Failed to fetch routes", map[string]interface{}{
			"school_id": SchoolObjID.Hex(),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(routes)
}

func GetSpecRoute(c *fiber.Ctx) error {
	RouteObjID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		logger.LogError(err, "Failed to convert route id", map[string]interface{}{
			"route_id": c.Params("id"),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	route, err := services.GetSpecRoute(RouteObjID)
	if err != nil {
		logger.LogError(err, "Failed to fetch route", map[string]interface{}{
			"route_id": RouteObjID.Hex(),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(route)
}

func AddRoute(c *fiber.Ctx) error {
	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
	if err != nil {
		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
			"school_id": c.Locals("schoolId").(string),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	username := c.Locals("username").(string)

	route := new(models.RoadRoute)
	if err := c.BodyParser(route); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, route); err != nil {
		return err
	}

	if err := services.AddRoute(*route, SchoolObjID, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to add route", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Route created successfully", nil)
}