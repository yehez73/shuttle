package controllers

import (
	"log"
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetAllStudentWithParents(c *fiber.Ctx) error {
	UserID, ok := c.Locals("userId").(string)
	if !ok || UserID == "" {
		return utils.UnauthorizedResponse(c, "Invalid token", nil)
	}

	SchoolID, err := services.CheckPermittedSchoolAccess(UserID)
	if err != nil {
		return utils.UnauthorizedResponse(c, "You don't have permission to access this resource", nil)
	}

	students, err := services.GetAllPermitedSchoolStudentsWithParents(SchoolID)
	if err != nil {
		log.Println(err)
		return utils.InternalServerErrorResponse(c, "Failed to fetch students", nil)
	}

	return c.Status(fiber.StatusOK).JSON(students)
}

func AddSchoolStudentWithParents(c *fiber.Ctx) error {
	UserID, ok := c.Locals("userId").(string)
	if !ok || UserID == "" {
		return utils.UnauthorizedResponse(c, "Invalid token", nil)
	}

	username, ok := c.Locals("username").(string)
	if !ok || username == "" {
		return utils.UnauthorizedResponse(c, "Invalid token", nil)
	}

	SchoolID, err := services.CheckPermittedSchoolAccess(UserID)
	if err != nil {
		return utils.UnauthorizedResponse(c, "You don't have permission to access this resource", nil)
	}
	
	student := new(models.SchoolStudentRequest)
	if err := c.BodyParser(student); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	validate := validator.New()
	if err := validate.Struct(student); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return utils.BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
		}
	}

	if err := services.AddPermittedSchoolStudentWithParents(*student, SchoolID, username); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create student: "+err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Student created successfully", nil)
}

func UpdateSchoolStudentWithParents(c *fiber.Ctx) error {
	UserID, ok := c.Locals("userId").(string)
	if !ok || UserID == "" {
		return utils.UnauthorizedResponse(c, "Invalid token", nil)
	}

	SchoolID, err := services.CheckPermittedSchoolAccess(UserID)
	println("SchoolID", SchoolID.String())
	if err != nil {
		return utils.UnauthorizedResponse(c, "You don't have permission to access this resource", nil)
	}

	id := c.Params("id")
	student := new(models.SchoolStudentRequest)
	if err := c.BodyParser(student); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	validate := validator.New()
	if err := validate.Struct(student); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return utils.BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
		}
	}

	if err := services.UpdatePermittedSchoolStudentWithParents(id, *student, SchoolID); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update student: "+err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Student updated successfully", nil)
}

func DeleteSchoolStudentWithParents(c *fiber.Ctx) error {
	UserID, ok := c.Locals("userId").(string)
	if !ok || UserID == "" {
		return utils.UnauthorizedResponse(c, "Invalid token", nil)
	}

	SchoolID, err := services.CheckPermittedSchoolAccess(UserID)
	if err != nil {
		return utils.UnauthorizedResponse(c, "You don't have permission to access this resource", nil)
	}

	id := c.Params("id")
	if err := services.DeletePermittedSchoolStudentWithParents(id, SchoolID); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete student: "+err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Student deleted successfully", nil)
}