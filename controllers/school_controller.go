package controllers

import (
	"log"
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetAllSchools(c *fiber.Ctx) error {
	schools, err := services.GetAllSchools()
	if err != nil {
		log.Print(err)
		return utils.InternalServerErrorResponse(c, "Internal server error", nil)
	}

	return c.Status(fiber.StatusOK).JSON(schools)
}

func GetSpecSchool(c *fiber.Ctx) error {
	id := c.Params("id")
	
	school, err := services.GetSpecSchool(id)
	if err != nil {
		log.Print(err)
		return utils.InternalServerErrorResponse(c, "Internal server error", nil)
	}

	return c.Status(fiber.StatusOK).JSON(school)
}

func AddSchool(c *fiber.Ctx) error {
	school := new(models.School)
	if err := c.BodyParser(school); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	validate := validator.New()
	if err := validate.Struct(school); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return utils.BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
		}
	}	

	if err := services.AddSchool(*school); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create school" + err.Error(), nil)
	}

	return utils.SuccessResponse(c, "School created successfully", nil)
}

func UpdateSchool(c *fiber.Ctx) error {
	id := c.Params("id")
	school := new(models.School)
	if err := c.BodyParser(school); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	validate := validator.New()
	if err := validate.Struct(school); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return utils.BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
		}
	}	

	if err := services.UpdateSchool(id, *school); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update school: " + err.Error(), nil)
	}

	return utils.SuccessResponse(c, "School updated successfully", nil)
}

func DeleteSchool(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := services.DeleteSchool(id); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete school: " + err.Error(), nil)
	}

	return utils.SuccessResponse(c, "School deleted successfully", nil)
}