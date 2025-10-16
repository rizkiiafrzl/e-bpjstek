package handlers

import (
	"fmt"
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

		// Optional filter parameter
		status := c.Query("status")

		var workers []models.Worker
		q := db.Where("user_id = ?", user.ID)

		// Filter berdasarkan status
		if status == "aktif" {
			// Hanya tampilkan worker yang aktif (bukan NONAKTIF)
			q = q.Where("status_pegawai != ? OR status_pegawai IS NULL", "NONAKTIF")
		} else if status == "nonaktif" {
			// Hanya tampilkan worker yang nonaktif
			q = q.Where("status_pegawai = ?", "NONAKTIF")
		}
		// Jika status = "semua" atau tidak ada parameter, tampilkan semua data

		if err := q.Order("created_at desc").Find(&workers).Error; err != nil {
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

// CreateWorkerRequest sesuai spesifikasi SIPP BPJamsostek
type CreateWorkerRequest struct {
	// A. Data Identitas Dasar (WAJIB)
	NoPegawai string `json:"noPegawai" validate:"required"`  // A: NO_PEGAWAI
	Nama      string `json:"nama" validate:"required,min=3"` // B: NAMA_LENGKAP
	Gelar     string `json:"gelar"`                          // C: GELAR

	// D-E. Kontak Rumah
	TeleponAreaRumah string `json:"teleponAreaRumah"` // D: TELEPON_AREA_RUMAH
	TeleponRumah     string `json:"teleponRumah"`     // E: TELEPON_RUMAH

	// G-I. Kontak Kantor
	TeleponAreaKantor string `json:"teleponAreaKantor"` // G: TELEPON_AREA_KANTOR
	TeleponKantor     string `json:"teleponKantor"`     // H: TELEPON_KANTOR
	TeleponExtKantor  string `json:"teleponExtKantor"`  // I: TELEPON_EXT_KANTOR

	// J-K. Kontak Utama
	Handphone string `json:"handphone" validate:"required"`   // J: HP
	Email     string `json:"email" validate:"required,email"` // K: EMAIL

	// L-M. Data Kelahiran
	TempatLahir string `json:"tempatLahir"` // L: TEMPAT_LAHIR
	DateOfBirth string `json:"dateOfBirth"` // M: TANGGAL_LAHIR

	// N. Data Keluarga
	IbuKandung string `json:"ibuKandung"` // N: NAMA_IBU_KANDUNG

	// O-Q. Identitas
	JenisIdentitas    string `json:"jenisIdentitas"`                 // O: JENIS_IDENTITAS
	NIK               string `json:"nik" validate:"required,len=16"` // P: NO_IDENTITAS (NIK)
	MasaLakuIdentitas string `json:"masaLakuIdentitas"`              // Q: MASA_LAKU_IDENTITAS

	// R. Jenis Kelamin
	JenisKelamin string `json:"jenisKelamin" validate:"required,oneof=Laki-laki Perempuan"` // R: JENIS_KELAMIN

	// S-T. Kepesertaan
	SuratMenyuratKe    string `json:"suratMenyuratKe"`    // S: SURAT_MENYURAT_KE
	TanggalKepesertaan string `json:"tanggalKepesertaan"` // T: TANGGAL_KEPESERTAAN

	// U-V. Status Personal
	StatusKawin string `json:"statusKawin"`                                  // U: STATUS_KAWIN
	GolDarah    string `json:"golDarah" validate:"omitempty,oneof=A B AB O"` // V: GOLONGAN_DARAH

	// W-X. Identitas Pajak & Negara
	NPWP       string `json:"npwp" validate:"omitempty,len=15"` // W: NPWP
	KodeNegara string `json:"kodeNegara"`                       // X: KODE_NEGARA

	// Y. Upah/Gaji
	Upah float64 `json:"upah" validate:"required,min=1000000"` // Y: UPAH (minimal 1 juta)

	// Z-AA. Alamat
	Alamat  string `json:"alamat"`  // Z: ALAMAT
	KodePos string `json:"kodePos"` // AA: KODE_POS

	// AB-AC. Pekerjaan
	LokasiPekerjaan string `json:"lokasiPekerjaan"`                                    // AB: LOKASI_PEKERJAAN
	StatusPegawai   string `json:"statusPegawai" validate:"required,oneof=PKWTT PKWT"` // AC: STATUS_PEGAWAI

	// AD-AE. Tanggal Kerja
	TanggalAwalBekerja  string `json:"tanggalAwalBekerja" validate:"required"` // AD: TGL_AWAL_BEKERJA
	TanggalAkhirKontrak string `json:"tanggalAkhirKontrak"`                    // AE: TGL_AKHIR_KONTRAK (wajib untuk PKWT)

	// AF-AG. Data Tambahan
	Rapel       float64 `json:"rapel"`       // AF: RAPEL
	Nationality string  `json:"nationality"` // AG: NATIONALITY

	// Data untuk WNA (jika bukan WNI)
	PassportNo         string `json:"passportNo"`
	PassportValidUntil string `json:"passportValidUntil"`

	// Data tambahan untuk backward compatibility
	KPJ       string `json:"kpj"`
	Kabupaten string `json:"kabupaten"`
}

// ProgramBPJSRequest untuk mengelola program kepesertaan
type ProgramBPJSRequest struct {
	JKK bool `json:"jkk"` // Jaminan Kecelakaan Kerja
	JKM bool `json:"jkm"` // Jaminan Kematian
	JHT bool `json:"jht"` // Jaminan Hari Tua
	JP  bool `json:"jp"`  // Jaminan Pensiun
}

// IuranCalculation untuk kalkulasi iuran otomatis
type IuranCalculation struct {
	Upah        float64            `json:"upah"`
	RisikoUsaha float64            `json:"risikoUsaha"` // 0.24% - 1.74% tergantung risiko
	Program     ProgramBPJSRequest `json:"program"`

	// Hasil Kalkulasi
	IuranJKK           float64 `json:"iuranJKK"`           // Ditanggung perusahaan 100%
	IuranJKM           float64 `json:"iuranJKM"`           // Ditanggung perusahaan 100%
	IuranJHTPerusahaan float64 `json:"iuranJHTPerusahaan"` // 3.7% dari upah
	IuranJHTTK         float64 `json:"iuranJHTTK"`         // 2% dari upah
	IuranJPPerusahaan  float64 `json:"iuranJPPerusahaan"`  // 2% dari upah
	IuranJPTK          float64 `json:"iuranJPTK"`          // 1% dari upah

	TotalDitanggungPerusahaan float64 `json:"totalDitanggungPerusahaan"`
	TotalDitanggungTK         float64 `json:"totalDitanggungTK"`
	TotalIuran                float64 `json:"totalIuran"`
}

// ValidasiSIPP dengan validasi lengkap sesuai SIPP BPJamsostek
func ValidasiSIPP(req *CreateWorkerRequest) []string {
	var errors []string

	// 1. Validasi Data Wajib
	if req.Nama == "" {
		errors = append(errors, "Nama lengkap wajib diisi")
	}
	if req.NIK == "" {
		errors = append(errors, "NIK wajib diisi")
	}
	if req.Email == "" {
		errors = append(errors, "Email wajib diisi")
	}
	if req.Handphone == "" {
		errors = append(errors, "Nomor HP wajib diisi")
	}
	if req.StatusPegawai == "" {
		errors = append(errors, "Status pegawai wajib diisi")
	}
	if req.TanggalAwalBekerja == "" {
		errors = append(errors, "Tanggal awal bekerja wajib diisi")
	}
	if req.Upah < 1000000 {
		errors = append(errors, "Upah minimal Rp 1.000.000")
	}

	// 2. Validasi Format Data
	if req.NIK != "" {
		nikRe := regexp.MustCompile(`^\d{16}$`)
		if !nikRe.MatchString(req.NIK) {
			errors = append(errors, "NIK harus 16 digit angka")
		}
	}

	if req.Email != "" {
		emailRe := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRe.MatchString(req.Email) {
			errors = append(errors, "Format email tidak valid")
		}
	}

	if req.NPWP != "" {
		npwpRe := regexp.MustCompile(`^\d{15}$`)
		if !npwpRe.MatchString(req.NPWP) {
			errors = append(errors, "NPWP harus 15 digit angka")
		}
	}

	// 3. Validasi Jenis Kelamin
	if req.JenisKelamin != "" && req.JenisKelamin != "Laki-laki" && req.JenisKelamin != "Perempuan" {
		errors = append(errors, "Jenis kelamin harus 'Laki-laki' atau 'Perempuan'")
	}

	// 4. Validasi Golongan Darah
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
			errors = append(errors, "Golongan darah harus A, B, AB, atau O")
		}
	}

	// 5. Validasi Status Pegawai dan Tanggal Kontrak
	if req.StatusPegawai == "PKWT" {
		if req.TanggalAkhirKontrak == "" {
			errors = append(errors, "Tanggal akhir kontrak wajib diisi untuk PKWT")
		} else {
			// Validasi tanggal kontrak
			tglAwal, err1 := time.Parse("2006-01-02", req.TanggalAwalBekerja)
			tglAkhir, err2 := time.Parse("2006-01-02", req.TanggalAkhirKontrak)

			if err1 == nil && err2 == nil {
				if tglAkhir.Before(tglAwal) || tglAkhir.Equal(tglAwal) {
					errors = append(errors, "Tanggal akhir kontrak harus setelah tanggal awal bekerja")
				}

				// Validasi maksimal kontrak 2 tahun
				duration := tglAkhir.Sub(tglAwal)
				if duration.Hours() > 24*365*2 {
					errors = append(errors, "Kontrak PKWT maksimal 2 tahun")
				}
			}
		}
	}

	// 6. Validasi Usia Kerja
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			age := time.Since(dob).Hours() / (24 * 365)
			if age < 17 {
				errors = append(errors, "Usia minimal 17 tahun")
			}
			if age > 65 {
				errors = append(errors, "Usia maksimal 65 tahun")
			}
		}
	}

	// 7. Validasi WNA
	if req.Nationality == "WNA" {
		if req.PassportNo == "" {
			errors = append(errors, "Nomor passport wajib untuk WNA")
		}
		if req.PassportValidUntil == "" {
			errors = append(errors, "Tanggal berlaku passport wajib untuk WNA")
		}
	}

	return errors
}

// KalkulasiIuran menghitung iuran BPJS sesuai spesifikasi SIPP
func KalkulasiIuran(upah float64, program ProgramBPJSRequest, risikoUsaha float64) IuranCalculation {
	calc := IuranCalculation{
		Upah:        upah,
		RisikoUsaha: risikoUsaha,
		Program:     program,
	}

	// Default risiko usaha jika tidak diisi (kategori rendah)
	if risikoUsaha == 0 {
		risikoUsaha = 0.54 // 0.54% untuk kategori risiko rendah
	}

	// Kalkulasi per program
	if program.JKK {
		calc.IuranJKK = upah * (risikoUsaha / 100) // Ditanggung perusahaan 100%
	}

	if program.JKM {
		calc.IuranJKM = upah * 0.003 // 0.3% dari upah, ditanggung perusahaan 100%
	}

	if program.JHT {
		calc.IuranJHTPerusahaan = upah * 0.037 // 3.7% dari upah
		calc.IuranJHTTK = upah * 0.02          // 2% dari upah
	}

	if program.JP {
		calc.IuranJPPerusahaan = upah * 0.02 // 2% dari upah
		calc.IuranJPTK = upah * 0.01         // 1% dari upah
	}

	// Total per pihak
	calc.TotalDitanggungPerusahaan = calc.IuranJKK + calc.IuranJKM + calc.IuranJHTPerusahaan + calc.IuranJPPerusahaan
	calc.TotalDitanggungTK = calc.IuranJHTTK + calc.IuranJPTK
	calc.TotalIuran = calc.TotalDitanggungPerusahaan + calc.TotalDitanggungTK

	return calc
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

		// Set default values
		if req.Nationality == "" {
			req.Nationality = "WNI"
		}

		// Validasi: dukung pembuatan draft minimal untuk alur "lanjutan"
		// Jika data belum lengkap (misal hanya NIK/Nama, upah 0), lewati ValidasiSIPP penuh
		// Minimal checks: Nama wajib, identitas dasar (WNI -> NIK 16 digit jika diisi; WNA -> passport jika diisi)
		isDraft := false
		if req.Upah == 0 || req.StatusPegawai == "" || req.TanggalAwalBekerja == "" || req.Email == "" || req.Handphone == "" {
			isDraft = true
		}

		if !isDraft {
			// Validasi lengkap bila bukan draft
			validationErrors := ValidasiSIPP(&req)
			if len(validationErrors) > 0 {
				return c.Status(400).JSON(fiber.Map{
					"error":   "validation failed",
					"details": validationErrors,
				})
			}
		} else {
			// Validasi minimal untuk draft
			var minimalErrors []string
			if strings.TrimSpace(req.Nama) == "" {
				minimalErrors = append(minimalErrors, "Nama lengkap wajib diisi")
			}
			if req.Nationality == "WNI" {
				if req.NIK != "" {
					nikRe := regexp.MustCompile(`^\d{16}$`)
					if !nikRe.MatchString(req.NIK) {
						minimalErrors = append(minimalErrors, "NIK harus 16 digit angka")
					}
				}
			} else if req.Nationality == "WNA" {
				// Jika WNA dan passport diisi, pastikan format minimal
				if req.PassportNo != "" && len(req.PassportNo) < 5 {
					minimalErrors = append(minimalErrors, "Nomor passport tidak valid")
				}
			}
			if len(minimalErrors) > 0 {
				return c.Status(400).JSON(fiber.Map{
					"error":   "validation failed",
					"details": minimalErrors,
				})
			}
		}

		// Cek NIK duplikat
		var exists int64
		db.Model(&models.Worker{}).Where("user_id = ? AND nik = ?", user.ID, req.NIK).Count(&exists)
		if exists > 0 {
			return c.Status(400).JSON(fiber.Map{"error": "NIK sudah terdaftar"})
		}

		// Kalkulasi iuran otomatis (default semua program aktif)
		program := ProgramBPJSRequest{
			JKK: true,
			JKM: true,
			JHT: true,
			JP:  true,
		}
		iuran := KalkulasiIuran(req.Upah, program, 0.54) // Default risiko rendah
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

		// Parse additional dates
		var tanggalAwalBekerja *time.Time
		if req.TanggalAwalBekerja != "" {
			if t, err := time.Parse("2006-01-02", req.TanggalAwalBekerja); err == nil {
				tanggalAwalBekerja = &t
			}
		}
		var tanggalAkhirKontrak *time.Time
		if req.TanggalAkhirKontrak != "" {
			if t, err := time.Parse("2006-01-02", req.TanggalAkhirKontrak); err == nil {
				tanggalAkhirKontrak = &t
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
			log.Printf("Error creating worker: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "failed to create worker"})
		}

		// Response dengan informasi iuran
		response := fiber.Map{
			"id":      w.ID, // expose id di top-level untuk kemudahan frontend
			"worker":  w,
			"iuran":   iuran,
			"message": "TK berhasil ditambahkan",
			"draft":   isDraft,
		}

		return c.Status(201).JSON(response)
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
		if body.Nationality == "" {
			body.Nationality = w.Nationality
		}
		if body.Nationality == "WNI" {
			nikRe := regexp.MustCompile(`^\d{16}$`)
			if body.NIK == "" || !nikRe.MatchString(body.NIK) {
				return c.Status(400).JSON(fiber.Map{"error": "Format NIK tidak valid"})
			}
		} else if body.Nationality == "WNA" {
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
		w.Nationality = body.Nationality
		w.PassportNo = body.PassportNo
		// Extra fields (only overwrite if provided)
		if body.TempatLahir != "" {
			w.TempatLahir = body.TempatLahir
		}
		if body.IbuKandung != "" {
			w.IbuKandung = body.IbuKandung
		}
		if body.JenisKelamin != "" {
			w.JenisKelamin = body.JenisKelamin
		}
		if body.GolDarah != "" {
			w.GolDarah = body.GolDarah
		}
		if body.StatusKawin != "" {
			w.StatusKawin = body.StatusKawin
		}
		if body.StatusPegawai != "" {
			w.StatusPegawai = body.StatusPegawai
		}
		if body.LokasiPekerjaan != "" {
			w.LokasiPekerjaan = body.LokasiPekerjaan
		}
		if body.Alamat != "" {
			w.Alamat = body.Alamat
		}
		if body.Kabupaten != "" {
			w.Kabupaten = body.Kabupaten
		}
		if body.KodePos != "" {
			w.KodePos = body.KodePos
		}
		if body.TeleponAreaRumah != "" {
			w.TeleponAreaRumah = body.TeleponAreaRumah
		}
		if body.TeleponRumah != "" {
			w.TeleponRumah = body.TeleponRumah
		}
		if body.TeleponAreaKantor != "" {
			w.TeleponAreaKantor = body.TeleponAreaKantor
		}
		if body.TeleponKantor != "" {
			w.TeleponKantor = body.TeleponKantor
		}
		if body.TeleponExtKantor != "" {
			w.TeleponExtKantor = body.TeleponExtKantor
		}
		if body.Handphone != "" {
			w.Handphone = body.Handphone
		}
		if body.NPWP != "" {
			w.NPWP = body.NPWP
		}
		if body.Email != "" {
			w.Email = body.Email
		}
		if body.JenisIdentitas != "" {
			w.JenisIdentitas = body.JenisIdentitas
		}
		if body.MasaLakuIdentitas != "" {
			w.MasaLakuIdentitas = body.MasaLakuIdentitas
		}
		if body.SuratMenyuratKe != "" {
			w.SuratMenyuratKe = body.SuratMenyuratKe
		}
		if body.TanggalKepesertaan != "" {
			w.TanggalKepesertaan = body.TanggalKepesertaan
		}
		if body.KodeNegara != "" {
			w.KodeNegara = body.KodeNegara
		}
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
		if body.TanggalAwalBekerja != "" {
			if t, perr := time.Parse("2006-01-02", body.TanggalAwalBekerja); perr == nil {
				w.TanggalAwalBekerja = &t
			}
		}
		if body.TanggalAkhirKontrak != "" {
			if t, perr := time.Parse("2006-01-02", body.TanggalAkhirKontrak); perr == nil {
				w.TanggalAkhirKontrak = &t
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

// KalkulasiIuranEndpoint untuk menghitung iuran BPJS
func KalkulasiIuranEndpoint(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := middleware.GetUserFromContext(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		var body struct {
			Upah        float64            `json:"upah" validate:"required,min=1000000"`
			Program     ProgramBPJSRequest `json:"program"`
			RisikoUsaha float64            `json:"risikoUsaha"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}

		// Default program jika tidak diisi
		if !body.Program.JKK && !body.Program.JKM && !body.Program.JHT && !body.Program.JP {
			body.Program = ProgramBPJSRequest{
				JKK: true,
				JKM: true,
				JHT: true,
				JP:  true,
			}
		}

		// Kalkulasi iuran
		iuran := KalkulasiIuran(body.Upah, body.Program, body.RisikoUsaha)

		return c.JSON(fiber.Map{
			"upah":        body.Upah,
			"program":     body.Program,
			"risikoUsaha": body.RisikoUsaha,
			"iuran":       iuran,
		})
	}
}

// GetProgramBPJS untuk mendapatkan daftar program BPJS
func GetProgramBPJS(db *database.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		programs := fiber.Map{
			"jkk": fiber.Map{
				"name":    "Jaminan Kecelakaan Kerja",
				"rate":    "0.24% - 1.74% (tergantung risiko usaha)",
				"company": "100%",
				"worker":  "0%",
			},
			"jkm": fiber.Map{
				"name":    "Jaminan Kematian",
				"rate":    "0.3%",
				"company": "100%",
				"worker":  "0%",
			},
			"jht": fiber.Map{
				"name":    "Jaminan Hari Tua",
				"rate":    "5.7%",
				"company": "3.7%",
				"worker":  "2%",
			},
			"jp": fiber.Map{
				"name":    "Jaminan Pensiun",
				"rate":    "3%",
				"company": "2%",
				"worker":  "1%",
			},
		}

		return c.JSON(fiber.Map{
			"programs":    programs,
			"description": "Program kepesertaan BPJS Ketenagakerjaan sesuai UU No. 24 Tahun 2011",
		})
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

			// Mapping sesuai template SIPP yang benar (33 kolom A-AG)
			req := CreateWorkerRequest{
				// A-C. Data Identitas Dasar
				NoPegawai: get(0), // A: NO_PEGAWAI
				Nama:      get(1), // B: NAMA_LENGKAP
				Gelar:     get(2), // C: GELAR

				// D-E. Kontak Rumah
				TeleponAreaRumah: get(3), // D: TELEPON_AREA_RUMAH
				TeleponRumah:     get(4), // E: TELEPON_RUMAH

				// G-I. Kontak Kantor (kolom F kosong)
				TeleponAreaKantor: get(6), // G: TELEPON_AREA_KANTOR
				TeleponKantor:     get(7), // H: TELEPON_KANTOR
				TeleponExtKantor:  get(8), // I: TELEPON_EXT_KANTOR

				// J-K. Kontak Utama
				Handphone: get(9),  // J: HP
				Email:     get(10), // K: EMAIL

				// L-M. Data Kelahiran
				TempatLahir: get(11), // L: TEMPAT_LAHIR
				DateOfBirth: get(12), // M: TANGGAL_LAHIR

				// N. Data Keluarga
				IbuKandung: get(13), // N: NAMA_IBU_KANDUNG

				// O-Q. Identitas
				JenisIdentitas:    get(14), // O: JENIS_IDENTITAS
				NIK:               get(15), // P: NO_IDENTITAS (NIK)
				MasaLakuIdentitas: get(16), // Q: MASA_LAKU_IDENTITAS

				// R. Jenis Kelamin
				JenisKelamin: get(17), // R: JENIS_KELAMIN

				// S-T. Kepesertaan
				SuratMenyuratKe:    get(18), // S: SURAT_MENYURAT_KE
				TanggalKepesertaan: get(19), // T: TANGGAL_KEPESERTAAN

				// U-V. Status Personal
				StatusKawin: get(20), // U: STATUS_KAWIN
				GolDarah:    get(21), // V: GOLONGAN_DARAH

				// W-X. Identitas Pajak & Negara
				NPWP:       get(22), // W: NPWP
				KodeNegara: get(23), // X: KODE_NEGARA

				// Y. Upah/Gaji
				Upah: parseUpah(get(24)), // Y: UPAH

				// Z-AA. Alamat
				Alamat:  get(25), // Z: ALAMAT
				KodePos: get(26), // AA: KODE_POS

				// AB-AC. Pekerjaan
				LokasiPekerjaan: get(27), // AB: LOKASI_PEKERJAAN
				StatusPegawai:   get(28), // AC: STATUS_PEGAWAI

				// AD-AE. Tanggal Kerja
				TanggalAwalBekerja:  get(29), // AD: TGL_AWAL_BEKERJA
				TanggalAkhirKontrak: get(30), // AE: TGL_AKHIR_KONTRAK

				// AF-AG. Data Tambahan
				Rapel:       parseUpah(get(31)), // AF: RAPEL
				Nationality: get(32),            // AG: NATIONALITY

				// Data tambahan untuk backward compatibility
				KPJ: "", // Tidak dipetakan dari template
			}

			if req.Nationality == "" {
				req.Nationality = "WNI"
			}

			// Validasi data sesuai SIPP
			validationErrors := ValidasiSIPP(&req)
			if len(validationErrors) > 0 {
				log.Printf("Validation failed for row %d: %v", i+1, validationErrors)
				invalid++
				continue
			}

			// Cek NIK duplikat
			var exists int64
			db.Model(&models.Worker{}).Where("user_id = ? AND nik = ?", user.ID, req.NIK).Count(&exists)
			if exists > 0 {
				log.Printf("NIK duplikat for row %d: %s", i+1, req.NIK)
				invalid++
				continue
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
			// Parse additional dates for mass upload
			var tanggalAwalBekerja *time.Time
			if req.TanggalAwalBekerja != "" {
				if t, err := time.Parse("2006-01-02", req.TanggalAwalBekerja); err == nil {
					tanggalAwalBekerja = &t
				}
			}
			var tanggalAkhirKontrak *time.Time
			if req.TanggalAkhirKontrak != "" {
				if t, err := time.Parse("2006-01-02", req.TanggalAkhirKontrak); err == nil {
					tanggalAkhirKontrak = &t
				}
			}

			// Kalkulasi iuran otomatis (default semua program aktif)
			program := ProgramBPJSRequest{
				JKK: true,
				JKM: true,
				JHT: true,
				JP:  true,
			}
			_ = KalkulasiIuran(req.Upah, program, 0.54) // Default risiko rendah

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

// parseUpah attempts to interpret various human inputs ("31,00", "5.000.000", "5,000,000.00")
// and returns rupiah value. Heuristic: if parsed value < 1000, we treat it as "in millions"
// commonly written as 31,00 (meaning 31 juta) and multiply by 1_000_000.
func parseUpah(s string) float64 {
	if s == "" {
		return 0
	}
	original := s
	s = strings.TrimSpace(s)
	// Normalize common formats
	// Case 1: European style 5.000.000,50 → remove dots, replace comma with dot
	if strings.Contains(s, ",") && strings.Count(s, ".") >= 1 {
		s = strings.ReplaceAll(s, ".", "")
		s = strings.ReplaceAll(s, ",", ".")
	} else if strings.Contains(s, ",") && !strings.Contains(s, ".") {
		// Case 2: "31,00" → 31.00
		s = strings.ReplaceAll(s, ",", ".")
	} else {
		// Case 3: remove thousand separators like 5,000,000.00
		s = strings.ReplaceAll(s, ",", "")
	}
	val := parseFloat(s)
	// Heuristic: if value extremely small (e.g., 31.00), interpret as millions
	if val > 0 && val < 1000 {
		// Only apply if original looked like compact format (had comma or no thousand separators)
		if strings.Contains(original, ",") || (!strings.Contains(original, ".") && !strings.Contains(original, ",")) {
			val = val * 1_000_000
		}
	}
	return val
}

// isLikelyKPJ returns true if a string looks like KPJ (mostly digits with length >= 6)
func isLikelyKPJ(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// remove spaces and dashes
	clean := strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "-", "")
	if len(clean) < 6 || len(clean) > 20 {
		return false
	}
	for _, ch := range clean {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
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

		// Build header index map to support wide Excel templates
		headerIdx := map[string]int{}
		if len(rows) > 0 {
			for i, h := range rows[0] {
				key := strings.TrimSpace(strings.ToUpper(h))
				headerIdx[key] = i
			}
		}

		// Helper to get value by header name (case-insensitive)
		getByHeader := func(key string, r []string) string {
			idx, ok := headerIdx[strings.TrimSpace(strings.ToUpper(key))]
			if !ok || idx >= len(r) {
				return ""
			}
			return strings.TrimSpace(r[idx])
		}

		// Skip header row
		dataRows := rows[1:]
		total := len(dataRows)
		valid := 0
		invalid := 0
		var rowErrors []string

		// Process each row
		for i, r := range dataRows {
			// Try wide header-based format first
			nik := getByHeader("NIK", r)
			nama := getByHeader("NAMA_LENGKAP", r)
			kpj := getByHeader("KPJ", r)
			kodeTK := getByHeader("KODE_TK", r)
			if kodeTK == "" { // alternate header names commonly used
				kodeTK = getByHeader("KODE", r)
			}
			upahStr := getByHeader("UPAH", r)
			rapelStr := getByHeader("RAPEL", r)

			// Fallback to legacy 4-column format if headers missing
			if nik == "" && len(r) >= 4 {
				nik = strings.TrimSpace(r[0])
				nama = strings.TrimSpace(r[1])
				upahStr = strings.TrimSpace(r[2])
				rapelStr = strings.TrimSpace(r[3])
			}

			// Normalization helpers
			digits := func(s string) string {
				var b strings.Builder
				for _, ch := range s {
					if ch >= '0' && ch <= '9' {
						b.WriteRune(ch)
					}
				}
				return b.String()
			}
			nik = digits(nik)
			kpj = digits(kpj)
			kodeTK = strings.TrimSpace(kodeTK)

			if nik == "" && kpj == "" && kodeTK == "" {
				invalid++
				rowErrors = append(rowErrors, fmt.Sprintf("Row %d: Tidak ada identitas (KPJ/NIK/KODE_TK) untuk pencarian", i+2))
				continue
			}
			if nama == "" {
				invalid++
				rowErrors = append(rowErrors, fmt.Sprintf("Row %d: Nama kosong", i+2))
				continue
			}

			// Parse upah dan rapel (terima berbagai format)
			upah := parseUpah(upahStr)
			rapel := parseUpah(rapelStr)

			// Find existing worker by priority: KPJ → NIK → NoPegawai (kode)
			var worker models.Worker
			var findErr error
			if kpj != "" {
				findErr = db.Where("user_id = ? AND kpj = ?", user.ID, kpj).First(&worker).Error
			}
			if findErr != nil && nik != "" {
				findErr = db.Where("user_id = ? AND nik = ?", user.ID, nik).First(&worker).Error
			}
			if findErr != nil && kodeTK != "" {
				findErr = db.Where("user_id = ? AND no_pegawai = ?", user.ID, kodeTK).First(&worker).Error
			}
			if findErr != nil {
				invalid++
				rowErrors = append(rowErrors, fmt.Sprintf("Row %d: Worker tidak ditemukan (KPJ=%s, NIK=%s, KODE_TK=%s)", i+2, kpj, nik, kodeTK))
				continue
			}

			// Update worker's upah and rapel
			worker.Upah = upah
			worker.Rapel = rapel

			if err := db.Save(&worker).Error; err != nil {
				invalid++
				rowErrors = append(rowErrors, fmt.Sprintf("Row %d: Gagal menyimpan upah untuk NIK %s - %v", i+2, nik, err))
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
			"errors":    rowErrors,
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

		// Process each row - Template lanjutan (update_data_tk)
		var rowErrors []string
		for i, r := range dataRows {
			get := func(idx int) string {
				if idx < len(r) {
					return strings.TrimSpace(r[idx])
				}
				return ""
			}

			// Kolom berdasarkan template yang Anda kirim
			kode := get(0)             // KODE (kami simpan ke NoPegawai)
			noPegawai := get(1)        // NOMOR_PEGAWAI (bila ada)
			nik := get(2)              // NOMOR_IDENTITAS / NIK
			kpj := get(3)              // KPJ
			nama := get(4)             // NAMA_TENAGA_KERJA
			tempatLahir := get(5)      // TEMPAT_LAHIR
			tglLahir := get(6)         // TGL_LAHIR (22-11-1990)
			jenisKelamin := get(7)     // JENIS_KELAMIN (L/P atau teks)
			ibuKandung := get(8)       // NAMA_IBU_KANDUNG
			npwp := get(9)             // NPWP
			alamat := get(10)          // ALAMAT_LENGKAP_DOMISILI
			handphone := get(11)       // HANDPHONE
			lokasiPekerjaan := get(12) // LOKASI_PEKERJAAN (kode)
			email := get(13)           // EMAIL
			// Bank fields (opsional)
			namaBank := get(14)
			_ = namaBank
			kodeBank := get(15)
			_ = kodeBank
			noRek := get(16)
			_ = noRek
			namaRek := get(17)
			_ = namaRek

			// Normalisasi KPJ dan NIK: hanya digit
			digits := func(s string) string {
				var b strings.Builder
				for _, ch := range s {
					if ch >= '0' && ch <= '9' {
						b.WriteRune(ch)
					}
				}
				return b.String()
			}
			kpj = digits(kpj)
			nik = digits(nik)

			// Normalisasi jenis kelamin: terima berbagai format dan map ke 'L' / 'P'
			if jk := strings.TrimSpace(strings.ToUpper(jenisKelamin)); jk != "" {
				switch {
				case jk == "L" || strings.HasPrefix(jk, "LAKI") || strings.HasPrefix(jk, "PRIA") || strings.HasPrefix(jk, "L-MALE"):
					jenisKelamin = "L"
				case jk == "P" || strings.HasPrefix(jk, "PEREM") || strings.HasPrefix(jk, "WANITA") || strings.HasPrefix(jk, "FEMALE"):
					jenisKelamin = "P"
				default:
					// Ambil huruf pertama bila hanya 1-2 char, kalau tidak, kosongkan agar tidak melanggar size:1
					if len(jk) > 0 {
						jenisKelamin = string(jk[0])
						if jenisKelamin != "L" && jenisKelamin != "P" {
							jenisKelamin = ""
						}
					}
				}
			}

			// Helper parse multi-format date
			parseDate := func(s string) (time.Time, bool) {
				layouts := []string{"2006-01-02", "02-01-2006", "02/01/2006", "2006/01/02"}
				for _, l := range layouts {
					if t, err := time.Parse(l, strings.TrimSpace(s)); err == nil {
						return t, true
					}
				}
				return time.Time{}, false
			}

			// Cari worker dengan prioritas: KPJ → NIK → NoPegawai(KODE atau NOMOR_PEGAWAI) → Nama+TglLahir
			var worker models.Worker
			var findErr error
			if kpj != "" {
				findErr = db.Where("user_id = ? AND kpj = ?", user.ID, kpj).First(&worker).Error
			}
			if findErr != nil && nik != "" {
				findErr = db.Where("user_id = ? AND nik = ?", user.ID, nik).First(&worker).Error
			}
			if findErr != nil {
				key := kode
				if key == "" {
					key = noPegawai
				}
				if key != "" {
					findErr = db.Where("user_id = ? AND no_pegawai = ?", user.ID, key).First(&worker).Error
				}
			}
			if findErr != nil && nama != "" {
				if dob, ok := parseDate(tglLahir); ok {
					findErr = db.Where("user_id = ? AND upper(nama) = ? AND date_of_birth = ?", user.ID, strings.ToUpper(nama), dob.Format("2006-01-02")).First(&worker).Error
				}
			}
			if findErr != nil {
				// UPSERT: buat baru jika tidak ditemukan
				newWorker := models.Worker{
					UserID: user.ID,
					KPJ:    kpj,
					NIK:    nik,
					NoPegawai: func() string {
						if noPegawai != "" {
							return noPegawai
						}
						return kode
					}(),
					Nama:            nama,
					TempatLahir:     tempatLahir,
					JenisKelamin:    jenisKelamin,
					IbuKandung:      ibuKandung,
					NPWP:            npwp,
					Alamat:          alamat,
					Handphone:       handphone,
					LokasiPekerjaan: lokasiPekerjaan,
					Email:           email,
					Nationality:     "WNI",
				}
				if dob, ok := parseDate(tglLahir); ok {
					newWorker.DateOfBirth = dob
				}
				if err := db.Create(&newWorker).Error; err != nil {
					invalid++
					rowErrors = append(rowErrors, fmt.Sprintf("Row %d: create failed - %v", i+2, err))
					continue
				}
				valid++
			} else {
				// Update
				if nama != "" {
					worker.Nama = nama
				}
				if kpj != "" {
					worker.KPJ = kpj
				}
				if kode != "" {
					worker.NoPegawai = kode
				}
				if noPegawai != "" {
					worker.NoPegawai = noPegawai
				}
				if tempatLahir != "" {
					worker.TempatLahir = tempatLahir
				}
				if jenisKelamin != "" {
					worker.JenisKelamin = jenisKelamin
				}
				if ibuKandung != "" {
					worker.IbuKandung = ibuKandung
				}
				if npwp != "" {
					worker.NPWP = npwp
				}
				if alamat != "" {
					worker.Alamat = alamat
				}
				if handphone != "" {
					worker.Handphone = handphone
				}
				if lokasiPekerjaan != "" {
					worker.LokasiPekerjaan = lokasiPekerjaan
				}
				if email != "" {
					worker.Email = email
				}
				if dob, ok := parseDate(tglLahir); ok {
					worker.DateOfBirth = dob
				}
				if err := db.Save(&worker).Error; err != nil {
					invalid++
					rowErrors = append(rowErrors, fmt.Sprintf("Row %d: update failed - %v", i+2, err))
					continue
				}
				valid++
			}
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
			"errors":    rowErrors,
		})
	}
}

// UploadTKNA menerima file .xlsx untuk menonaktifkan TK yang sudah ada
func UploadTKNA(db *database.DB) fiber.Handler {
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
		tempPath := "upload_tk_na_tmp_" + time.Now().Format("20060102150405") + "_" + fileHeader.Filename
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

		// Process each row - Format TK NA (SIPP): KPJ, NAMA_LENGKAP, TGL_LAHIR, SEBAB_NA, TGL_KEJADIAN, KETERANGAN
		var rowErrors []string
		for i, r := range dataRows {
			get := func(idx int) string {
				if idx < len(r) {
					return strings.TrimSpace(r[idx])
				}
				return ""
			}

			// Ambil data dari kolom template
			kpj := get(0)
			// Bersihkan KPJ: ambil hanya digit
			{
				var b strings.Builder
				for _, ch := range kpj {
					if ch >= '0' && ch <= '9' {
						b.WriteRune(ch)
					}
				}
				kpj = b.String()
			}
			nama := get(1)
			tglLahirStr := get(2)
			sebabNA := get(3)
			tglKejadianStr := get(4)
			keterangan := get(5)

			if kpj == "" && (nama == "" || tglLahirStr == "") {
				msg := fmt.Sprintf("Row %d: KPJ kosong dan data identifikasi (Nama/Tgl Lahir) tidak lengkap", i+2)
				log.Printf(msg)
				rowErrors = append(rowErrors, msg)
				invalid++
				continue
			}

			// Helper parse tanggal fleksibel
			parseDate := func(s string) (*time.Time, bool) {
				s = strings.TrimSpace(s)
				if s == "" {
					return nil, false
				}
				layouts := []string{"2006-01-02", "02-01-2006", "02/01/2006", "2006/01/02"}
				for _, l := range layouts {
					if t, err := time.Parse(l, s); err == nil {
						tt := t
						return &tt, true
					}
				}
				return nil, false
			}

			// Cari worker: prioritas KPJ, fallback Nama + Tgl Lahir
			var worker models.Worker
			var findErr error
			if kpj != "" {
				findErr = db.Where("user_id = ? AND kpj = ?", user.ID, kpj).First(&worker).Error
			}
			if findErr != nil {
				if nama != "" {
					if dob, ok := parseDate(tglLahirStr); ok {
						findErr = db.Where("user_id = ? AND upper(nama) = ? AND date_of_birth = ?", user.ID, strings.ToUpper(nama), dob.Format("2006-01-02")).First(&worker).Error
					}
				}
			}
			if findErr != nil {
				msg := fmt.Sprintf("Row %d: Worker tidak ditemukan (KPJ=%s, Nama=%s)", i+2, kpj, nama)
				log.Printf(msg)
				rowErrors = append(rowErrors, msg)
				invalid++
				continue
			}

			// Validasi nama: jika tidak sama, beri peringatan tapi lanjutkan (kasus ejaan/kapitalisasi)
			if nama != "" && strings.TrimSpace(strings.ToUpper(worker.Nama)) != strings.TrimSpace(strings.ToUpper(nama)) {
				log.Printf("Row %d: Nama berbeda. DB: %s, File: %s - tetap dinonaktifkan berdasarkan identifikasi lain", i+2, worker.Nama, nama)
			}

			// PILIHAN 1: Soft Delete (Recommended) - Tambah field status
			// Update status menjadi nonaktif
			worker.StatusPegawai = "NONAKTIF"
			// Set tanggal akhir kontrak = TGL_KEJADIAN bila ada, kalau tidak pakai hari ini
			if tkj, ok := parseDate(tglKejadianStr); ok {
				worker.TanggalAkhirKontrak = tkj
			} else {
				today := time.Now()
				t := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
				worker.TanggalAkhirKontrak = &t
			}
			if sebabNA != "" || keterangan != "" {
				// Bisa tambah field alasan_nonaktif jika diperlukan
				log.Printf("Worker %s dinonaktifkan. Sebab: %s. Ket: %s", worker.Nama, sebabNA, keterangan)
			}

			if err := db.Save(&worker).Error; err != nil {
				log.Printf("Row %d: Gagal menonaktifkan worker %s: %v", i+2, nama, err)
				invalid++
				continue
			}

			// PILIHAN 2: Hard Delete (Alternative) - Hapus dari database
			// Uncomment jika ingin menghapus data sepenuhnya
			// if err := db.Delete(&worker).Error; err != nil {
			//     log.Printf("Row %d: Gagal menghapus worker %s: %v", i+2, nama, err)
			//     invalid++
			//     continue
			// }

			log.Printf("Worker %s (KPJ: %s) berhasil dinonaktifkan", worker.Nama, worker.KPJ)
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
			Type:             "tk_na", // Type khusus untuk TK NA
		}
		if err := db.Create(&hist).Error; err != nil {
			log.Printf("Failed to save upload history: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error":   "failed to save history",
				"details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message":   "Upload TK Nonaktif berhasil",
			"totalData": total,
			"valid":     valid,
			"invalid":   invalid,
			"historyId": hist.ID,
			"action":    "nonaktif", // Indikator bahwa ini adalah proses nonaktif
			"errors":    rowErrors,
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
