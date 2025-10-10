package handlers

import (
	"fmt"
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

		// Optional filters & pagination
		status := c.Query("status")
		page := c.QueryInt("page", 1)
		pageSize := c.QueryInt("pageSize", 10)
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 10
		}

		var periods []models.ReportPeriod
		q := db.Where("user_id = ?", user.ID)
		if status != "" && status != "all" {
			q = q.Where("status = ?", status)
		}
		var total int64
		if err := q.Model(&models.ReportPeriod{}).Count(&total).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to count"})
		}
		if err := q.Order("year DESC, month DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&periods).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to fetch data"})
		}
		return c.JSON(fiber.Map{
			"items":    periods,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		})
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

// DeleteReportPeriod menghapus periode milik user (hanya bila status Draft)
func DeleteReportPeriod(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		id := c.Params("id")
		var rp models.ReportPeriod
		if err := db.Where("id = ? AND user_id = ?", id, user.ID).First(&rp).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		}
		if rp.Status != "Draft" {
			return c.Status(400).JSON(fiber.Map{"error": "hanya status Draft yang bisa dihapus"})
		}
		if err := db.Delete(&rp).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to delete"})
		}
		return c.SendStatus(204)
	}
}

// CalculateReportPeriod melakukan perhitungan dummy total iuran & denda
func CalculateReportPeriod(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		id := c.Params("id")
		var rp models.ReportPeriod
		if err := db.Where("id = ? AND user_id = ?", id, user.ID).First(&rp).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		}
		// Dummy rule: total_iuran = jumlah TK * 29700, denda 2%
		// Dalam implementasi nyata, hitung berdasarkan data worker di bulan tsb
		if rp.TotalTk < 0 {
			rp.TotalTk = 0
		}
		rp.TotalIuran = float64(rp.TotalTk) * 29700
		rp.TotalDenda = rp.TotalIuran * 0.02
		if err := db.Save(&rp).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to save"})
		}
		return c.JSON(rp)
	}
}

// FinalizeReportPeriod mengubah status menjadi Finalisasi
func FinalizeReportPeriod(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		id := c.Params("id")
		var rp models.ReportPeriod
		if err := db.Where("id = ? AND user_id = ?", id, user.ID).First(&rp).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		}
		if rp.Status == "Posting" {
			return c.Status(400).JSON(fiber.Map{"error": "sudah Posting"})
		}
		rp.Status = "Finalisasi"
		if err := db.Save(&rp).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to save"})
		}
		return c.JSON(rp)
	}
}

// PrintReportPeriod placeholder: kembalikan JSON yang nanti bisa dipakai client untuk cetak
func PrintReportPeriod(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		id := c.Params("id")
		var rp models.ReportPeriod
		if err := db.Where("id = ? AND user_id = ?", id, user.ID).First(&rp).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		}
		return c.JSON(fiber.Map{
			"message": "print-ready",
			"period":  rp,
		})
	}
}

// SummaryDashboard mengembalikan ringkasan untuk bagian atas dashboard
func SummaryDashboard(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		var rp models.ReportPeriod
		found := (db.Where("user_id = ?", user.ID).Order("year DESC, month DESC").First(&rp).Error == nil)

		uid := user.ID.String()
		tail := uid
		if len(uid) >= 4 {
			tail = uid[len(uid)-4:]
		}
		yyyymm := 0
		if found {
			yyyymm = rp.Year*100 + rp.Month
		} else {
			now := time.Now()
			yyyymm = now.Year()*100 + int(now.Month())
		}
		kode := fmt.Sprintf("4241%s%04d", tail, yyyymm)

		totalIuranDanDenda := 0.0
		if found {
			totalIuranDanDenda = rp.TotalIuran + rp.TotalDenda
		}
		sisaSebelumnya := 0.0
		totalTagihan := totalIuranDanDenda + sisaSebelumnya

		return c.JSON(fiber.Map{
			"kodeTagihan":              kode,
			"totalIuranDanDenda":       totalIuranDanDenda,
			"sisaPembayaranSebelumnya": sisaSebelumnya,
			"totalTagihan":             totalTagihan,
		})
	}
}
