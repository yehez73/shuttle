package handler

import (
	"encoding/json"
	"fmt"
	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandlerInterface interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	GetMyProfile(c *fiber.Ctx) error
	UpdateMyProfile(c *fiber.Ctx) error
	IssueNewAccessToken(c *fiber.Ctx) error
	AddDeviceToken(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
}

type authHandler struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthHttpHandler(authService services.AuthService, userService services.UserService) AuthHandlerInterface {
	return &authHandler{
		authService: authService,
		userService: userService,
	}
}

func (handler *authHandler) Login(c *fiber.Ctx) error {
	loginRequest := new(dto.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	userDataOnLogin, err := handler.authService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		logger.LogError(err, "Failed to login", map[string]interface{}{
			"email": loginRequest.Email,
		})
		return utils.UnauthorizedResponse(c, "Invalid email or password", nil)
	}

	logger.LogInfo("User logged in", map[string]interface{}{
		"id":    userDataOnLogin.UserID,
		"email": loginRequest.Email,
	})

	// Access token (short expiration)
	accessToken, err := utils.GenerateToken(fmt.Sprintf("%d", userDataOnLogin.UserID), userDataOnLogin.UserUUID, userDataOnLogin.Username, userDataOnLogin.RoleCode)
	if err != nil {
		logger.LogError(err, "Failed to generate access token", map[string]interface{}{
			"user_id": userDataOnLogin.UserID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	// Refresh token (long expiration)
	refreshToken, err := utils.GenerateRefreshToken(fmt.Sprintf("%d", userDataOnLogin.UserID), userDataOnLogin.UserUUID, userDataOnLogin.Username, userDataOnLogin.RoleCode)
	if err != nil {
		logger.LogError(err, "Failed to generate refresh token", map[string]interface{}{
			"user_id": userDataOnLogin.UserID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	// Save refresh token in the database
	err = utils.SaveRefreshToken(userDataOnLogin.UserUUID, refreshToken)
	if err != nil {
		logger.LogError(err, "Failed to save refresh token", map[string]interface{}{
			"user_id": userDataOnLogin.UserID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	responseData := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return utils.SuccessResponse(c, "User logged in successfully", responseData)
}

func (handler *authHandler) Logout(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok {
		return utils.UnauthorizedResponse(c, "Token is invalid", nil)
	}

	// Delete WebSocket connection if exists
	// conn, exists := utils.GetConnection(userUUID)
	// if exists {
	// 	conn.Close()
	// 	utils.RemoveConnection(userUUID)
	// 	logger.LogInfo("WebSocket connection closed", map[string]interface{}{
	// 		"user_uuid": userUUID,
	// 	})
	// }

	err := handler.authService.DeleteRefreshTokenOnLogout(c.Context(), userUUID)
	if err != nil {
		logger.LogError(err, "Failed to delete refresh token", map[string]interface{}{
			"user_uuid": userUUID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	utils.InvalidateToken(c.Get("Authorization"))

	err = handler.authService.UpdateUserStatus(userUUID, "offline", time.Now())
	if err != nil {
		logger.LogError(err, "Failed to update user status", map[string]interface{}{
			"user_uuid": userUUID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "User logged out successfully", nil)
}

func (handler *authHandler) GetMyProfile(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok {
		return utils.UnauthorizedResponse(c, "Token is invalid", nil)
	}

	roleCode, ok := c.Locals("role_code").(string)
	if !ok {
		return utils.UnauthorizedResponse(c, "Token is invalid", nil)
	}

	user, err := handler.authService.GetMyProfile(userUUID, roleCode)
	if err != nil {
		logger.LogError(err, "Failed to get user profile", map[string]interface{}{
			"user_uuid": userUUID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "User profile retrieved", user)
}

func (handler *authHandler) UpdateMyProfile(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok {
		return utils.UnauthorizedResponse(c, "Token is invalid", nil)
	}

	username, ok := c.Locals("user_name").(string)
	if !ok {
		return utils.UnauthorizedResponse(c, "Token is invalid", nil)
	}

	updateRequest := new(dto.UserRequestsDTO)
	if err := c.BodyParser(updateRequest); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	existingUser, err := handler.userService.GetSpecUserWithDetails(userUUID)
	if err != nil {
		logger.LogError(err, "Failed to fetch user", nil)
		return utils.NotFoundResponse(c, "User not found", nil)
	}

	updateRequest.Password = existingUser.User.Password
	updateRequest.Role = dto.Role(existingUser.User.Role)
	updateRequest.RoleCode = existingUser.User.RoleCode
	// Marshall existing user details to JSON.RawMessage whether its a super admin, school admin, or driver
	mergedDetails, err := mergeDetails(updateRequest.Details, existingUser.Details)
	if err != nil {
		logger.LogError(err, "Failed to merge user details", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}
	updateRequest.Details = mergedDetails

	err = handler.userService.UpdateUser(userUUID, *updateRequest, username, nil)
	if err != nil {
		logger.LogError(err, "Failed to update user profile", map[string]interface{}{
			"user_uuid": userUUID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "User profile updated", nil)
}

func (handler *authHandler) ChangePassword(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok {
		return utils.UnauthorizedResponse(c, "Token is invalid", nil)
	}

	changePasswordRequest := new(dto.ChangePasswordRequest)
	if err := c.BodyParser(changePasswordRequest); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, changePasswordRequest); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}

	existingUser, err := handler.userService.GetSpecUserWithDetails(userUUID)
	if err != nil {
		logger.LogError(err, "Failed to fetch user", nil)
		return utils.NotFoundResponse(c, "User not found", nil)
	}

	if !utils.ValidatePassword(changePasswordRequest.OldPassword, existingUser.User.Password) {
		return utils.BadRequestResponse(c, "Invalid old password", nil)
	}

	err = handler.authService.ChangePassword(userUUID, changePasswordRequest.NewPassword)
	if err != nil {
		logger.LogError(err, "Failed to change password", map[string]interface{}{
			"user_uuid": userUUID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Password changed successfully", nil)
}

// Reissue a new access token
func (handler *authHandler) IssueNewAccessToken(c *fiber.Ctx) error {
	refreshToken := c.Get("Authorization")
	if refreshToken == "" {
		return utils.UnauthorizedResponse(c, "Missing refresh token", nil)
	}

	// Remove "Bearer " prefix
	const bearerPrefix = "Bearer "
	if len(refreshToken) > len(bearerPrefix) && refreshToken[:len(bearerPrefix)] == bearerPrefix {
		refreshToken = refreshToken[len(bearerPrefix):]
	}

	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		logger.LogWarn("Invalid refresh token", map[string]interface{}{
			"error": err.Error(),
		})
		return utils.UnauthorizedResponse(c, "Invalid refresh token", nil)
	}

	userID := claims["sub"].(string)
	userUUID := claims["user_uuid"].(string)

	tokenErr := handler.authService.CheckStoredRefreshToken(userUUID, refreshToken)
	if tokenErr != nil {
		logger.LogError(tokenErr, "Failed to get stored refresh token", map[string]interface{}{
			"user_id": userID,
			"token":   refreshToken,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	username := claims["user_name"].(string)
	roleCode := claims["role_code"].(string)

	newRefreshToken, err := utils.RegenerateRefreshToken(refreshToken)
	if err != nil {
		logger.LogError(err, "Failed to regenerate refresh token", map[string]interface{}{
			"user_id": userID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	err = handler.authService.UpdateRefreshToken(userUUID, refreshToken, newRefreshToken)
	if err != nil {
		logger.LogError(err, "Failed to update refresh token", map[string]interface{}{
			"user_uuid": userUUID,
		})
		return utils.UnauthorizedResponse(c, "Something went wrong, please try again later", nil)
	}

	// Generate new access token
	newAccessToken, err := utils.GenerateToken(userID, userUUID, username, roleCode)
	if err != nil {
		logger.LogError(err, "Failed to generate access token", map[string]interface{}{
			"user_id": userID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Access token refreshed", map[string]interface{}{
		"reissued_access_token":  newAccessToken,
		"reiussed_refresh_token": newRefreshToken,
	})
}

func (handler *authHandler) AddDeviceToken(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok {
		return utils.UnauthorizedResponse(c, "Token is invalid", nil)
	}

	tokenRequest := new(dto.DeviceTokenRequest)
	if err := c.BodyParser(tokenRequest); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	// Generate FCM Token di backend menggunakan Firebase Admin SDK
	// fcmToken, err := handler.authService.GenerateFCMToken()
	// if err != nil {
	// 	return utils.InternalServerErrorResponse(c, "Failed to generate FCM token", nil)
	// }

	deviceToken := tokenRequest.Token

	// Simpan FCM Token di database
	err := handler.authService.AddDeviceToken(userUUID, deviceToken)
	if err != nil {
		logger.LogError(err, "Failed to save FCM token", map[string]interface{}{
			"user_uuid":    userUUID,
			"device_token": tokenRequest.Token,
		})
		return utils.InternalServerErrorResponse(c, "Failed to save FCM token", nil)
	}

	// Return the generated token as part of the response

	return utils.SuccessResponse(c, "Device token added successfully", nil)
}

func mergeDetails(updateRequestDetails, existingDetails json.RawMessage) (json.RawMessage, error) {
	// Unmarshal `existingDetails` into a map
	var existingDetailsMap map[string]interface{}
	if len(existingDetails) > 0 {
		if err := json.Unmarshal(existingDetails, &existingDetailsMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal existingDetails: %w", err)
		}
	}

	// Unmarshal `updateRequestDetails` into a map
	var updateDetailsMap map[string]interface{}
	if len(updateRequestDetails) > 0 {
		if err := json.Unmarshal(updateRequestDetails, &updateDetailsMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal updateRequestDetails: %w", err)
		}
	}

	// Create the merged map
	mergedDetailsMap := make(map[string]interface{})
	for pascalKey, existingValue := range existingDetailsMap {
		// Convert PascalCase key to snake_case for matching
		snakeKey := pascalToSnakeCase(pascalKey)

		// If there's a value in updateRequestDetails for the snake_case key, use it
		if updateValue, exists := updateDetailsMap[snakeKey]; exists {
			mergedDetailsMap[pascalKey] = updateValue
		} else {
			// Otherwise, keep the existing value
			mergedDetailsMap[pascalKey] = existingValue
		}
	}

	// Marshal the merged map back to json.RawMessage
	mergedDetails, err := json.Marshal(mergedDetailsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged details: %w", err)
	}

	return json.RawMessage(mergedDetails), nil
}

// Utility function to convert PascalCase to snake_case
func pascalToSnakeCase(input string) string {
	var result []rune
	for i, r := range input {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}