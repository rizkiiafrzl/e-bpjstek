package handlers

import (
    "time"

    "farm-management-backend/database"
    "farm-management-backend/middleware"
    "farm-management-backend/models"

    "github.com/gofiber/fiber/v2"
)

type CreateReportPeriodRequest struct {
    Year  int `json:"year"`
    Month int `json:"month"`
}

// ListReportPeriods menampilkan periode pelaporan milik user login
func ListReportPeriods(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user, err := middleware.GetUserFromContext(c)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }

        var periods []models.ReportPeriod
        if err := db.Where("user_id = ?", user.ID).Order("year DESC, month DESC").Find(&periods).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to fetch data"})
        }
        return c.JSON(periods)
    }
}

// CreateReportPeriod membuat periode baru, maksimal 1x per bulan
func CreateReportPeriod(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user, err := middleware.GetUserFromContext(c)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }

        var req CreateReportPeriodRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
        }

        // Default ke bulan & tahun saat ini jika tidak diisi
        now := time.Now()
        if req.Year == 0 {
            req.Year = now.Year()
        }
        if req.Month == 0 {
            req.Month = int(now.Month())
        }
        if req.Month < 1 || req.Month > 12 {
            return c.Status(400).JSON(fiber.Map{"error": "month must be 1..12"})
        }

        // Cek existing
        var count int64
        if err := db.Model(&models.ReportPeriod{}).
            Where("user_id = ? AND year = ? AND month = ?", user.ID, req.Year, req.Month).
            Count(&count).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to check existing"})
        }
        if count > 0 {
            return c.Status(409).JSON(fiber.Map{"error": "Periode bulan ini sudah ada"})
        }

        rp := models.ReportPeriod{
            UserID: user.ID,
            Year:   req.Year,
            Month:  req.Month,
            Status: "Draft",
        }
        if err := db.Create(&rp).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to create"})
        }
        return c.Status(201).JSON(rp)
    }
}


