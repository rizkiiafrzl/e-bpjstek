package main

import (
    "farm-management-backend/config"
    "farm-management-backend/database"
    "farm-management-backend/handlers"
    "farm-management-backend/middleware"
    "log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

    // Initialize database
    db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate database
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetEnv("CORS_ORIGIN", "http://localhost:5173"),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

    // Inject db into context for middleware access
    app.Use(func(c *fiber.Ctx) error {
        c.Locals("db", db)
        return c.Next()
    })

    // Routes
    setupRoutes(app, db)

	// Start server
	port := config.GetEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func setupRoutes(app *fiber.App, db *database.DB) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Farm Management API is running",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Auth routes (no auth required)
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register(db))
	auth.Post("/login", handlers.Login(db))
	auth.Post("/google", handlers.GoogleAuth(db))

	// Protected routes
	protected := api.Group("/", middleware.AuthRequired())
	protected.Get("/profile", handlers.GetProfile(db))
	protected.Put("/profile", handlers.UpdateProfile(db))
	protected.Post("/logout", handlers.Logout())

    // (farm routes removed)
}
