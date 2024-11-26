package controllers

import (
	"log"
	"shuttle/models"
	"shuttle/services"
	"shuttle/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(c *fiber.Ctx) error {
	loginReq := new(models.LoginRequest)
	if err := c.BodyParser(loginReq); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	user, err := services.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid email or password", nil)
	}

	Fullname := user.FirstName + " " + user.LastName

	// Access token (short expiration)
	accessToken, err := utils.GenerateToken(user.ID.Hex(), Fullname, string(user.Role), user.RoleCode)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate access token", nil)
	}

	// Refresh token (long expiration)
	refreshToken, err := utils.GenerateRefreshToken(user.ID.Hex(), Fullname, string(user.Role), user.RoleCode)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate refresh token", nil)
	}

	// Save refresh token in the database
	err = utils.SaveRefreshToken(user.ID.Hex(), refreshToken)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to save refresh token", nil)
	}

	responseData := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return utils.SuccessResponse(c, "User logged in successfully", responseData)
}

func Logout(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token missing", nil)
	}

	UserID, err := utils.GetUserIDFromToken(token)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	// Delete WebSocket connection if exists
	conn, exists := utils.GetConnection(UserID)
	if exists {
		conn.Close()
		utils.RemoveConnection(UserID)
		log.Println("User disconnected from WebSocket:", UserID)
	}

	err = services.DeleteRefreshTokenOnLogout(UserID)
	if err != nil {
		log.Println("Error deleting refresh token:", err)
	}

	utils.InvalidateToken(token)

	ObjectID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error updating user status", nil)
	}

	err = services.UpdateUserStatus(ObjectID, "offline", time.Now())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error updating user status", nil)
	}

	return utils.SuccessResponse(c, "User logged out successfully", map[string]interface{}{})
}

func GetMyProfile(c *fiber.Ctx) error {
	token := string(c.Request().Header.Peek("Authorization"))
	UserID, err := utils.GetUserIDFromToken(token)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}
	RoleCode, err := utils.GetRoleCodeFromToken(token)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	user, err := services.GetMyProfile(UserID, RoleCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// Reissue a new access token
func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Get("Authorization")
	if refreshToken == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Refresh token missing", nil)
	}

	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid refresh token", nil)
	}

	userID := claims["userId"].(string)

	storedRefreshToken, err := services.GetStoredRefreshToken(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Refresh token not found or invalid", nil)
	}

	if storedRefreshToken != refreshToken {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Refresh token mismatch", nil)
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expirationTime) {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Refresh token expired", nil)
	}

	userId := claims["userId"].(string)
	name := claims["name"].(string)
	role := claims["role"].(string)
	role_code := claims["role_code"].(string)

	// Generate new access token
	accessToken, err := utils.GenerateToken(userId, name, role, role_code)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to generate access token", nil)
	}

	return utils.SuccessResponse(c, "Access token refreshed", map[string]interface{}{
		"access_token": accessToken,
	})
}