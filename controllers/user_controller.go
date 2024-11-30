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

func GetAllSuperAdmin(c *fiber.Ctx) error {
	users, err := services.GetAllSuperAdmin()
	if err != nil {
		logger.LogError(err, "Failed to fetch all super admins", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func GetSpecSuperAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := services.GetSpecSuperAdmin(id)
	if err != nil {
		logger.LogError(err, "Failed to fetch specific super admin", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func GetAllSchoolAdmin(c *fiber.Ctx) error {
	users, err := services.GetAllSchoolAdmin()
	if err != nil {
		logger.LogError(err, "Failed to fetch all school admins", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func GetSpecSchoolAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := services.GetSpecSchoolAdmin(id)
	if err != nil {
		logger.LogError(err, "Failed to fetch specific school admin", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func GetAllPermittedDriver(c *fiber.Ctx) error {
	role := c.Locals("role_code").(string)

    var users []models.UserResponse
    var err error

    if role == "SA" {
        users, err = services.GetAllDriverFromAllSchools()
    } else if role == "AS" {
		schoolObjID, parseErr := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
		if parseErr != nil {
			logger.LogError(parseErr, "Failed to convert school id", map[string]interface{}{
				"school_id": c.Locals("schoolId"),
			})
			return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
		}

        users, err = services.GetAllDriverForPermittedSchool(schoolObjID)
    }

    if err != nil {
		logger.LogError(err, "Failed to fetch all drivers", nil)
        return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
    }

    return c.Status(fiber.StatusOK).JSON(users)
}

func GetSpecPermittedDriver(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("role_code").(string)

	var user models.UserResponse
	var err error

	if role == "SA" {
		user, err = services.GetSpecDriverFromAllSchools(id)
	} else if role == "AS" {
		schoolObjID, parseErr := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
		if parseErr != nil {
			logger.LogError(parseErr, "Failed to convert school id", map[string]interface{}{
				"school_id": c.Locals("schoolId"),
			})
			return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
		}
		
		user, err = services.GetSpecDriverForPermittedSchool(id, schoolObjID)
	}

	if err != nil {
		logger.LogError(err, "Failed to fetch specific driver", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func AddUser(c *fiber.Ctx) error {
	username := c.Locals("username").(string)

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, user); err != nil {
		return err
	}

	if _, err := services.AddUser(*user, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to create user", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "User created successfully", nil)
}

func AddSchoolDriver(c *fiber.Ctx) error {
	username := c.Locals("username").(string)

	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
	if err != nil {
		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
			"school_id": c.Locals("schoolId"),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, user); err != nil {
		return err
	}

	driverDetails := map[string]interface{}{
		"license_number": user.Details.(map[string]interface{})["license_number"],
		"school_id": SchoolObjID.Hex(),
		"vehicle_id": user.Details.(map[string]interface{})["vehicle_id"],
	}

	user.Role = models.Driver
	user.Details = driverDetails

	if _, err := services.AddUser(*user, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to create driver", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Driver created successfully", nil)
}

func UpdateUser(c *fiber.Ctx) error {
    id := c.Params("id")
    username := c.Locals("username").(string)

    existingUser, err := services.GetSpecUser(id)
    if err != nil {
        return utils.NotFoundResponse(c, "User not found", nil)
    }

    user, err := parseFormData(c, &existingUser)
    if err != nil {
        return utils.BadRequestResponse(c, "Invalid request type", nil)
    }

    roleStr := c.FormValue("role", string(existingUser.Role))
    role, err := models.ParseRole(roleStr)
    if err != nil {
        return utils.BadRequestResponse(c, "Invalid role", nil)
    }
    user.Role = role

    if err := utils.ValidateStruct(c, user); err != nil {
		return err
	}

    user.Picture, err = utils.HandleAssetsOnUpdate(c, existingUser.Picture)
    if err != nil {
		logger.LogError(err, "Failed to handle assets", nil)
        return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
    }

    if err := services.UpdateUser(id, *user, username, nil); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to update user", nil)
        return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
    }

    return utils.SuccessResponse(c, "User updated successfully", nil)
}

func UpdateSchoolDriver(c *fiber.Ctx) error {
	id := c.Params("id")
    username := c.Locals("username").(string)

    existingDriver, err := services.GetSpecUser(id)
    if err != nil {
        return utils.NotFoundResponse(c, "Driver not found", nil)
    }

    SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
	if err != nil {
		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
			"school_id": c.Locals("schoolId"),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

    user := new(models.User)
    if err := c.BodyParser(user); err != nil {
        return utils.BadRequestResponse(c, "Invalid request data", nil)
    }

	user, err = parseFormData(c, &existingDriver)
    if err != nil {
        return utils.BadRequestResponse(c, "Invalid request type", nil)
    }

    licenseNumber := c.FormValue("license_number", existingDriver.Details.(map[string]interface{})["license_number"].(string))
    vehicleID := c.FormValue("vehicle_id", "")
    if vehicleID != "" {
        _, err := services.GetSpecVehicle(vehicleID)
        if err != nil {
            return utils.BadRequestResponse(c, "Vehicle not found or invalid", nil)
        }
    }

    user.Role = models.Driver
    driverDetails := map[string]interface{}{
        "license_number": licenseNumber,
        "school_id":      SchoolObjID.Hex(),
        "vehicle_id":     vehicleID,
    }
    user.Details = driverDetails

    if err := utils.ValidateStruct(c, user); err != nil {
		return err
	}

    user.Picture, err = utils.HandleAssetsOnUpdate(c, existingDriver.Picture)
    if err != nil {
		logger.LogError(err, "Failed to handle assets", nil)
        return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
    }

    if err := services.UpdateUser(id, *user, username, nil); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to update driver", nil)
        return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
    }

    return utils.SuccessResponse(c, "Driver updated successfully", nil)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := services.DeleteUser(id); err != nil {
		logger.LogError(err, "Failed to delete user", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "User deleted successfully", nil)
}

func DeleteSchoolDriver(c *fiber.Ctx) error {
	id := c.Params("id")

	SchoolObjID, err := primitive.ObjectIDFromHex(c.Locals("schoolId").(string))
	if err != nil {
		logger.LogError(err, "Failed to convert school id", map[string]interface{}{
			"school_id": c.Locals("schoolId"),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}
	
	if err := services.DeleteSchoolDriver(id, SchoolObjID); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to delete driver", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Driver deleted successfully", nil)
}

func parseFormData(c *fiber.Ctx, existingUser *models.User) (*models.User, error) {
    user := new(models.User)
    if err := c.BodyParser(user); err != nil {
        return nil, nil
    }

    user.FirstName = c.FormValue("first_name", existingUser.FirstName)
    user.LastName = c.FormValue("last_name", existingUser.LastName)
	user.Email = c.FormValue("email", existingUser.Email)
    user.Password = c.FormValue("password", existingUser.Password)
	
    genderStr := c.FormValue("gender", string(existingUser.Gender))
    gender, err := models.ParseGender(genderStr)
    if err != nil {
        return nil, nil
    }
    user.Gender = gender

    user.Phone = c.FormValue("phone", existingUser.Phone)
    user.Address = c.FormValue("address", existingUser.Address)

    return user, nil
}