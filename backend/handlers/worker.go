package handlers

import (
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"farm-management-backend/database"
	"farm-management-backend/middleware"
	"farm-management-backend/models"

	"github.com/gofiber/fiber/v2"
	excelize "github.com/xuri/excelize/v2"
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
	NIK                string  `json:"nik"`
	KPJ                string  `json:"kpj"`
	NoPegawai          string  `json:"noPegawai"`
	Nama               string  `json:"nama"`
	DateOfBirth        string  `json:"dateOfBirth"` // format: YYYY-MM-DD
	Upah               float64 `json:"upah"`
	Rapel              float64 `json:"rapel"`
	Nationality        string  `json:"nationality"` // "WNI" (default) atau "WNA"
	PassportNo         string  `json:"passportNo"`
	PassportValidUntil string  `json:"passportValidUntil"` // YYYY-MM-DD
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
		// Validations
		if req.Nama == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Nama wajib diisi"})
		}
		nat := req.Nationality
		if nat == "" {
			nat = "WNI"
		}
		if nat == "WNI" {
			if req.NIK == "" {
				return c.Status(400).JSON(fiber.Map{"error": "NIK wajib diisi untuk WNI"})
			}
			nikRe := regexp.MustCompile(`^\d{16}$`)
			if !nikRe.MatchString(req.NIK) {
				return c.Status(400).JSON(fiber.Map{"error": "Format NIK tidak valid"})
			}
		} else if nat == "WNA" {
			if req.PassportNo == "" {
				return c.Status(400).JSON(fiber.Map{"error": "Nomor paspor wajib untuk WNA"})
			}
		}
		// Uniqueness per user
		var exists int64
		if req.NIK != "" {
			if err := db.Model(&models.Worker{}).Where("user_id = ? AND nik = ?", user.ID, req.NIK).Count(&exists).Error; err == nil && exists > 0 {
				return c.Status(409).JSON(fiber.Map{"error": "NIK sudah terdaftar"})
			}
		}
		// parse date if provided
		var dob time.Time
		if req.DateOfBirth != "" {
			if t, perr := time.Parse("2006-01-02", req.DateOfBirth); perr == nil {
				dob = t
			}
		}

		// Parse passport valid date if provided
		var passValid *time.Time
		if req.PassportValidUntil != "" {
			if t, perr := time.Parse("2006-01-02", req.PassportValidUntil); perr == nil {
				passValid = &t
			}
		}

		w := models.Worker{
			UserID:             user.ID,
			NIK:                req.NIK,
			KPJ:                req.KPJ,
			NoPegawai:          req.NoPegawai,
			Nama:               req.Nama,
			DateOfBirth:        dob,
			Upah:               req.Upah,
			Rapel:              req.Rapel,
			Nationality:        nat,
			PassportNo:         req.PassportNo,
			PassportValidUntil: passValid,
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
		// Apply updates with validations similar to create
		if body.Nama == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Nama wajib diisi"})
		}
		nat := body.Nationality
		if nat == "" {
			nat = w.Nationality
		}
		if nat == "WNI" {
			nikRe := regexp.MustCompile(`^\d{16}$`)
			if body.NIK == "" || !nikRe.MatchString(body.NIK) {
				return c.Status(400).JSON(fiber.Map{"error": "Format NIK tidak valid"})
			}
		} else if nat == "WNA" {
			if body.PassportNo == "" {
				return c.Status(400).JSON(fiber.Map{"error": "Nomor paspor wajib untuk WNA"})
			}
		}
		// Check uniqueness if NIK changed
		if body.NIK != w.NIK {
			var exists int64
			if err := db.Model(&models.Worker{}).Where("user_id = ? AND nik = ?", user.ID, body.NIK).Count(&exists).Error; err == nil && exists > 0 {
				return c.Status(409).JSON(fiber.Map{"error": "NIK sudah terdaftar"})
			}
		}
		w.NIK = body.NIK
		w.KPJ = body.KPJ
		w.NoPegawai = body.NoPegawai
		w.Nama = body.Nama
		w.Nationality = nat
		w.PassportNo = body.PassportNo
		if body.PassportValidUntil != "" {
			if t, perr := time.Parse("2006-01-02", body.PassportValidUntil); perr == nil {
				w.PassportValidUntil = &t
			}
		}
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

// UploadTK menerima file .xlsx untuk upload TK (mendaftar/lanjutan)
func UploadTK(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "file is required"})
		}

		// Validate file extension
		filename := strings.ToLower(fileHeader.Filename)
		if !strings.HasSuffix(filename, ".xlsx") && !strings.HasSuffix(filename, ".xls") {
			return c.Status(400).JSON(fiber.Map{
				"error":    "File harus berformat Excel (.xls atau .xlsx)",
				"received": fileHeader.Filename,
			})
		}

		// Validate file size (max 10MB)
		if fileHeader.Size > 10*1024*1024 {
			return c.Status(400).JSON(fiber.Map{
				"error": "File terlalu besar. Maksimal 10MB",
				"size":  fileHeader.Size,
			})
		}

		// Read file data for database storage
		fileData, err := fileHeader.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to read file"})
		}
		defer fileData.Close()

		// Read all file content
		fileBytes, err := io.ReadAll(fileData)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to read file content"})
		}

		// Save temp file for parsing
		tempPath := "upload_tk_tmp_" + time.Now().Format("20060102150405") + "_" + fileHeader.Filename
		if err := c.SaveFile(fileHeader, tempPath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to save upload"})
		}

		xf, err := excelize.OpenFile(tempPath)
		if err != nil {
			// Clean up temp file
			os.Remove(tempPath)
			return c.Status(400).JSON(fiber.Map{
				"error":   "File Excel tidak valid atau rusak",
				"details": err.Error(),
			})
		}
		defer func() {
			xf.Close()
			// Clean up temp file after processing
			os.Remove(tempPath)
		}()

		sheet := xf.GetSheetName(0)
		if sheet == "" {
			return c.Status(400).JSON(fiber.Map{"error": "File Excel tidak memiliki sheet"})
		}

		rows, err := xf.GetRows(sheet)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Gagal membaca data dari Excel",
				"details": err.Error(),
			})
		}

		if len(rows) < 2 {
			return c.Status(400).JSON(fiber.Map{
				"error": "File Excel harus memiliki minimal 2 baris (header + data)",
				"rows":  len(rows),
			})
		}

		total := 0
		valid := 0
		invalid := 0
		for i, r := range rows {
			if i == 0 {
				continue
			} // header
			if len(r) == 0 {
				continue
			}
			total++
			get := func(idx int) string {
				if idx < len(r) {
					return strings.TrimSpace(r[idx])
				}
				return ""
			}

			// Format TK: NIK, Nama, KPJ, NoPegawai, DateOfBirth, Upah, Rapel, Nationality, PassportNo, PassportValidUntil
			req := CreateWorkerRequest{
				NIK:                get(0),
				Nama:               get(1),
				KPJ:                get(2),
				NoPegawai:          get(3),
				DateOfBirth:        get(4),
				Upah:               parseFloat(get(5)),
				Rapel:              parseFloat(get(6)),
				Nationality:        get(7),
				PassportNo:         get(8),
				PassportValidUntil: get(9),
			}

			if req.Nationality == "" {
				req.Nationality = "WNI"
			}

			// Validate data
			if req.Nama == "" {
				invalid++
				continue
			}

			if req.Nationality == "WNI" {
				nikRe := regexp.MustCompile(`^\d{16}$`)
				if req.NIK == "" || !nikRe.MatchString(req.NIK) {
					invalid++
					continue
				}
				// Check uniqueness
				var exists int64
				db.Model(&models.Worker{}).Where("user_id = ? AND nik = ?", user.ID, req.NIK).Count(&exists)
				if exists > 0 {
					invalid++
					continue
				}
			} else if req.Nationality == "WNA" {
				if req.PassportNo == "" {
					invalid++
					continue
				}
			}

			// Parse dates
			var dob time.Time
			if req.DateOfBirth != "" {
				if t, perr := time.Parse("2006-01-02", req.DateOfBirth); perr == nil {
					dob = t
				}
			}
			var passValid *time.Time
			if req.PassportValidUntil != "" {
				if t, perr := time.Parse("2006-01-02", req.PassportValidUntil); perr == nil {
					passValid = &t
				}
			}

			// Create worker
			w := models.Worker{
				UserID:             user.ID,
				NIK:                req.NIK,
				KPJ:                req.KPJ,
				NoPegawai:          req.NoPegawai,
				Nama:               req.Nama,
				DateOfBirth:        dob,
				Upah:               req.Upah,
				Rapel:              req.Rapel,
				Nationality:        req.Nationality,
				PassportNo:         req.PassportNo,
				PassportValidUntil: passValid,
			}

			if err := db.Create(&w).Error; err != nil {
				invalid++
				continue
			}
			valid++
		}

		// Create history with file data
		hist := models.UploadHistory{
			UserID:           user.ID,
			FileName:         fileHeader.Filename,
			FileData:         fileBytes,
			FileSize:         fileHeader.Size,
			TotalValid:       valid,
			TotalInvalid:     invalid,
			TotalData:        total,
			ValidationStatus: "Selesai",
			DataSource:       "Upload",
			Type:             "tk",
		}
		if err := db.Create(&hist).Error; err != nil {
			log.Printf("Failed to save upload history: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error":   "failed to save history",
				"details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message":   "Upload berhasil",
			"totalData": total,
			"valid":     valid,
			"invalid":   invalid,
			"historyId": hist.ID,
		})
	}
}

// parseFloat helper function
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	// Simple float parsing - in production, use strconv.ParseFloat
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return 0
}

// UploadUpah menerima file .xlsx untuk upload data upah
func UploadUpah(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "file is required"})
		}

		// Validate file extension
		filename := strings.ToLower(fileHeader.Filename)
		if !strings.HasSuffix(filename, ".xlsx") && !strings.HasSuffix(filename, ".xls") {
			return c.Status(400).JSON(fiber.Map{
				"error":    "File harus berformat Excel (.xls atau .xlsx)",
				"received": fileHeader.Filename,
			})
		}

		// Validate file size (max 10MB)
		if fileHeader.Size > 10*1024*1024 {
			return c.Status(400).JSON(fiber.Map{
				"error": "File terlalu besar. Maksimal 10MB",
				"size":  fileHeader.Size,
			})
		}

		// Read file data for database storage
		fileData, err := fileHeader.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to read file"})
		}
		defer fileData.Close()

		// Read all file content
		fileBytes, err := io.ReadAll(fileData)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to read file content"})
		}

		// Save temp file for parsing
		tempPath := "upload_upah_tmp_" + time.Now().Format("20060102150405") + "_" + fileHeader.Filename
		if err := c.SaveFile(fileHeader, tempPath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to save upload"})
		}

		xf, err := excelize.OpenFile(tempPath)
		if err != nil {
			// Clean up temp file
			os.Remove(tempPath)
			return c.Status(400).JSON(fiber.Map{
				"error":   "File Excel tidak valid atau rusak",
				"details": err.Error(),
			})
		}
		defer func() {
			xf.Close()
			// Clean up temp file after processing
			os.Remove(tempPath)
		}()

		sheet := xf.GetSheetName(0)
		if sheet == "" {
			return c.Status(400).JSON(fiber.Map{"error": "File Excel tidak memiliki sheet"})
		}

		rows, err := xf.GetRows(sheet)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Gagal membaca data Excel",
				"details": err.Error(),
			})
		}

		if len(rows) < 2 {
			return c.Status(400).JSON(fiber.Map{"error": "File Excel harus memiliki minimal 2 baris (header + data)"})
		}

		// Skip header row
		dataRows := rows[1:]
		total := len(dataRows)
		valid := 0
		invalid := 0

		// Process each row
		for _, r := range dataRows {
			// Format Upah: NIK, Nama, Upah, Rapel
			get := func(idx int) string {
				if idx < len(r) {
					return strings.TrimSpace(r[idx])
				}
				return ""
			}

			// Validate required fields
			nik := get(0)
			nama := get(1)
			upahStr := get(2)
			rapelStr := get(3)

			if nik == "" || nama == "" {
				invalid++
				continue
			}

			// Parse upah and rapel
			upah := parseFloat(upahStr)
			rapel := parseFloat(rapelStr)

			// Find existing worker by NIK
			var worker models.Worker
			if err := db.Where("user_id = ? AND nik = ?", user.ID, nik).First(&worker).Error; err != nil {
				// Worker not found, skip this row
				invalid++
				continue
			}

			// Update worker's upah and rapel
			worker.Upah = upah
			worker.Rapel = rapel

			if err := db.Save(&worker).Error; err != nil {
				invalid++
				continue
			}
			valid++
		}

		// Create history with file data
		hist := models.UploadHistory{
			UserID:           user.ID,
			FileName:         fileHeader.Filename,
			FileData:         fileBytes,
			FileSize:         fileHeader.Size,
			TotalValid:       valid,
			TotalInvalid:     invalid,
			TotalData:        total,
			ValidationStatus: "Selesai",
			DataSource:       "Upload",
			Type:             "upah",
		}
		if err := db.Create(&hist).Error; err != nil {
			log.Printf("Failed to save upload history: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error":   "failed to save history",
				"details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message":   "Upload upah berhasil",
			"totalData": total,
			"valid":     valid,
			"invalid":   invalid,
			"historyId": hist.ID,
		})
	}
}

// KoreksiTK menerima file .xlsx untuk koreksi data TK yang sudah ada
func KoreksiTK(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "file is required"})
		}

		// Validate file extension
		filename := strings.ToLower(fileHeader.Filename)
		if !strings.HasSuffix(filename, ".xlsx") && !strings.HasSuffix(filename, ".xls") {
			return c.Status(400).JSON(fiber.Map{
				"error":    "File harus berformat Excel (.xls atau .xlsx)",
				"received": fileHeader.Filename,
			})
		}

		// Validate file size (max 10MB)
		if fileHeader.Size > 10*1024*1024 {
			return c.Status(400).JSON(fiber.Map{
				"error": "File terlalu besar. Maksimal 10MB",
				"size":  fileHeader.Size,
			})
		}

		// Read file data for database storage
		fileData, err := fileHeader.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to read file"})
		}
		defer fileData.Close()

		// Read all file content
		fileBytes, err := io.ReadAll(fileData)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to read file content"})
		}

		// Save temp file for parsing
		tempPath := "koreksi_tk_tmp_" + time.Now().Format("20060102150405") + "_" + fileHeader.Filename
		if err := c.SaveFile(fileHeader, tempPath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to save upload"})
		}

		xf, err := excelize.OpenFile(tempPath)
		if err != nil {
			// Clean up temp file
			os.Remove(tempPath)
			return c.Status(400).JSON(fiber.Map{
				"error":   "File Excel tidak valid atau rusak",
				"details": err.Error(),
			})
		}
		defer func() {
			xf.Close()
			// Clean up temp file after processing
			os.Remove(tempPath)
		}()

		sheet := xf.GetSheetName(0)
		if sheet == "" {
			return c.Status(400).JSON(fiber.Map{"error": "File Excel tidak memiliki sheet"})
		}

		rows, err := xf.GetRows(sheet)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Gagal membaca data Excel",
				"details": err.Error(),
			})
		}

		if len(rows) < 2 {
			return c.Status(400).JSON(fiber.Map{"error": "File Excel harus memiliki minimal 2 baris (header + data)"})
		}

		// Skip header row
		dataRows := rows[1:]
		total := len(dataRows)
		valid := 0
		invalid := 0

		// Process each row
		for _, r := range dataRows {
			// Format Koreksi TK: NIK, Nama, KPJ, NoPegawai, DateOfBirth, Upah, Rapel, Nationality, PassportNo, PassportValidUntil
			get := func(idx int) string {
				if idx < len(r) {
					return strings.TrimSpace(r[idx])
				}
				return ""
			}

			// Validate required fields
			nik := get(0)
			nama := get(1)
			kpj := get(2)
			noPegawai := get(3)
			dateOfBirth := get(4)
			upahStr := get(5)
			rapelStr := get(6)
			nationality := get(7)
			passportNo := get(8)
			passportValidUntil := get(9)

			if nik == "" || nama == "" {
				invalid++
				continue
			}

			// Set default nationality
			if nationality == "" {
				nationality = "WNI"
			}

			// Find existing worker by NIK
			var worker models.Worker
			if err := db.Where("user_id = ? AND nik = ?", user.ID, nik).First(&worker).Error; err != nil {
				// Worker not found, skip this row
				invalid++
				continue
			}

			// Update worker data
			worker.Nama = nama
			worker.KPJ = kpj
			worker.NoPegawai = noPegawai
			worker.Upah = parseFloat(upahStr)
			worker.Rapel = parseFloat(rapelStr)
			worker.Nationality = nationality
			worker.PassportNo = passportNo

			// Parse dates
			if dateOfBirth != "" {
				if t, perr := time.Parse("2006-01-02", dateOfBirth); perr == nil {
					worker.DateOfBirth = t
				}
			}

			if passportValidUntil != "" {
				if t, perr := time.Parse("2006-01-02", passportValidUntil); perr == nil {
					worker.PassportValidUntil = &t
				}
			}

			if err := db.Save(&worker).Error; err != nil {
				invalid++
				continue
			}
			valid++
		}

		// Create history with file data
		hist := models.UploadHistory{
			UserID:           user.ID,
			FileName:         fileHeader.Filename,
			FileData:         fileBytes,
			FileSize:         fileHeader.Size,
			TotalValid:       valid,
			TotalInvalid:     invalid,
			TotalData:        total,
			ValidationStatus: "Selesai",
			DataSource:       "Upload",
			Type:             "koreksi",
		}
		if err := db.Create(&hist).Error; err != nil {
			log.Printf("Failed to save upload history: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error":   "failed to save history",
				"details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message":   "Koreksi data TK berhasil",
			"totalData": total,
			"valid":     valid,
			"invalid":   invalid,
			"historyId": hist.ID,
		})
	}
}

// UploadWorkers menerima file .xlsx di field form "file" dan membuat pekerja secara massal
func UploadWorkers(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "file is required"})
		}

		// Validate file extension
		filename := strings.ToLower(fileHeader.Filename)
		if !strings.HasSuffix(filename, ".xlsx") && !strings.HasSuffix(filename, ".xls") {
			return c.Status(400).JSON(fiber.Map{
				"error":    "File harus berformat Excel (.xls atau .xlsx)",
				"received": fileHeader.Filename,
			})
		}

		// Validate file size (max 10MB)
		if fileHeader.Size > 10*1024*1024 {
			return c.Status(400).JSON(fiber.Map{
				"error": "File terlalu besar. Maksimal 10MB",
				"size":  fileHeader.Size,
			})
		}

		// Read file data for database storage
		fileData, err := fileHeader.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to read file"})
		}
		defer fileData.Close()

		// Read all file content
		fileBytes, err := io.ReadAll(fileData)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to read file content"})
		}

		// Save temp file for parsing
		tempPath := "upload_tmp_" + time.Now().Format("20060102150405") + "_" + fileHeader.Filename
		if err := c.SaveFile(fileHeader, tempPath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to save upload"})
		}

		xf, err := excelize.OpenFile(tempPath)
		if err != nil {
			// Clean up temp file
			os.Remove(tempPath)
			return c.Status(400).JSON(fiber.Map{
				"error":   "File Excel tidak valid atau rusak",
				"details": err.Error(),
			})
		}
		defer func() {
			xf.Close()
			// Clean up temp file after processing
			os.Remove(tempPath)
		}()

		sheet := xf.GetSheetName(0)
		if sheet == "" {
			return c.Status(400).JSON(fiber.Map{"error": "File Excel tidak memiliki sheet"})
		}

		rows, err := xf.GetRows(sheet)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "Gagal membaca data dari Excel",
				"details": err.Error(),
			})
		}

		if len(rows) < 2 {
			return c.Status(400).JSON(fiber.Map{
				"error": "File Excel harus memiliki minimal 2 baris (header + data)",
				"rows":  len(rows),
			})
		}

		total := 0
		valid := 0
		invalid := 0
		for i, r := range rows {
			if i == 0 {
				continue
			} // header
			if len(r) == 0 {
				continue
			}
			total++
			get := func(idx int) string {
				if idx < len(r) {
					return strings.TrimSpace(r[idx])
				}
				return ""
			}
			req := CreateWorkerRequest{
				NIK: get(0), Nama: get(1), KPJ: get(2), NoPegawai: get(3), DateOfBirth: get(4),
				Nationality: get(7), PassportNo: get(8), PassportValidUntil: get(9),
			}
			if req.Nationality == "" {
				req.Nationality = "WNI"
			}
			if req.Nationality == "WNI" {
				nikRe := regexp.MustCompile(`^\d{16}$`)
				if req.NIK == "" || !nikRe.MatchString(req.NIK) || req.Nama == "" {
					invalid++
					continue
				}
				var exists int64
				db.Model(&models.Worker{}).Where("user_id = ? AND nik = ?", user.ID, req.NIK).Count(&exists)
				if exists > 0 {
					invalid++
					continue
				}
			} else {
				if req.Nama == "" || req.PassportNo == "" {
					invalid++
					continue
				}
			}
			var dob time.Time
			if req.DateOfBirth != "" {
				if t, perr := time.Parse("2006-01-02", req.DateOfBirth); perr == nil {
					dob = t
				}
			}
			var passValid *time.Time
			if req.PassportValidUntil != "" {
				if t, perr := time.Parse("2006-01-02", req.PassportValidUntil); perr == nil {
					passValid = &t
				}
			}
			w := models.Worker{UserID: user.ID, NIK: req.NIK, KPJ: req.KPJ, NoPegawai: req.NoPegawai, Nama: req.Nama, DateOfBirth: dob, Nationality: req.Nationality, PassportNo: req.PassportNo, PassportValidUntil: passValid}
			if err := db.Create(&w).Error; err != nil {
				invalid++
				continue
			}
			valid++
		}

		// Create history with file data
		hist := models.UploadHistory{
			UserID:           user.ID,
			FileName:         fileHeader.Filename,
			FileData:         fileBytes,
			FileSize:         fileHeader.Size,
			TotalValid:       valid,
			TotalInvalid:     invalid,
			TotalData:        total,
			ValidationStatus: "Selesai",
			DataSource:       "Upload",
			Type:             "workers",
		}
		if err := db.Create(&hist).Error; err != nil {
			log.Printf("Failed to save upload history: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error":   "failed to save history",
				"details": err.Error(),
			})
		}

		return c.JSON(hist)
	}
}

// List upload history
func ListUploadHistory(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		var items []models.UploadHistory
		if err := db.Where("user_id = ?", user.ID).Order("created_at DESC").Find(&items).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to fetch"})
		}
		return c.JSON(items)
	}
}

// DownloadUploadedFile downloads the original Excel file from database
func DownloadUploadedFile(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		id := c.Params("id")
		var hist models.UploadHistory
		if err := db.Where("id = ? AND user_id = ?", id, user.ID).First(&hist).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "upload history not found"})
		}

		if len(hist.FileData) == 0 {
			return c.Status(404).JSON(fiber.Map{"error": "file data not found"})
		}

		// Set headers for file download
		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Set("Content-Disposition", "attachment; filename=\""+hist.FileName+"\"")
		c.Set("Content-Length", string(rune(len(hist.FileData))))

		return c.Send(hist.FileData)
	}
}

// DeleteUploadHistory deletes upload history record
func DeleteUploadHistory(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		id := c.Params("id")
		var hist models.UploadHistory
		if err := db.Where("id = ? AND user_id = ?", id, user.ID).First(&hist).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "upload history not found"})
		}

		if err := db.Delete(&hist).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to delete upload history"})
		}

		return c.SendStatus(204)
	}
}
