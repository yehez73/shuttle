package controllers

import (
	"shuttle/logger"
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllSchools(c *fiber.Ctx) error {
	schools, err := services.GetAllSchools()
	if err != nil {
		logger.LogError(err, "Failed to fetch all schools", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(schools)
}

func GetSpecSchool(c *fiber.Ctx) error {
	id := c.Params("id")
	
	school, err := services.GetSpecSchool(id)
	if err != nil {
		logger.LogError(err, "Failed to fetch specific school", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(school)
}

func AddSchool(c *fiber.Ctx) error {
	school := new(models.School)
	if err := c.BodyParser(school); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, school); err != nil {
		return err
	}

	if err := services.AddSchool(*school); err != nil {
		logger.LogError(err, "Failed to create school", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "School created successfully", nil)
}

func UpdateSchool(c *fiber.Ctx) error {
	id := c.Params("id")
	school := new(models.School)
	if err := c.BodyParser(school); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, school); err != nil {
		return err
	}

	if err := services.UpdateSchool(id, *school); err != nil {
		logger.LogError(err, "Failed to update school", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "School updated successfully", nil)
}

func DeleteSchool(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := services.DeleteSchool(id); err != nil {
		logger.LogError(err, "Failed to delete school", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "School deleted successfully", nil)
}