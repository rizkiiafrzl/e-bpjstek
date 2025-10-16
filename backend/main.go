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

	// Seed initial data
	if err := database.SeedTemplates(db); err != nil {
		log.Printf("Failed to seed templates: %v", err)
	}

	// Start server
	port := config.GetEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func setupRoutes(app *fiber.App, db *database.DB) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Farm Management API is running",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Public templates download
	templates := api.Group("/templates")
	templates.Get(":key", handlers.DownloadTemplate())

	// Auth routes (no auth required)
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register(db))
	auth.Post("/login", handlers.Login(db))
	auth.Post("/google", handlers.GoogleAuth(db))
	// SSO Google
	auth.Get("/google/login", handlers.GoogleLoginRedirect())
	auth.Get("/google/callback", handlers.GoogleLoginCallback(db))

	// Protected routes (apply middleware only to these routes)
	api.Get("/profile", middleware.AuthRequired(), handlers.GetProfile(db))
	api.Put("/profile", middleware.AuthRequired(), handlers.UpdateProfile(db))
	api.Post("/logout", middleware.AuthRequired(), handlers.Logout())

	// Report periods (protected)
	api.Get("/report-periods", middleware.AuthRequired(), handlers.ListReportPeriods(db))
	api.Post("/report-periods", middleware.AuthRequired(), handlers.CreateReportPeriod(db))
	api.Delete("/report-periods/:id", middleware.AuthRequired(), handlers.DeleteReportPeriod(db))
	api.Post("/report-periods/:id/calculate", middleware.AuthRequired(), handlers.CalculateReportPeriod(db))
	api.Post("/report-periods/:id/finalize", middleware.AuthRequired(), handlers.FinalizeReportPeriod(db))
	api.Get("/report-periods/:id/print", middleware.AuthRequired(), handlers.PrintReportPeriod(db))
	api.Get("/report-periods-summary", middleware.AuthRequired(), handlers.SummaryDashboard(db))

	// Workers (protected)
	// Specific routes must come BEFORE parameter routes to avoid conflicts
	api.Post("/workers/upload", middleware.AuthRequired(), handlers.UploadWorkers(db))
	api.Post("/workers/upload-tk", middleware.AuthRequired(), handlers.UploadTK(db))
	api.Post("/workers/upload-upah", middleware.AuthRequired(), handlers.UploadUpah(db))
	api.Post("/koreksi-data-tk", middleware.AuthRequired(), handlers.KoreksiTK(db))
	api.Get("/workers/upload-history", middleware.AuthRequired(), handlers.ListUploadHistory(db))
	api.Get("/workers/upload-history/:id/download", middleware.AuthRequired(), handlers.DownloadUploadedFile(db))
	api.Delete("/workers/upload-history/:id", middleware.AuthRequired(), handlers.DeleteUploadHistory(db))
	// General routes
	api.Get("/workers", middleware.AuthRequired(), handlers.ListWorkers(db))
	api.Post("/workers", middleware.AuthRequired(), handlers.CreateWorker(db))
	// Parameter routes must come LAST to avoid catching specific routes
	api.Get("/workers/:id", middleware.AuthRequired(), handlers.GetWorker(db))
	api.Put("/workers/:id", middleware.AuthRequired(), handlers.UpdateWorker(db))
	api.Delete("/workers/:id", middleware.AuthRequired(), handlers.DeleteWorker(db))

	// (farm routes removed)
}
