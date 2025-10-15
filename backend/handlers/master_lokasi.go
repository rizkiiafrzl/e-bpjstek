package handlers

import (
    "farm-management-backend/database"
    "farm-management-backend/models"

    "github.com/gofiber/fiber/v2"
)

// ListMasterLokasi menampilkan semua master lokasi
func ListMasterLokasi(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var lokasi []models.MasterLokasi
        if err := db.Where("is_active = true").Order("nama").Find(&lokasi).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to fetch lokasi"})
        }
        return c.JSON(lokasi)
    }
}

// CreateMasterLokasi membuat master lokasi baru
func CreateMasterLokasi(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var lokasi models.MasterLokasi
        if err := c.BodyParser(&lokasi); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
        }

        if err := db.Create(&lokasi).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to create lokasi"})
        }

        return c.Status(201).JSON(lokasi)
    }
}

// UpdateMasterLokasi mengupdate master lokasi
func UpdateMasterLokasi(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id := c.Params("id")
        var lokasi models.MasterLokasi
        if err := db.Where("id = ?", id).First(&lokasi).Error; err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "lokasi not found"})
        }

        if err := c.BodyParser(&lokasi); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
        }

        if err := db.Save(&lokasi).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to update lokasi"})
        }

        return c.JSON(lokasi)
    }
}

// DeleteMasterLokasi menghapus master lokasi (soft delete)
func DeleteMasterLokasi(db *database.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        id := c.Params("id")
        var lokasi models.MasterLokasi
        if err := db.Where("id = ?", id).First(&lokasi).Error; err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "lokasi not found"})
        }

        // Soft delete - set is_active = false
        lokasi.IsActive = false
        if err := db.Save(&lokasi).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "failed to delete lokasi"})
        }

        return c.SendStatus(204)
    }
}
