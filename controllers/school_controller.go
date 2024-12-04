package controllers

import (
	"net/mail"
	"regexp"

	"shuttle/errors"
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

	if err := validateCommonFields(c, school); err != nil {
		return utils.BadRequestResponse(c, err.Error(), nil)
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

	if err := validateCommonFields(c, school); err != nil {
		return utils.BadRequestResponse(c, err.Error(), nil)
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

func validateCommonFields(c *fiber.Ctx, school *models.School) error {
	if school.Name == "" {
		logger.LogError(nil, "School name is required", nil)
		return errors.New("School name is required", 0)
	}

	if school.Address == "" {
		logger.LogError(nil, "School address is required", nil)
		return errors.New("School address is required", 0)
	}

	if school.Contact == "" {
		logger.LogError(nil, "School contact number is required", nil)
		return errors.New("School contact number is required", 0)
	}

	phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	if !phoneRegex.MatchString(school.Contact) {
		logger.LogError(nil, "Invalid contact number format", nil)
		return errors.New("Invalid contact number format", 0)
	}

	if len(school.Contact) < 12 || len(school.Contact) > 15 {
		logger.LogError(nil, "Contact number should be between 12 to 15 characters", nil)
		return errors.New("Contact number should be between 12 to 15 characters", 0)
	}

	_, err := mail.ParseAddress(school.Email)
	if err != nil {
		logger.LogError(err, "Invalid email address", nil)
		return errors.New("Invalid email address", 0)
	}

	return nil
}