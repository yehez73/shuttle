package controllers

import (
	"io"
	"log"
	"net/http"
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetAllUser(c *fiber.Ctx) error {
	users, err := services.GetAllUser()
	if err != nil {
		log.Print(err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func GetSpecUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := services.GetSpecUser(id)
	if err != nil {
		log.Print(err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Internal server error", nil)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func AddUser(c *fiber.Ctx) error {
	token := string(c.Request().Header.Peek("Authorization"))
	username, err := utils.GetUsernameFromToken(token)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
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
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user: " + err.Error(), nil)
	}

	return utils.SuccessResponse(c, "User created successfully", nil)
}

func UpdateUser(c *fiber.Ctx) error {
	token := string(c.Request().Header.Peek("Authorization"))
	username, err := utils.GetUsernameFromToken(token)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	id := c.Params("id")
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

	file, err := c.FormFile("picture")
	if err != nil {
		return utils.BadRequestResponse(c, "Picture is required", nil)
	}

	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error", nil)
	}

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process image", nil)
	}

	if err := services.UpdateUser(id, *user, username, fileBytes); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user: " + err.Error(), nil)
	}

	return utils.SuccessResponse(c, "User updated successfully", nil)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := services.DeleteUser(id); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user: " + err.Error(), nil)
	}

	return utils.SuccessResponse(c, "User deleted successfully", nil)
}