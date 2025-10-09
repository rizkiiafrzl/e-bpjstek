package handlers

import (
	"farm-management-backend/config"
	"farm-management-backend/database"
	"farm-management-backend/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents registration request
type RegisterRequest struct {
    FullName string `json:"full_name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Register handles user registration
func Register(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RegisterRequest
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

		// Check if user already exists
		var existingUser models.User
		if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return c.Status(409).JSON(fiber.Map{
				"error": "User with this email already exists",
			})
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to hash password",
			})
		}

        // Create user (matches existing columns)
        user := models.User{
            FullName:     req.FullName,
            Email:        req.Email,
            PasswordHash: string(hashedPassword),
        }

		if err := db.Create(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}

		// Generate JWT token
        token, err := generateJWT(user.ID.String())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

        // Hide password in response
        user.PasswordHash = ""

		return c.Status(201).JSON(AuthResponse{
			Token: token,
			User:  user,
		})
	}
}

// Login handles user login
func Login(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
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

		// Find user
		var user models.User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}

        // (no is_active column in your users table)

		// Verify password
        if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}

		// Generate JWT token
        token, err := generateJWT(user.ID.String())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

        // Hide password in response
        user.PasswordHash = ""

		return c.JSON(AuthResponse{
			Token: token,
			User:  user,
		})
	}
}

// GoogleAuth handles Google OAuth (placeholder)
func GoogleAuth(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement Google OAuth
		return c.Status(501).JSON(fiber.Map{
			"error": "Google OAuth not implemented yet",
		})
	}
}

// Logout handles user logout
func Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// In JWT, logout is handled on client side by removing token
		return c.JSON(fiber.Map{
			"message": "Logged out successfully",
		})
	}
}

// generateJWT generates JWT token for user
func generateJWT(userID interface{}) (string, error) {
	expireHours, _ := strconv.Atoi(config.GetEnv("JWT_EXPIRE_HOURS", "24"))
	expireTime := time.Now().Add(time.Duration(expireHours) * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expireTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetEnv("JWT_SECRET", "your-secret-key")))
}
