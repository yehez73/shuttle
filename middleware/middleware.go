package middleware

import (
	"shuttle/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		_, exists := utils.InvalidTokens[token]
		if exists {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token or you have logged out", nil)
		}

		if token == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing or invalid token", nil)
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
		}

		c.Locals("userId", claims["userId"])
		return c.Next()
	}
}

func SuperAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		_, exists := utils.InvalidTokens[token]
		if exists {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token or you have logged out", nil)
		}

		if token == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing or invalid token", nil)
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
		}

		if claims["role"] != "superadmin" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to access this resource", nil)
		}

		c.Locals("userId", claims["userId"])
		return c.Next()
	}
}

func SchoolAdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		_, exists := utils.InvalidTokens[token]
		if exists {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token or you have logged out", nil)
		}

		if token == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing or invalid token", nil)
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", nil)
		}

		if claims["role"] != "schooladmin" && claims["role"] != "superadmin" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to access this resource", nil)
		}

		c.Locals("userId", claims["userId"])
		return c.Next()
	}
}