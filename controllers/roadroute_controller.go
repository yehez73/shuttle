package controllers

import (
	"net/http"
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func AddRoadRoute(c *fiber.Ctx) error {
	token := string(c.Request().Header.Peek("Authorization"))
	UserID, err := utils.GetUserIDFromToken(token)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	SchoolID, err := services.CheckPermittedSchoolAccess(UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "You don't have permission to access this resource", nil)
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
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create route: "+err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Route created successfully", nil)
}