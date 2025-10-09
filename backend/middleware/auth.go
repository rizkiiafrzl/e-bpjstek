package middleware

import (
	"farm-management-backend/config"
	"farm-management-backend/database"
	"farm-management-backend/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// AuthRequired middleware to protect routes
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(401, "Invalid signing method")
			}
			return []byte(config.GetEnv("JWT_SECRET", "your-secret-key")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Extract user ID from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

    userID, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid user ID in token",
			})
		}

		// Check if token is expired
		exp, ok := claims["exp"].(float64)
		if !ok || float64(time.Now().Unix()) > exp {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token expired",
			})
		}

		// Get user from database
		db := c.Locals("db").(*database.DB)
		var user models.User
    if err := db.First(&user, "id = ?", userID).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Store user in context
		c.Locals("user", user)
    c.Locals("user_id", userID)

		return c.Next()
	}
}

// GetUserFromContext extracts user from context
func GetUserFromContext(c *fiber.Ctx) (*models.User, error) {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return nil, fiber.NewError(401, "User not found in context")
	}
	return &user, nil
}
