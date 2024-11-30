package controllers

import (
	"shuttle/logger"
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddRoadRoute(c *fiber.Ctx) error {
	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("school_id").(string))
	if err != nil {
		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
			"school_id": c.Locals("school_id").(string),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	route := new(models.RoadRoute)
	if err := c.BodyParser(route); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, route); err != nil {
		return err
	}

	if err := services.AddRoadRoute(*route, SchoolObjID); err != nil {
		logger.LogError(err, "Failed to create route", map[string]interface{}{
			"school_id": SchoolObjID.Hex(),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Route created successfully", nil)
}