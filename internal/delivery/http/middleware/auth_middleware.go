package middleware

import (
	"shwetaik-sqlacc-stock-api/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	if claims.Service != "stock-api" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	c.Locals("service", claims.Service)
	return c.Next()
}
