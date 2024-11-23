package controllers

import (
	"log"
	"net/http"
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetAllVehicles (c *fiber.Ctx) error {
	vehicles, err := services.GetAllVehicles()
	if err != nil {
		log.Println(err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}

	return c.Status(fiber.StatusOK).JSON(vehicles)
}

func GetSpecVehicle (c *fiber.Ctx) error {
	id := c.Params("id")
	vehicle, err := services.GetSpecVehicle(id)
	if err != nil {
		log.Println(err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}

	return c.Status(fiber.StatusOK).JSON(vehicle)
}

func AddVehicle (c *fiber.Ctx) error {
	vehicle := new(models.Vehicle)
	if err := c.BodyParser(vehicle); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	validate := validator.New()
	if err := validate.Struct(vehicle); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return utils.BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
		}
	}

	if err := services.AddVehicle(*vehicle); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create vehicle: "+err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Vehicle created successfully", nil)
}

func UpdateVehicle (c *fiber.Ctx) error {
	id := c.Params("id")
	vehicle := new(models.Vehicle)
	if err := c.BodyParser(vehicle); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	validate := validator.New()
	if err := validate.Struct(vehicle); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return utils.BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
		}
	}

	if err := services.UpdateVehicle(*vehicle, id); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update vehicle: "+err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Vehicle updated successfully", nil)
}

func DeleteVehicle (c *fiber.Ctx) error {
	id := c.Params("id")
	if err := services.DeleteVehicle(id); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete vehicle: "+err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Vehicle deleted successfully", nil)
}