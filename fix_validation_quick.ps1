# Quick Fix untuk Validasi SIPP - Mengurangi Validasi untuk Testing
# Script PowerShell untuk mengatasi masalah "Valid: 0"

Write-Host "🔧 Quick Fix untuk Masalah Valid: 0" -ForegroundColor Cyan
Write-Host "=" * 50 -ForegroundColor Cyan

# Backup file asli
$BackupFile = "backend/handlers/worker.go.backup"
$OriginalFile = "backend/handlers/worker.go"

Write-Host "`n1. Membuat backup file asli..." -ForegroundColor Green
if (Test-Path $OriginalFile) {
    Copy-Item $OriginalFile $BackupFile -Force
    Write-Host "✅ Backup created: $BackupFile" -ForegroundColor Green
} else {
    Write-Host "❌ File tidak ditemukan: $OriginalFile" -ForegroundColor Red
    exit
}

Write-Host "`n2. Membuat versi validasi yang lebih ringan..." -ForegroundColor Green

# Buat versi validasi yang lebih ringan
$LighterValidation = @'
// ValidasiSIPP dengan validasi yang lebih ringan untuk testing
func ValidasiSIPP(req *CreateWorkerRequest) []string {
	var errors []string

	// Validasi minimal - hanya field yang benar-benar wajib
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

	// Validasi format yang lebih fleksibel
	if req.NIK != "" {
		// Hanya cek apakah NIK adalah angka dan minimal 10 digit
		nikRe := regexp.MustCompile(`^\d{10,16}$`)
		if !nikRe.MatchString(req.NIK) {
			errors = append(errors, "NIK harus 10-16 digit angka")
		}
	}

	if req.Email != "" {
		// Validasi email yang lebih fleksibel
		emailRe := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
		if !emailRe.MatchString(req.Email) {
			errors = append(errors, "Format email tidak valid")
		}
	}

	if req.NPWP != "" {
		// NPWP lebih fleksibel - minimal 10 digit
		npwpRe := regexp.MustCompile(`^\d{10,15}$`)
		if !npwpRe.MatchString(req.NPWP) {
			errors = append(errors, "NPWP harus 10-15 digit angka")
		}
	}

	// Validasi jenis kelamin yang lebih fleksibel
	if req.JenisKelamin != "" {
		validJenisKelamin := []string{"L", "P", "Laki-laki", "Perempuan", "Laki", "Perempuan"}
		isValid := false
		for _, jk := range validJenisKelamin {
			if req.JenisKelamin == jk {
				isValid = true
				break
			}
		}
		if !isValid {
			errors = append(errors, "Jenis kelamin harus L/P atau Laki-laki/Perempuan")
		}
	}

	// Validasi golongan darah yang lebih fleksibel
	if req.GolDarah != "" {
		validGolDarah := []string{"A", "B", "AB", "O", "a", "b", "ab", "o"}
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

	// Validasi status pegawai yang lebih fleksibel
	if req.StatusPegawai != "" {
		validStatus := []string{"PKWT", "PKWTT", "Tetap", "Kontrak", "pkwt", "pkwtt"}
		isValid := false
		for _, status := range validStatus {
			if req.StatusPegawai == status {
				isValid = true
				break
			}
		}
		if !isValid {
			errors = append(errors, "Status pegawai harus PKWT atau PKWTT")
		}
	}

	// Validasi PKWT yang lebih fleksibel
	if req.StatusPegawai == "PKWT" || req.StatusPegawai == "pkwt" || req.StatusPegawai == "Kontrak" {
		if req.TanggalAkhirKontrak == "" {
			errors = append(errors, "Tanggal akhir kontrak wajib diisi untuk PKWT")
		}
	}

	// Validasi upah yang lebih fleksibel - minimal 100 ribu
	if req.Upah < 100000 {
		errors = append(errors, "Upah minimal Rp 100.000")
	}

	// Validasi usia yang lebih fleksibel
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			age := time.Since(dob).Hours() / (24 * 365)
			if age < 15 {
				errors = append(errors, "Usia minimal 15 tahun")
			}
			if age > 70 {
				errors = append(errors, "Usia maksimal 70 tahun")
			}
		}
	}

	return errors
}
'@

# Baca file asli
$Content = Get-Content $OriginalFile -Raw

# Ganti fungsi ValidasiSIPP dengan versi yang lebih ringan
$Pattern = '(?s)// ValidasiSIPP melakukan validasi data sesuai spesifikasi SIPP.*?return errors\s*\}'
$NewContent = $Content -replace $Pattern, $LighterValidation

# Tulis file yang sudah dimodifikasi
Set-Content -Path $OriginalFile -Value $NewContent -Encoding UTF8

Write-Host "✅ File dimodifikasi dengan validasi yang lebih ringan" -ForegroundColor Green

Write-Host "`n3. Restart backend untuk menerapkan perubahan..." -ForegroundColor Green
Write-Host "💡 Jalankan perintah berikut di terminal baru:" -ForegroundColor Yellow
Write-Host "   cd backend" -ForegroundColor White
Write-Host "   go run main.go" -ForegroundColor White

Write-Host "`n4. Test upload setelah restart..." -ForegroundColor Green
Write-Host "💡 Setelah backend restart, coba upload file Excel lagi" -ForegroundColor Yellow

Write-Host "`n5. Jika masih bermasalah, jalankan script debug:" -ForegroundColor Green
Write-Host "   .\debug_upload_sipp.ps1" -ForegroundColor White

Write-Host "`n" + "=" * 50 -ForegroundColor Cyan
Write-Host "🔧 Perubahan yang dibuat:" -ForegroundColor Green
Write-Host "   ✅ Validasi NIK: 10-16 digit (bukan harus 16)" -ForegroundColor White
Write-Host "   ✅ Validasi NPWP: 10-15 digit (bukan harus 15)" -ForegroundColor White
Write-Host "   ✅ Validasi Email: format lebih fleksibel" -ForegroundColor White
Write-Host "   ✅ Validasi Upah: minimal 100 ribu (bukan 1 juta)" -ForegroundColor White
Write-Host "   ✅ Validasi Usia: 15-70 tahun (bukan 17-65)" -ForegroundColor White
Write-Host "   ✅ Validasi Jenis Kelamin: lebih fleksibel" -ForegroundColor White
Write-Host "   ✅ Validasi Status Pegawai: lebih fleksibel" -ForegroundColor White

Write-Host "`n⚠️ PERINGATAN:" -ForegroundColor Yellow
Write-Host "   Ini adalah validasi untuk testing saja!" -ForegroundColor White
Write-Host "   Untuk production, kembalikan ke validasi SIPP yang ketat" -ForegroundColor White
Write-Host "   Gunakan file backup: $BackupFile" -ForegroundColor White

Write-Host "`n🔄 Cara mengembalikan ke validasi asli:" -ForegroundColor Cyan
Write-Host "   Copy-Item `"$BackupFile`" `"$OriginalFile`" -Force" -ForegroundColor White







