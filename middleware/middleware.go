package middleware

import (
	"shuttle/logger"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
)

func SchoolAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userId").(string)
		if !ok || userID == "" {
			return utils.UnauthorizedResponse(c, "User ID is missing or invalid", nil)
		}

		schoolID, err := services.CheckPermittedSchoolAccess(userID)
		if err != nil {
			return utils.ForbiddenResponse(c, "You don't have permission to any school, please contact the support team", nil)
		}

		c.Locals("schoolId", schoolID)

		return c.Next()
	}
}

func AuthenticationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return utils.UnauthorizedResponse(c, "Missing token", nil)
		}
        
		const bearerPrefix = "Bearer "
		if len(token) > len(bearerPrefix) && token[:len(bearerPrefix)] == bearerPrefix {
			token = token[len(bearerPrefix):]
		}

		_, exists := utils.InvalidTokens[token]
		if exists {
			return utils.UnauthorizedResponse(c, "Invalid token or you have been logged out", nil)
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			logger.LogWarn("Invalid token", map[string]interface{}{})
			return utils.UnauthorizedResponse(c, "Invalid token", nil)
		}

		userID, ok := claims["userId"].(string)
		if !ok || userID == "" {
			return utils.UnauthorizedResponse(c, "User ID is missing or invalid", nil)
		}

		roleCode, ok := claims["role_code"].(string)
		if !ok || roleCode == "" {
			return utils.UnauthorizedResponse(c, "Role code is missing or invalid", nil)
		}

		username, ok := claims["name"].(string)
		if !ok || username == "" {
			return utils.UnauthorizedResponse(c, "Username is missing or invalid", nil)
		}

		c.Locals("userId", userID)
		c.Locals("role_code", roleCode)
		c.Locals("username", username)

		return c.Next()
	}
}

func AuthorizationMiddleware(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role_code").(string)
		if !ok || role == "" {
			return utils.UnauthorizedResponse(c, "Role code is missing or invalid", nil)
		}

		if !contains(allowedRoles, role) {
			return utils.ForbiddenResponse(c, "You don't have permission to access this resource", nil)
		}

		return c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
