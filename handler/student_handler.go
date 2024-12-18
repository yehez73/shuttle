package handler

import (
	"fmt"
	"reflect"
	"shuttle/errors"
	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type StudentHandlerInterface interface {
	GetAllStudentWithParents(c *fiber.Ctx) error
	GetSpecStudentWithParents(c *fiber.Ctx) error
	AddSchoolStudentWithParents(c *fiber.Ctx) error
	UpdateSchoolStudentWithParents(c *fiber.Ctx) error
	DeleteSchoolStudentWithParentsIfNeccessary(c *fiber.Ctx) error
}

type studentHandler struct {
	studentService services.StudentService
}

func NewStudentHttpHandler(studentService services.StudentService) StudentHandlerInterface {
	return &studentHandler{
		studentService: studentService,
	}
}

func (handler *studentHandler) GetAllStudentWithParents(c *fiber.Ctx) error {
	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain school uuid", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return utils.BadRequestResponse(c, "Invalid page number", nil)
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		return utils.BadRequestResponse(c, "Invalid limit number", nil)
	}

	sortField := c.Query("sort_by", "student_id")
	sortDirection := c.Query("direction", "asc")

	if sortDirection != "asc" && sortDirection != "desc" {
		return utils.BadRequestResponse(c, "Invalid sort direction, use 'asc' or 'desc'", nil)
	}

	if !isValidSortFieldForStudents(sortField) {
		return utils.BadRequestResponse(c, "Invalid sort field", nil)
	}

	students, totalItems, err := handler.studentService.GetAllStudentsWithParents(page, limit, sortField, sortDirection, schoolUUIDStr)
	if err != nil {
		logger.LogError(err, "Failed to fetch paginated students", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	totalPages := (totalItems + limit - 1) / limit

	if page > totalPages {
		if totalItems > 0 {
			return utils.BadRequestResponse(c, "Page number out of range", nil)
		} else {
			page = 1
		}
	}

	start := (page-1)*limit + 1
	if totalItems == 0 || start > totalItems {
		start = 0
	}

	end := start + len(students) - 1
	if end > totalItems {
		end = totalItems
	}

	if len(students) == 0 {
		start = 0
		end = 0
	}

	response := fiber.Map{
		"data": students,
		"meta": fiber.Map{
			"current_page":   page,
			"total_pages":    totalPages,
			"per_page_items": limit,
			"total_items":    totalItems,
			"showing":        fmt.Sprintf("Showing %d-%d of %d", start, end, totalItems),
		},
	}

	return utils.SuccessResponse(c, "Students fetched successfully", response)
}

func (handler *studentHandler) GetSpecStudentWithParents(c *fiber.Ctx) error {
	id := c.Params("id")

	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain school uuid", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	students, err := handler.studentService.GetSpecStudentWithParents(id, schoolUUIDStr)
	if err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to fetch students", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Students fetched successfully", students)
}

func (handler *studentHandler) AddSchoolStudentWithParents(c *fiber.Ctx) error {
	username, ok := c.Locals("user_name").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain username", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain school uuid", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	student := new(dto.SchoolStudentParentRequestDTO)
	if err := c.BodyParser(student); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, student); err != nil {
		return utils.BadRequestResponse(c, err.Error(), nil)
	}

	if reflect.DeepEqual(dto.UserRequestsDTO{}, student.Parent) {
		return utils.BadRequestResponse(c, "Parent details are required", nil)
	}

	if err := handler.studentService.AddSchoolStudentWithParents(*student, schoolUUIDStr, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to add student", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Student created successfully", nil)
}

func (handler *studentHandler) UpdateSchoolStudentWithParents(c *fiber.Ctx) error {
	username, ok := c.Locals("user_name").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain username", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain school uuid", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	id := c.Params("id")

	student := new(dto.SchoolStudentParentRequestDTO)
	if err := c.BodyParser(student); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, student); err != nil {
		return utils.BadRequestResponse(c, err.Error(), nil)
	}

	if reflect.DeepEqual(dto.UserRequestsDTO{}, student.Parent) {
		return utils.BadRequestResponse(c, "Parent details are required", nil)
	}

	if err := handler.studentService.UpdateSchoolStudentWithParents(id, *student, schoolUUIDStr, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to update student", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Student updated successfully", nil)
}

func (handler *studentHandler) DeleteSchoolStudentWithParentsIfNeccessary(c *fiber.Ctx) error {
	id := c.Params("id")

	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain school uuid", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	username, ok := c.Locals("user_name").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain username", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	if err := handler.studentService.DeleteSchoolStudentWithParentsIfNeccessary(id, schoolUUIDStr, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to delete student", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Student deleted successfully", nil)
}

func isValidSortFieldForStudents(field string) bool {
	allowedFields := map[string]bool{
		"student_id":        true,
		"student_grade":     true,
		"student_first_name": true,
		"student_last_name":  true,
	}
	return allowedFields[field]
}