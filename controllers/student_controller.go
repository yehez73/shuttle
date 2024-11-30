package controllers

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

func GetAllStudentWithParents(c *fiber.Ctx) error {
	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("school_id").(string))
	if err != nil {
		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
			"school_id": c.Locals("school_id"),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	students, err := services.GetAllPermitedSchoolStudentsWithParents(SchoolObjID)
	if err != nil {
		logger.LogError(err, "Failed to fetch all students", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(students)
}

func AddSchoolStudentWithParents(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("school_id").(string))
	if err != nil {
		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
			"school_id": c.Locals("school_id"),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}
	
	student := new(models.SchoolStudentRequest)
	if err := c.BodyParser(student); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, student); err != nil {
		return err
	}

	if (models.User{}) == student.Parent {
		return utils.BadRequestResponse(c, "Parent details are required", nil)
	}

	if student.Parent.Phone == "" || student.Parent.Address == "" || student.Parent.Email == "" {
		return utils.BadRequestResponse(c, "Parent details are required", nil)
	}

	if err := services.AddPermittedSchoolStudentWithParents(*student, SchoolObjID, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to add student", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Student created successfully", nil)
}

func UpdateSchoolStudentWithParents(c *fiber.Ctx) error {
	id := c.Params("id")
	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("school_id").(string))
	if err != nil {
		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
			"school_id": c.Locals("school_id"),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	student := new(models.SchoolStudentRequest)
	if err := c.BodyParser(student); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, student); err != nil {
		return err
	}

	if err := services.UpdatePermittedSchoolStudentWithParents(id, *student, SchoolObjID); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to update student", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Student updated successfully", nil)
}

func DeleteSchoolStudentWithParents(c *fiber.Ctx) error {
	id := c.Params("id")
	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("school_id").(string))
	if err != nil {
		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
			"school_id": c.Locals("school_id"),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	if err := services.DeletePermittedSchoolStudentWithParents(id, SchoolObjID); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to delete student", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Student deleted successfully", nil)
}