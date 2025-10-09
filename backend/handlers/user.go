package handlers

import (
    "farm-management-backend/database"
    "farm-management-backend/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// GetProfile gets user profile
func GetProfile(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

        // Hide password
        user.PasswordHash = ""

		return c.JSON(user)
	}
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"omitempty,min=2,max=100"`
	Avatar   string `json:"avatar" validate:"omitempty,url"`
}

// UpdateProfile updates user profile
func UpdateProfile(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var req UpdateProfileRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Validation failed",
				"details": err.Error(),
			})
		}

		// Update user fields
		updates := make(map[string]interface{})
		if req.FullName != "" {
			updates["full_name"] = req.FullName
		}
		if req.Avatar != "" {
			updates["avatar"] = req.Avatar
		}

		// Update in database
		if err := db.Model(user).Updates(updates).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to update profile",
			})
		}

		// Get updated user
		if err := db.First(user, user.ID).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to fetch updated profile",
			})
		}

        // Hide password
        user.PasswordHash = ""

		return c.JSON(user)
	}
}
