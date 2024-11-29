package controllers

import (
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func AddRoadRoute(c *fiber.Ctx) error {
	UserID, ok := c.Locals("userId").(string)
	if !ok || UserID == "" {
		return utils.UnauthorizedResponse(c, "Invalid Token", nil)
	}

	SchoolID, err := services.CheckPermittedSchoolAccess(UserID)
	if err != nil {
		return utils.UnauthorizedResponse(c, "You are not permitted to access this school", nil)
	}

	route := new(models.RoadRoute)
	if err := c.BodyParser(route); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	validate := validator.New()
	if err := validate.Struct(route); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return utils.BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
		}
	}

	if err := services.AddRoadRoute(*route, SchoolID); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to add route", nil)
	}

	return utils.SuccessResponse(c, "Route created successfully", nil)
}