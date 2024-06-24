package middleware

import (
	"strings"

	"github.com/ajay/exptest/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Skip authentication for specific endpoints
	if c.Path() == "/register" {
		return c.Next()
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	userIDHex, ok := claims["userID"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Set user ID in context for further use
	c.Locals("userID", userID)

	return c.Next()
}
