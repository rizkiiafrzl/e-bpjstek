package handlers

import (
    "time"

    "farm-management-backend/database"
    "farm-management-backend/middleware"
    "farm-management-backend/models"

    "github.com/gofiber/fiber/v2"
)

// List workers milik user login
func ListWorkers(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user, err := middleware.GetUserFromContext(c)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        var workers []models.Worker
        if err := db.Where("user_id = ?", user.ID).Order("created_at desc").Find(&workers).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to fetch workers"})
        }
        return c.JSON(workers)
    }
}

// Get worker by id
func GetWorker(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user, err := middleware.GetUserFromContext(c)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        id := c.Params("id")
        var w models.Worker
        if err := db.Where("id = ? AND user_id = ?", id, user.ID).First(&w).Error; err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "not found"})
        }
        return c.JSON(w)
    }
}

type CreateWorkerRequest struct {
    NIK       string  `json:"nik"`
    KPJ       string  `json:"kpj"`
    NoPegawai string  `json:"noPegawai"`
    Nama      string  `json:"nama"`
    DateOfBirth string `json:"dateOfBirth"` // format: YYYY-MM-DD
    Upah      float64 `json:"upah"`
    Rapel     float64 `json:"rapel"`
}

// Create worker baru
func CreateWorker(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user, err := middleware.GetUserFromContext(c)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        var req CreateWorkerRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
        }
        if req.NIK == "" || req.Nama == "" {
            return c.Status(400).JSON(fiber.Map{"error": "NIK dan Nama wajib diisi"})
        }
        // parse date if provided
        var dob time.Time
        if req.DateOfBirth != "" {
            if t, perr := time.Parse("2006-01-02", req.DateOfBirth); perr == nil {
                dob = t
            }
        }

        w := models.Worker{
            UserID:    user.ID,
            NIK:       req.NIK,
            KPJ:       req.KPJ,
            NoPegawai: req.NoPegawai,
            Nama:      req.Nama,
            DateOfBirth: dob,
            Upah:      req.Upah,
            Rapel:     req.Rapel,
        }
        if err := db.Create(&w).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to create worker"})
        }
        return c.Status(201).JSON(w)
    }
}

// Update worker
func UpdateWorker(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user, err := middleware.GetUserFromContext(c)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        id := c.Params("id")
        var body CreateWorkerRequest
        if err := c.BodyParser(&body); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
        }
        var w models.Worker
        if err := db.Where("id = ? AND user_id = ?", id, user.ID).First(&w).Error; err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "not found"})
        }
        // Apply updates
        w.NIK = body.NIK
        w.KPJ = body.KPJ
        w.NoPegawai = body.NoPegawai
        w.Nama = body.Nama
        if body.DateOfBirth != "" {
            if t, perr := time.Parse("2006-01-02", body.DateOfBirth); perr == nil {
                w.DateOfBirth = t
            }
        }
        w.Upah = body.Upah
        w.Rapel = body.Rapel
        if err := db.Save(&w).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to update"})
        }
        return c.JSON(w)
    }
}

// Delete worker
func DeleteWorker(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user, err := middleware.GetUserFromContext(c)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
        }
        id := c.Params("id")
        if err := db.Where("id = ? AND user_id = ?", id, user.ID).Delete(&models.Worker{}).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to delete"})
        }
        return c.SendStatus(204)
    }
}


