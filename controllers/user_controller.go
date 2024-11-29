package controllers

import (
	"log"
	"shuttle/errors"
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetAllSuperAdmin(c *fiber.Ctx) error {
	users, err := services.GetAllSuperAdmin()
	if err != nil {
		log.Print(err)
		return utils.InternalServerErrorResponse(c, "Internal server error", nil)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func GetSpecSuperAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := services.GetSpecSuperAdmin(id)
	if err != nil {
		log.Print(err)
		return utils.InternalServerErrorResponse(c, "Internal server error", nil)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func GetAllSchoolAdmin(c *fiber.Ctx) error {
	users, err := services.GetAllSchoolAdmin()
	if err != nil {
		log.Print(err)
		return utils.InternalServerErrorResponse(c, "Internal server error", nil)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func GetSpecSchoolAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := services.GetSpecSchoolAdmin(id)
	if err != nil {
		log.Print(err)
		return utils.InternalServerErrorResponse(c, "Internal server error", nil)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func GetAllPermittedDriver(c *fiber.Ctx) error {
    userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return utils.UnauthorizedResponse(c, "User ID is missing or invalid", nil)
	}

    role, ok := c.Locals("role_code").(string)
	if !ok || role == "" {
		return utils.UnauthorizedResponse(c, "Role code is missing or invalid", nil)
	}

    var users []models.UserResponse
    var err error

    if role == "SA" {
        users, err = services.GetAllDriverFromAllSchools()
    } else if role == "AS" {
        schoolID, permErr := services.CheckPermittedSchoolAccess(userID)
        if permErr != nil {
            return utils.ForbiddenResponse(c, "You don't have permission to access this resource" + permErr.Error(), nil)
        }
		
        users, err = services.GetAllDriverForSchool(schoolID)
    } else {
        return utils.ForbiddenResponse(c, "You don't have permission to access this resource", nil)
    }

    if err != nil {
        return utils.InternalServerErrorResponse(c, "Failed to fetch drivers: " + err.Error(), nil)
    }

    return c.Status(fiber.StatusOK).JSON(users)
}

func GetSpecPermittedDriver(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return utils.UnauthorizedResponse(c, "User ID is missing or invalid", nil)
	}

	role, ok := c.Locals("role_code").(string)
	if !ok || role == "" {
		return utils.UnauthorizedResponse(c, "Role code is missing or invalid", nil)
	}

	id := c.Params("id")
	var user models.UserResponse
	var err error

	if role == "SA" {
		user, err = services.GetSpecDriverFromAllSchools(id)
	} else if role == "AS" {
		schoolID, permErr := services.CheckPermittedSchoolAccess(userID)
		if permErr != nil {
			return utils.ForbiddenResponse(c, "You don't have permission to access this resource" + permErr.Error(), nil)
		}
		
		user, err = services.GetSpecDriverForSchool(id, schoolID)
	} else {
		return utils.ForbiddenResponse(c, "You don't have permission to access this resource", nil)
	}

	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch driver: " + err.Error(), nil)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func AddUser(c *fiber.Ctx) error {
	username, ok := c.Locals("username").(string)
	if !ok || username == "" {
		return utils.UnauthorizedResponse(c, "Username is missing or invalid", nil)
	}

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return utils.BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
		}
	}

	if _, err := services.AddUser(*user, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		return utils.InternalServerErrorResponse(c, "Failed to create user: "+err.Error(), nil)
	}

	return utils.SuccessResponse(c, "User created successfully", nil)
}

func AddSchoolDriver(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return utils.UnauthorizedResponse(c, "User ID is missing or invalid", nil)
	}

	username, ok := c.Locals("username").(string)
	if !ok || username == "" {
		return utils.UnauthorizedResponse(c, "Username is missing or invalid", nil)
	}

	SchoolID, err := services.CheckPermittedSchoolAccess(userID)
	if err != nil {
		return utils.UnauthorizedResponse(c, "You don't have permission to create driver for this school", nil)
	}
	
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return utils.BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
		}
	}

	driverDetails := map[string]interface{}{
		"license_number": user.Details.(map[string]interface{})["license_number"],
		"school_id": SchoolID.Hex(),
		"vehicle_id": user.Details.(map[string]interface{})["vehicle_id"],
	}

	user.Role = models.Driver
	user.Details = driverDetails

	if _, err := services.AddUser(*user, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		return utils.InternalServerErrorResponse(c, "Failed to create driver: "+err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Driver created successfully", nil)
}

func UpdateUser(c *fiber.Ctx) error {
    username, ok := c.Locals("username").(string)
    if !ok || username == "" {
        return utils.UnauthorizedResponse(c, "Username is missing or invalid", nil)
    }

    id := c.Params("id")

    existingUser, err := services.GetSpecUser(id)
    if err != nil {
        return utils.InternalServerErrorResponse(c, "Failed to fetch user: "+err.Error(), nil)
    }

    user, err := parseFormData(c, &existingUser)
    if err != nil {
        return utils.BadRequestResponse(c, err.Error(), nil)
    }

    roleStr := c.FormValue("role", string(existingUser.Role))
    role, err := models.ParseRole(roleStr)
    if err != nil {
        return utils.BadRequestResponse(c, "Invalid role", nil)
    }
    user.Role = role

    validate := validator.New()
    if err := validate.Struct(user); err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            return utils.BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
        }
    }

    user.Picture, err = utils.HandleAssetsOnUpdate(c, existingUser.Picture)
    if err != nil {
        return utils.InternalServerErrorResponse(c, err.Error(), nil)
    }

    if err := services.UpdateUser(id, *user, username, nil); err != nil {
        return utils.InternalServerErrorResponse(c, "Failed to update user: "+err.Error(), nil)
    }

    return utils.SuccessResponse(c, "User updated successfully", nil)
}

func UpdateSchoolDriver(c *fiber.Ctx) error {
    username, ok := c.Locals("username").(string)
	if !ok || username == "" {
		return utils.UnauthorizedResponse(c, "Username is missing or invalid", nil)
	}

    id := c.Params("id")
    existingDriver, err := services.GetSpecUser(id)
    if err != nil {
        return utils.NotFoundResponse(c, "Driver not found", nil)
    }

    UserID, ok := c.Locals("userId").(string)
	if !ok || UserID == "" {
		return utils.UnauthorizedResponse(c, "User ID is missing or invalid", nil)
	}

    SchoolID, err := services.CheckPermittedSchoolAccess(UserID)
    if err != nil {
        return utils.UnauthorizedResponse(c, "You don't have permission to update driver for this school", nil)
    }

    user := new(models.User)
    if err := c.BodyParser(user); err != nil {
        return utils.BadRequestResponse(c, "Invalid request data", nil)
    }

	user, err = parseFormData(c, &existingDriver)
    if err != nil {
        return utils.BadRequestResponse(c, err.Error(), nil)
    }

    licenseNumber := c.FormValue("license_number", existingDriver.Details.(map[string]interface{})["license_number"].(string))
    vehicleID := c.FormValue("vehicle_id", "")
    if vehicleID != "" {
        _, err := services.GetSpecVehicle(vehicleID)
        if err != nil {
            return utils.BadRequestResponse(c, "Vehicle not found or invalid vehicle_id", nil)
        }
    }

    user.Role = models.Driver
    driverDetails := map[string]interface{}{
        "license_number": licenseNumber,
        "school_id":      SchoolID.Hex(),
        "vehicle_id":     vehicleID,
    }
    user.Details = driverDetails

    validate := validator.New()
    if err := validate.Struct(user); err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            return utils.BadRequestResponse(c, err.Field()+" is "+err.Tag(), nil)
        }
    }

    user.Picture, err = utils.HandleAssetsOnUpdate(c, existingDriver.Picture)
    if err != nil {
        return utils.InternalServerErrorResponse(c, err.Error(), nil)
    }

    if err := services.UpdateUser(id, *user, username, nil); err != nil {
        return utils.InternalServerErrorResponse(c, "Failed to update driver: " + err.Error(), nil)
    }

    return utils.SuccessResponse(c, "Driver updated successfully", nil)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := services.DeleteUser(id); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete user: " + err.Error(), nil)
	}

	return utils.SuccessResponse(c, "User deleted successfully", nil)
}

func DeleteSchoolDriver(c *fiber.Ctx) error {
	UserID, ok := c.Locals("userId").(string)
	if !ok || UserID == "" {
		return utils.UnauthorizedResponse(c, "Invalid Token", nil)
	}

	SchoolID, err := services.CheckPermittedSchoolAccess(UserID)
	if err != nil {
		return utils.UnauthorizedResponse(c, "You don't have permission to access this resource", nil)
	}

	id := c.Params("id")
	if err := services.DeleteSchoolDriver(id, SchoolID); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to delete driver: " + err.Error(), nil)
	}

	return utils.SuccessResponse(c, "Driver deleted successfully", nil)
}

func parseFormData(c *fiber.Ctx, existingUser *models.User) (*models.User, error) {
    user := new(models.User)
    if err := c.BodyParser(user); err != nil {
        return nil, errors.New("Invalid request data", 400)
    }

    user.FirstName = c.FormValue("first_name", existingUser.FirstName)
    user.LastName = c.FormValue("last_name", existingUser.LastName)
	user.Email = c.FormValue("email", existingUser.Email)
    user.Password = c.FormValue("password", existingUser.Password)
	
    genderStr := c.FormValue("gender", string(existingUser.Gender))
    gender, err := models.ParseGender(genderStr)
    if err != nil {
        return nil, errors.New("Invalid gender", 400)
    }
    user.Gender = gender

    user.Phone = c.FormValue("phone", existingUser.Phone)
    user.Address = c.FormValue("address", existingUser.Address)

    return user, nil
}