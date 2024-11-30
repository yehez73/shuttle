package controllers

import (
	"shuttle/logger"
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
		return utils.BadRequestResponse(c, "Invalid request data", 400)
	}

	user, err := services.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		logger.LogWarn("Login attempt failed for user", map[string]interface{}{
			"email": loginReq.Email,
			"error": err.Error(),
		})
		return utils.UnauthorizedResponse(c, "Invalid email or password", nil)
	}

	logger.LogInfo("User logged in", map[string]interface{}{
		"user_id": user.ID.Hex(),
		"email":   loginReq.Email,
	})

	Fullname := user.FirstName + " " + user.LastName

	// Access token (short expiration)
	accessToken, err := utils.GenerateToken(user.ID.Hex(), Fullname, string(user.Role), user.RoleCode)
	if err != nil {
		logger.LogError(err, "Failed to generate access token", map[string]interface{}{
			"user_id": user.ID.Hex(),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	// Refresh token (long expiration)
	refreshToken, err := utils.GenerateRefreshToken(user.ID.Hex(), Fullname, string(user.Role), user.RoleCode)
	if err != nil {
		logger.LogError(err, "Failed to generate refresh token", map[string]interface{}{
			"user_id": user.ID.Hex(),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	// Save refresh token in the database
	err = utils.SaveRefreshToken(user.ID.Hex(), refreshToken)
	if err != nil {
		logger.LogError(err, "Failed to save refresh token", map[string]interface{}{
			"user_id": user.ID.Hex(),
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	responseData := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return utils.SuccessResponse(c, "User logged in successfully", responseData)
}

func Logout(c *fiber.Ctx) error {
	UserID := c.Locals("userId").(string)

	// Delete WebSocket connection if exists
	conn, exists := utils.GetConnection(UserID)
	if exists {
		conn.Close()
		utils.RemoveConnection(UserID)
		logger.LogInfo("WebSocket connection closed", map[string]interface{}{
			"user_id": UserID,
		})
	}

	err := services.DeleteRefreshTokenOnLogout(UserID)
	if err != nil {
		logger.LogError(err, "Failed to delete refresh token", map[string]interface{}{
			"user_id": UserID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	utils.InvalidateToken(c.Get("Authorization"))

	ObjectID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		logger.LogError(err, "Invalid user ID", map[string]interface{}{
			"user_id": UserID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	err = services.UpdateUserStatus(ObjectID, "offline", time.Now())
	if err != nil {
		logger.LogError(err, "Failed to update user status", map[string]interface{}{
			"user_id": UserID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "User logged out successfully", nil)
}

func GetMyProfile(c *fiber.Ctx) error {
	UserID := c.Locals("userId").(string)
	RoleCode := c.Locals("roleCode").(string)
	
	user, err := services.GetMyProfile(UserID, RoleCode)
	if err != nil {
		logger.LogError(err, "Failed to get user profile", map[string]interface{}{
			"user_id": UserID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// Reissue a new access token
func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Get("Authorization")
	if refreshToken == "" {
		return utils.UnauthorizedResponse(c, "Missing refresh token", nil)
	}

	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		logger.LogWarn("Invalid refresh token", map[string]interface{}{})
		return utils.UnauthorizedResponse(c, "Invalid refresh token", nil)
	}

	userID := claims["userId"].(string)

	storedRefreshToken, err := services.GetStoredRefreshToken(userID)
	if err != nil {
		logger.LogError(err, "Failed to get stored refresh token", map[string]interface{}{
			"user_id": userID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	if storedRefreshToken != refreshToken {
		logger.LogWarn("Invalid refresh token", map[string]interface{}{})
		return utils.UnauthorizedResponse(c, "Invalid refresh token", nil)
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expirationTime) {
		return utils.UnauthorizedResponse(c, "Refresh token has expired", nil)
	}

	userId := claims["userId"].(string)
	name := claims["name"].(string)
	role := claims["role"].(string)
	role_code := claims["role_code"].(string)

	// Generate new access token
	accessToken, err := utils.GenerateToken(userId, name, role, role_code)
	if err != nil {
		logger.LogError(err, "Failed to generate access token", map[string]interface{}{
			"user_id": userId,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Access token refreshed", map[string]interface{}{
		"access_token": accessToken,
	})
}