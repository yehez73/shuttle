package controllers

import (
	"net/http"
	"shuttle/logger"
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllVehicles (c *fiber.Ctx) error {
	vehicles, err := services.GetAllVehicles()
	if err != nil {
		logger.LogError(err, "Failed to fetch all vehicles", nil)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(vehicles)
}

func GetSpecVehicle (c *fiber.Ctx) error {
	id := c.Params("id")
	vehicle, err := services.GetSpecVehicle(id)
	if err != nil {
		logger.LogError(err, "Failed to fetch specific vehicle", nil)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(vehicle)
}

func AddVehicle (c *fiber.Ctx) error {
	vehicle := new(models.Vehicle)
	if err := c.BodyParser(vehicle); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, vehicle); err != nil {
		return err
	}

	if err := services.AddVehicle(*vehicle); err != nil {
		logger.LogError(err, "Failed to create vehicle", nil)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Vehicle created successfully", nil)
}

func UpdateVehicle (c *fiber.Ctx) error {
	id := c.Params("id")
	vehicle := new(models.Vehicle)
	if err := c.BodyParser(vehicle); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, vehicle); err != nil {
		return err
	}

	if err := services.UpdateVehicle(*vehicle, id); err != nil {
		logger.LogError(err, "Failed to update vehicle", nil)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Vehicle updated successfully", nil)
}

func DeleteVehicle (c *fiber.Ctx) error {
	id := c.Params("id")
	if err := services.DeleteVehicle(id); err != nil {
		logger.LogError(err, "Failed to delete vehicle", nil)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Vehicle deleted successfully", nil)
}