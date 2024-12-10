package handler

import (
	"encoding/json"
	"fmt"
	"strconv"

	"shuttle/errors"
	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UserHandlerInterface interface {
	GetAllSuperAdmin(c *fiber.Ctx) error
	GetAllSchoolAdmin(c *fiber.Ctx) error
	GetAllPermittedDriver(c *fiber.Ctx) error

	GetSpecSuperAdmin(c *fiber.Ctx) error
	GetSpecSchoolAdmin(c *fiber.Ctx) error
	GetSpecPermittedDriver(c *fiber.Ctx) error

	AddUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error

	DeleteSuperAdmin(c *fiber.Ctx) error
	DeleteSchoolAdmin(c *fiber.Ctx) error
	DeleteDriver(c *fiber.Ctx) error
}

type userHandler struct {
	userService   services.UserService
	schoolService services.SchoolService
}

func NewUserHttpHandler(userService services.UserService, schoolService services.SchoolService) UserHandlerInterface {
	return &userHandler{
		userService:   userService,
		schoolService: schoolService,
	}
}

func (handler *userHandler) GetAllSuperAdmin(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return utils.BadRequestResponse(c, "Invalid page number", nil)
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		return utils.BadRequestResponse(c, "Invalid limit number", nil)
	}

	sortField := c.Query("sort_by", "user_id")
	sortDirection := c.Query("direction", "desc")

	if sortDirection != "asc" && sortDirection != "desc" {
		return utils.BadRequestResponse(c, "Invalid sort direction, use 'asc' or 'desc'", nil)
	}

	if !isValidSortFieldForUsers(sortField) {
		return utils.BadRequestResponse(c, "Invalid sort field", nil)
	}

	users, totalItems, err := handler.userService.GetAllSuperAdmin(page, limit, sortField, sortDirection)
	if err != nil {
		logger.LogError(err, "Failed to fetch paginated super admins", nil)
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

	end := start + len(users) - 1
	if end > totalItems {
		end = totalItems
	}

	if len(users) == 0 {
		start = 0
		end = 0
	}

	response := fiber.Map{
		"data": users,
		"meta": fiber.Map{
			"current_page":   page,
			"total_pages":    totalPages,
			"per_page_items": limit,
			"total_items":    totalItems,
			"showing":        fmt.Sprintf("Showing %d-%d of %d", start, end, totalItems),
		},
	}

	return utils.SuccessResponse(c, "Users fetched successfully", response)
}

func (handler *userHandler) GetAllSchoolAdmin(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return utils.BadRequestResponse(c, "Invalid page number", nil)
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		return utils.BadRequestResponse(c, "Invalid limit number", nil)
	}

	sortField := c.Query("sort_by", "user_id")
	sortDirection := c.Query("direction", "desc")

	if sortDirection != "asc" && sortDirection != "desc" {
		return utils.BadRequestResponse(c, "Invalid sort direction, use 'asc' or 'desc'", nil)
	}

	if !isValidSortFieldForUsers(sortField) {
		return utils.BadRequestResponse(c, "Invalid sort field", nil)
	}

	users, totalItems, err := handler.userService.GetAllSchoolAdmin(page, limit, sortField, sortDirection)
	if err != nil {
		logger.LogError(err, "Failed to fetch paginated school admins", nil)
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

	end := start + len(users) - 1
	if end > totalItems {
		end = totalItems
	}

	if len(users) == 0 {
		start = 0
		end = 0
	}

	response := fiber.Map{
		"data": users,
		"meta": fiber.Map{
			"current_page":   page,
			"total_pages":    totalPages,
			"per_page_items": limit,
			"total_items":    totalItems,
			"showing":        fmt.Sprintf("Showing %d-%d of %d", start, end, totalItems),
		},
	}

	return utils.SuccessResponse(c, "Users fetched successfully", response)
}

func (handler *userHandler) GetAllPermittedDriver(c *fiber.Ctx) error {
	role := c.Locals("role_code").(string)

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return utils.BadRequestResponse(c, "Invalid page number", nil)
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		return utils.BadRequestResponse(c, "Invalid limit number", nil)
	}

	sortField := c.Query("sort_by", "user_id")
	sortDirection := c.Query("direction", "desc")

	if sortDirection != "asc" && sortDirection != "desc" {
		return utils.BadRequestResponse(c, "Invalid sort direction, use 'asc' or 'desc'", nil)
	}

	if !isValidSortFieldForUsers(sortField) {
		return utils.BadRequestResponse(c, "Invalid sort field", nil)
	}

	switch role {
	case "SA":
		users, totalItems, err := handler.userService.GetAllDriverFromAllSchools(page, limit, sortField, sortDirection)
		if err != nil {
			logger.LogError(err, "Failed to fetch all drivers", nil)
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

		end := start + len(users) - 1
		if end > totalItems {
			end = totalItems
		}

		if len(users) == 0 {
			start = 0
			end = 0
		}

		response := fiber.Map{
			"data": users,
			"meta": fiber.Map{
				"current_page":   page,
				"total_pages":    totalPages,
				"per_page_items": limit,
				"total_items":    totalItems,
				"showing":        fmt.Sprintf("Showing %d-%d of %d", start, end, totalItems),
			},
		}

		return utils.SuccessResponse(c, "Users fetched successfully", response)
	case "AS":
		schoolUUID, ok := c.Locals("schoolUUID").(string)
		if !ok {
			return utils.BadRequestResponse(c, "Token is invalid", nil)
		}

		users, totalItems, err := handler.userService.GetAllDriverForPermittedSchool(page, limit, sortField, sortDirection, schoolUUID)
		if err != nil {
			logger.LogError(err, "Failed to fetch all drivers", nil)
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

		end := limit

		if len(users) == 0 {
			start = 0
			end = 0
		}

		response := fiber.Map{
			"data": users,
			"meta": fiber.Map{
				"current_page":   page,
				"total_pages":    totalPages,
				"per_page_items": limit,
				"total_items":    totalItems,
				"showing":        fmt.Sprintf("Showing %d-%d of %d", start, end, totalItems),
			},
		}

		return utils.SuccessResponse(c, "Users fetched successfully", response)
	default:
		return utils.BadRequestResponse(c, "Invalid role", nil)
	}
}

func (handler *userHandler) GetSpecPermittedDriver(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role_code").(string)

	var user dto.UserResponseDTO
	var err error

	switch role {
	case "SA":
		user, err = handler.userService.GetSpecDriverFromAllSchools(id)
	// case "AS":
	// 	schoolUUID, ok := c.Locals("schoolUUID").(string)
	// 	if !ok {
	// 		return utils.BadRequestResponse(c, "Token is invalid", nil)
	// 	}

	// 	user, err = handler.userService.GetSpecDriverForPermittedSchool(id, schoolUUID)
	default:
		return utils.BadRequestResponse(c, "Invalid role", nil)
	}

	if err != nil {
		logger.LogError(err, "Failed to fetch specific driver", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "User fetched successfully", user)
}

func (handler *userHandler) DeleteDriver(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Locals("user_name").(string)

	forceDelete := c.Query("force_delete")

	existingUser, checkErr := handler.userService.GetSpecUserWithDetails(id)
	if checkErr != nil {
		return utils.NotFoundResponse(c, "User not found", nil)
	}

	detailsMap, mapErr := convertToMap(existingUser.Details)
	if mapErr != nil {
		logger.LogError(mapErr, "Failed to convert details to map", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	if (detailsMap["SchoolUUID"] != "N/A" && detailsMap["SchoolUUID"] != "") && forceDelete != "true" {
		return utils.BadRequestResponse(c, "Warning: This driver is still associated with a school, continue?", nil)
	}

	err := handler.userService.DeleteDriver(id, username)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete driver", nil)
	}

	return utils.SuccessResponse(c, "Driver deleted successfully", nil)
}

func (handler *userHandler) GetSpecSuperAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := handler.userService.GetSpecSuperAdmin(id)
	if err != nil {
		logger.LogError(err, "Failed to fetch specific super admin", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "User fetched successfully", user)
}

func (handler *userHandler) GetSpecSchoolAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := handler.userService.GetSpecSchoolAdmin(id)
	if err != nil {
		logger.LogError(err, "Failed to fetch specific school admin", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "User fetched successfully", user)
}

// func GetAllPermittedDriver(c *fiber.Ctx) error {
// 	role := c.Locals("role_code").(string)

//     var users []models.UserResponse
//     var err error

//     if role == "SA" {
//         users, err = services.GetAllDriverFromAllSchools()
//     } else if role == "AS" {
// 		schoolObjID, parseErr := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
// 		if parseErr != nil {
// 			logger.LogError(parseErr, "Failed to convert school id", map[string]interface{}{
// 				"school_id": c.Locals("schoolId"),
// 			})
// 			return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 		}

//         users, err = services.GetAllDriverForPermittedSchool(schoolObjID)
//     }

//     if err != nil {
// 		logger.LogError(err, "Failed to fetch all drivers", nil)
//         return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
//     }

//     return c.Status(fiber.StatusOK).JSON(users)
// }

// func GetSpecPermittedDriver(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	role := c.Locals("role_code").(string)

// 	var user models.UserResponse
// 	var err error

// 	if role == "SA" {
// 		user, err = services.GetSpecDriverFromAllSchools(id)
// 	} else if role == "AS" {
// 		schoolObjID, parseErr := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
// 		if parseErr != nil {
// 			logger.LogError(parseErr, "Failed to convert school id", map[string]interface{}{
// 				"school_id": c.Locals("schoolId"),
// 			})
// 			return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 		}

// 		user, err = services.GetSpecDriverForPermittedSchool(id, schoolObjID)
// 	}

// 	if err != nil {
// 		logger.LogError(err, "Failed to fetch specific driver", nil)
// 		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 	}

// 	return c.Status(fiber.StatusOK).JSON(user)
// }

func (handler *userHandler) AddUser(c *fiber.Ctx) error {
	username := c.Locals("user_name").(string)

	userReqDTO := new(dto.UserRequestsDTO)
	if err := c.BodyParser(userReqDTO); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, userReqDTO); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}

	if err := validateUserRoleDetails(c, userReqDTO, *handler); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}

	if _, err := handler.userService.AddUser(*userReqDTO, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to create user", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "User created successfully", nil)
}

func (handler *userHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Locals("user_name").(string)

	userReqDTO := new(dto.UserRequestsDTO)
	if err := c.BodyParser(userReqDTO); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	existingUser, err := handler.userService.GetSpecUserWithDetails(id)
	if err != nil {
		logger.LogError(err, "Failed to fetch user", nil)
		return utils.NotFoundResponse(c, "User not found", nil)
	}

	userReqDTO.Password = existingUser.Password

	if err := utils.ValidateStruct(c, userReqDTO); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}

	if err := validateUserRoleDetails(c, userReqDTO, *handler); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}

	detailsMap, err := convertToMap(existingUser.Details)
	if err != nil {
		logger.LogError(err, "Failed to convert details to map", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	if err := handler.userService.UpdateUser(id, *userReqDTO, username, detailsMap, nil); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to update user", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "User updated successfully", nil)
}

func (handler *userHandler) DeleteSuperAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Locals("user_name").(string)

	if err := handler.userService.DeleteSuperAdmin(id, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to delete super admin", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Super admin deleted successfully", nil)
}

func (handler *userHandler) DeleteSchoolAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Locals("user_name").(string)

	forceDelete := c.Query("force_delete")

	existingUser, checkErr := handler.userService.GetSpecUserWithDetails(id)
	if checkErr != nil {
		return utils.NotFoundResponse(c, "User not found", nil)
	}

	detailsMap, mapErr := convertToMap(existingUser.Details)
	if mapErr != nil {
		logger.LogError(mapErr, "Failed to convert details to map", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	if (detailsMap["SchoolUUID"] != "N/A" && detailsMap["SchoolUUID"] != "") && forceDelete != "true" {
		return utils.BadRequestResponse(c, "Warning: This school admin is still associated with a school, continue?", nil)
	}

	err := handler.userService.DeleteSchoolAdmin(id, username)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete school admin", nil)
	}

	return utils.SuccessResponse(c, "School admin deleted successfully", nil)
}

// func validateCommonUserFields(c *fiber.Ctx, user *dto.UserRequestsDTO, handler userHandler) error {
// 	if !regexp.MustCompile(`^[a-zA-Z0-9_-]{3,}$`).MatchString(user.Username) {
// 		return fmt.Errorf("username must be at least 3 characters and contain only alphanumeric characters, hyphens, and underscores")
// 	}

// 	if len(user.Password) < 8 {
// 		return fmt.Errorf("password must be at least 8 characters")
// 	}

// 	if _, err := mail.ParseAddress(user.Email); err != nil {
// 		return fmt.Errorf("invalid email address format")
// 	}

// 	if !regexp.MustCompile(`^\+?[0-9]{12,15}$`).MatchString(user.Phone) {
// 		return fmt.Errorf("invalid phone number format")
// 	}

// 	validRoles := map[dto.Role]bool{
// 		dto.SuperAdmin: true,
// 		dto.SchoolAdmin: true,
// 		dto.Parent: true,
// 		dto.Driver: true,
// 	}
// 	if !validRoles[user.Role] {
// 		return fmt.Errorf("invalid role")
// 	}

// 	if user.Details != nil {
// 		if err := validateUserRoleDetails(c, user, handler); err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	return nil
// }

func validateUserRoleDetails(c *fiber.Ctx, user *dto.UserRequestsDTO, handler userHandler) error {
	switch user.Role {
	case dto.SuperAdmin:
		user.RoleCode = "SA"

	case dto.SchoolAdmin:
		details, err := parseDetails[dto.SchoolAdminDetailsRequestsDTO](user.Details)
		if err != nil {
			logger.LogError(err, "Invalid details format for SchoolAdmin", map[string]interface{}{
				"details": user.Details,
			})
			return errors.New("invalid details format for SchoolAdmin", 400)
		}

		if details.SchoolUUID == "" {
			return errors.New("school is required for SchoolAdmin", 400)
		}

		_, errSchool := handler.schoolService.GetSpecSchool(details.SchoolUUID)
		if errSchool != nil {
			return errors.New("school is not found", 404)
		}

		user.RoleCode = "AS"
		user.Details = details

	case dto.Parent:
		if user.Details == nil {
			return errors.New("parent details are required", 400)
		}
		user.RoleCode = "P"

	case dto.Driver:
		details, err := parseDetails[dto.DriverDetailsRequestsDTO](user.Details)
		if err != nil {
			logger.LogError(err, "Invalid details format for Driver", map[string]interface{}{
				"details": user.Details,
			})
			return errors.New("invalid details format for Driver", 400)
		}

		// if details.VehicleID != 0 {
		//     _, err := GetSpecVehicle(details.VehicleID)
		//     if err != nil {
		//         return errors.New("vehicle not found")
		//     }
		// }

		// if details.SchoolID != 0 {
		//     _, err := GetSpecSchool(details.SchoolID)
		//     if err != nil {
		//         return errors.New("school not found")
		//     }
		// }

		user.RoleCode = "D"
		user.Details = details

	default:
		return errors.New("invalid role specified", 400)
	}
	return nil
}

func parseDetails[T any](details interface{}) (T, error) {
	var parsedDetails T

	detailsMap, ok := details.(map[string]interface{})
	if !ok {
		return parsedDetails, fmt.Errorf("invalid details format: expected map[string]interface{}, got %T", details)
	}

	detailsBytes, err := json.Marshal(detailsMap)
	if err != nil {
		return parsedDetails, fmt.Errorf("failed to marshal details to JSON: %w", err)
	}

	err = json.Unmarshal(detailsBytes, &parsedDetails)
	if err != nil {
		return parsedDetails, fmt.Errorf("failed to unmarshal details to struct: %w", err)
	}

	return parsedDetails, nil
}

func isValidSortFieldForUsers(field string) bool {
	allowedFields := map[string]bool{
		"user_id":         true,
		"user_username":   true,
		"user_first_name": true,
		"user_last_name":  true,
		"created_at":      true,
	}
	return allowedFields[field]
}

func safeString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

func convertToMap(data interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
