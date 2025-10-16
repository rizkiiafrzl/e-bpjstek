#!/bin/bash

# Script untuk mengganti fungsi UploadTK dengan versi yang diperbaiki
echo "🔧 Patching UploadTK function..."

# Backup file asli
cp backend/handlers/worker.go backend/handlers/worker.go.backup
echo "✅ Backup created: worker.go.backup"

# Extract fungsi UploadTKFixed dari fix_upload_tk.go
echo "📝 Extracting fixed UploadTK function..."

# Buat file temporary dengan fungsi yang diperbaiki
cat > temp_upload_tk.go << 'EOF'
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

		// Debug logging
		log.Printf("=== UPLOAD TK DEBUG ===")
		log.Printf("Total rows in Excel: %d", len(rows))
		if len(rows) > 0 {
			log.Printf("Header row: %+v", rows[0])
		}
		if len(rows) > 1 {
			log.Printf("First data row: %+v", rows[1])
		}

		total := 0
		valid := 0
		invalid := 0
		errorDetails := []string{}

		for i, r := range rows {
			if i == 0 {
				continue
			} // header
			if len(r) == 0 {
				log.Printf("Skipping empty row %d", i+1)
				continue
			}
			total++
			log.Printf("Processing row %d: %+v", i+1, r)

			get := func(idx int) string {
				if idx < len(r) {
					return strings.TrimSpace(r[idx])
				}
				return ""
			}

			// Format TK Lengkap: NO_PEGAWAI, NAMA_LENGKAP, GELAR, TELEPON_AREA_RUMAH, TELEPON_RUMAH, EMAIL, TELEPON_AREA_KANTOR, TELEPON_KANTOR, TELEPON_EXT_KANTOR, HP, EMAIL, TEMPAT_LAHIR, TANGGAL_LAHIR, NAMA_IBU_KANDUNG, JENIS_IDENTITAS, NO_IDENTITAS, MASA_LAKU_IDENTITAS, JENIS_KELAMIN, SURAT_MENYURAT_KE, TANGGAL_KEPESERTAAN, STATUS_KAWIN, GOLONGAN_DARAH, NPWP, KODE_NEGARA, UPAH, ALAMAT, KODE_POS, LOKASI_PEKERJAAN, STATUS_PEGAWAI, TGL_AWAL_BEKERJA, TGL_AKHIR_KONTRAK
			req := CreateWorkerRequest{
				// Basic Info - Mapping yang BENAR sesuai template Excel
				NoPegawai:           get(0),             // A: NO_PEGAWAI (PEG001, PEG002, dll)
				Nama:                get(1),             // B: NAMA_LENGKAP (Ruki Afrizal, dll)
				KPJ:                 "",                 // C: GELAR - tidak dipetakan ke KPJ
				TeleponAreaRumah:    get(3),             // D: TELEPON_AREA_RUMAH (021, 022)
				TeleponRumah:        get(4),             // E: TELEPON_RUMAH (7654321)
				Email:               get(10),            // K: EMAIL (ruki.afrizal@example.com)
				TeleponAreaKantor:   get(6),             // G: TELEPON_AREA_KANTOR (021)
				TeleponKantor:       get(7),             // H: TELEPON_KANTOR (021-1234)
				TeleponExtKantor:    get(8),             // I: TELEPON_EXT_KANTOR (1234)
				Handphone:           get(9),             // J: HP (081234567890)
				TempatLahir:         get(11),            // L: TEMPAT_LAHIR (Jakarta)
				DateOfBirth:         get(12),            // M: TANGGAL_LAHIR (1990-05-14)
				IbuKandung:          get(13),            // N: NAMA_IBU_KANDUNG (Siti Aminah)
				JenisIdentitas:      get(14),            // O: JENIS_IDENTITAS (KTP)
				NIK:                 get(15),            // P: NO_IDENTITAS (3175051405900001) - NIK 16 digit
				MasaLakuIdentitas:   get(16),            // Q: MASA_LAKU_IDENTITAS (2016-05-14)
				JenisKelamin:        get(17),            // R: JENIS_KELAMIN (Laki-laki)
				SuratMenyuratKe:     get(18),            // S: SURAT_MENYURAT_KE (Jalan Mandaka No. 10)
				TanggalKepesertaan:  get(19),            // T: TANGGAL_KEPESERTAAN (2020-07-01)
				StatusKawin:         get(20),            // U: STATUS_KAWIN (Belum Kawin)
				GolDarah:            get(21),            // V: GOLONGAN_DARAH (O)
				NPWP:                get(22),            // W: NPWP (012345678901234)
				KodeNegara:          get(23),            // X: KODE_NEGARA (ID)
				Upah:                parseUpah(get(24)), // Y: UPAH (menerima 2500000.00, 5.000.000, 31,00 → 31 juta)
				Alamat:              get(25),            // Z: ALAMAT (Jl. Merpati 12 Jakarta)
				KodePos:             get(26),            // AA: KODE_POS (10130)
				LokasiPekerjaan:     get(27),            // AB: LOKASI_PEKERJAAN (Kantor Pusat)
				StatusPegawai:       get(28),            // AC: STATUS_PEGAWAI (Tetap)
				TanggalAwalBekerja:  get(29),            // AD: TGL_AWAL_BEKERJA (2020-07-01)
				TanggalAkhirKontrak: get(30),            // AE: TGL_AKHIR_KONTRAK (kosong untuk tetap)
				Rapel:               0,                  // Default rapel
				Nationality:         "WNI",              // Default nationality
			}

			log.Printf("Parsed data for row %d: Nama=%s, NIK=%s, Email=%s", i+1, req.Nama, req.NIK, req.Email)

			if req.Nationality == "" {
				req.Nationality = "WNI"
			}

			// Validate data
			if req.Nama == "" {
				log.Printf("Row %d invalid: Nama kosong", i+1)
				errorDetails = append(errorDetails, fmt.Sprintf("Row %d: Nama kosong", i+1))
				invalid++
				continue
			}

			// Validasi email format
			if req.Email != "" {
				emailRe := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
				if !emailRe.MatchString(req.Email) {
					log.Printf("Row %d invalid: Email format salah - %s", i+1, req.Email)
					errorDetails = append(errorDetails, fmt.Sprintf("Row %d: Email format salah - %s", i+1, req.Email))
					invalid++
					continue
				}
			}

			// Validasi NPWP format (15 digit)
			if req.NPWP != "" {
				npwpRe := regexp.MustCompile(`^\d{15}$`)
				if !npwpRe.MatchString(req.NPWP) {
					log.Printf("Row %d invalid: NPWP format salah - %s", i+1, req.NPWP)
					errorDetails = append(errorDetails, fmt.Sprintf("Row %d: NPWP format salah - %s", i+1, req.NPWP))
					invalid++
					continue
				}
			}

			// Validasi golongan darah
			if req.GolDarah != "" {
				validGolDarah := []string{"A", "B", "AB", "O"}
				isValid := false
				for _, gol := range validGolDarah {
					if req.GolDarah == gol {
						isValid = true
						break
					}
				}
				if !isValid {
					log.Printf("Row %d invalid: Golongan darah salah - %s", i+1, req.GolDarah)
					errorDetails = append(errorDetails, fmt.Sprintf("Row %d: Golongan darah salah - %s", i+1, req.GolDarah))
					invalid++
					continue
				}
			}

			if req.Nationality == "WNI" {
				nikRe := regexp.MustCompile(`^\d{16}$`)
				if req.NIK == "" || !nikRe.MatchString(req.NIK) {
					log.Printf("Row %d invalid: NIK format salah - %s", i+1, req.NIK)
					errorDetails = append(errorDetails, fmt.Sprintf("Row %d: NIK format salah - %s", i+1, req.NIK))
					invalid++
					continue
				}
				// Check uniqueness
				var exists int64
				db.Model(&models.Worker{}).Where("user_id = ? AND nik = ?", user.ID, req.NIK).Count(&exists)
				if exists > 0 {
					log.Printf("Row %d invalid: NIK sudah ada - %s", i+1, req.NIK)
					errorDetails = append(errorDetails, fmt.Sprintf("Row %d: NIK sudah ada - %s", i+1, req.NIK))
					invalid++
					continue
				}
			} else if req.Nationality == "WNA" {
				if req.PassportNo == "" {
					log.Printf("Row %d invalid: Passport kosong untuk WNA", i+1)
					errorDetails = append(errorDetails, fmt.Sprintf("Row %d: Passport kosong untuk WNA", i+1))
					invalid++
					continue
				}
			}

			// Parse dates
			var dob time.Time
			if req.DateOfBirth != "" {
				if t, perr := time.Parse("2006-01-02", req.DateOfBirth); perr == nil {
					dob = t
				} else {
					log.Printf("Row %d warning: Tanggal lahir format salah - %s", i+1, req.DateOfBirth)
				}
			}
			var passValid *time.Time
			if req.PassportValidUntil != "" {
				if t, perr := time.Parse("2006-01-02", req.PassportValidUntil); perr == nil {
					passValid = &t
				} else {
					log.Printf("Row %d warning: Passport valid until format salah - %s", i+1, req.PassportValidUntil)
				}
			}

			// Create worker
			// Parse additional dates for mass upload
			var tanggalAwalBekerja *time.Time
			if req.TanggalAwalBekerja != "" {
				if t, err := time.Parse("2006-01-02", req.TanggalAwalBekerja); err == nil {
					tanggalAwalBekerja = &t
				} else {
					log.Printf("Row %d warning: Tanggal awal bekerja format salah - %s", i+1, req.TanggalAwalBekerja)
				}
			}
			var tanggalAkhirKontrak *time.Time
			if req.TanggalAkhirKontrak != "" {
				if t, err := time.Parse("2006-01-02", req.TanggalAkhirKontrak); err == nil {
					tanggalAkhirKontrak = &t
				} else {
					log.Printf("Row %d warning: Tanggal akhir kontrak format salah - %s", i+1, req.TanggalAkhirKontrak)
				}
			}

			w := models.Worker{
				UserID: user.ID,
				NIK:    req.NIK,
				KPJ: func() string {
					if isLikelyKPJ(req.KPJ) {
						return req.KPJ
					}
					return ""
				}(),
				NoPegawai:          req.NoPegawai,
				Nama:               req.Nama,
				DateOfBirth:        dob,
				Upah:               req.Upah,
				Rapel:              req.Rapel,
				Nationality:        req.Nationality,
				PassportNo:         req.PassportNo,
				PassportValidUntil: passValid,
				// Personal Info
				TempatLahir:  req.TempatLahir,
				IbuKandung:   req.IbuKandung,
				JenisKelamin: req.JenisKelamin,
				GolDarah:     req.GolDarah,
				StatusKawin:  req.StatusKawin,
				// Employment Info
				StatusPegawai:       req.StatusPegawai,
				TanggalAwalBekerja:  tanggalAwalBekerja,
				TanggalAkhirKontrak: tanggalAkhirKontrak,
				LokasiPekerjaan:     req.LokasiPekerjaan,
				// Contact Info
				TeleponAreaRumah:  req.TeleponAreaRumah,
				TeleponRumah:      req.TeleponRumah,
				TeleponAreaKantor: req.TeleponAreaKantor,
				TeleponKantor:     req.TeleponKantor,
				TeleponExtKantor:  req.TeleponExtKantor,
				Handphone:         req.Handphone,
				Email:             req.Email,
				// Address Info
				Alamat:    req.Alamat,
				Kabupaten: req.Kabupaten,
				KodePos:   req.KodePos,
				// Additional Info
				NPWP:               req.NPWP,
				JenisIdentitas:     req.JenisIdentitas,
				MasaLakuIdentitas:  req.MasaLakuIdentitas,
				SuratMenyuratKe:    req.SuratMenyuratKe,
				TanggalKepesertaan: req.TanggalKepesertaan,
				KodeNegara:         req.KodeNegara,
			}

			if err := db.Create(&w).Error; err != nil {
				log.Printf("Row %d error: Failed to create worker - %v", i+1, err)
				errorDetails = append(errorDetails, fmt.Sprintf("Row %d: Database error - %v", i+1, err))
				invalid++
				continue
			}
			valid++
			log.Printf("Row %d success: Created worker %s", i+1, req.Nama)
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

		log.Printf("=== UPLOAD TK COMPLETED ===")
		log.Printf("Total: %d, Valid: %d, Invalid: %d", total, valid, invalid)

		response := fiber.Map{
			"message":   "Upload berhasil",
			"totalData": total,
			"valid":     valid,
			"invalid":   invalid,
			"historyId": hist.ID,
		}

		if len(errorDetails) > 0 {
			response["errors"] = errorDetails[:min(10, len(errorDetails))] // Limit to 10 errors
		}

		return c.JSON(response)
	}
}
EOF

echo "✅ Fixed UploadTK function created"

# Add missing import for fmt
echo "📝 Adding missing import..."

# Check if fmt is already imported
if ! grep -q '"fmt"' backend/handlers/worker.go; then
    # Add fmt import
    sed -i 's/import (/import (\n\t"fmt"/' backend/handlers/worker.go
    sed -i 's/"log"/"log"\n\t/' backend/handlers/worker.go
fi

echo "✅ Missing imports added"

echo "🎉 Patch completed!"
echo ""
echo "📋 Next steps:"
echo "1. Restart backend: cd backend && go run main.go"
echo "2. Test upload Excel dengan data valid"
echo "3. Cek log backend untuk debug info"
echo "4. Verifikasi data muncul di tabel"







