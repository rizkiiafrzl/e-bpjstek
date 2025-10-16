# PowerShell script untuk cek error detail upload
Write-Host "🔍 Checking Upload Errors..." -ForegroundColor Red

# Set your API details
$API_URL = "http://localhost:8080/api/v1"
$TOKEN = "YOUR_TOKEN_HERE"  # Ganti dengan token yang valid

$headers = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

Write-Host "`n🔍 Step 1: Cek Upload History Terbaru" -ForegroundColor Yellow
try {
    $history = Invoke-RestMethod -Uri "$API_URL/workers/upload-history" -Headers $headers -Method GET
    if ($history.Count -gt 0) {
        $latest = $history[0]
        Write-Host "📊 Upload terbaru:" -ForegroundColor Cyan
        Write-Host "   File: $($latest.fileName)" -ForegroundColor White
        Write-Host "   Type: $($latest.type)" -ForegroundColor White
        Write-Host "   Total Data: $($latest.totalData)" -ForegroundColor White
        Write-Host "   Valid: $($latest.totalValid)" -ForegroundColor Green
        Write-Host "   Invalid: $($latest.totalInvalid)" -ForegroundColor Red
        Write-Host "   Status: $($latest.validationStatus)" -ForegroundColor White
        Write-Host "   Date: $($latest.created_at)" -ForegroundColor White
    } else {
        Write-Host "❌ Tidak ada upload history" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ Error getting upload history: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🔍 Step 2: Cek Data Workers Saat Ini" -ForegroundColor Yellow
try {
    $workers = Invoke-RestMethod -Uri "$API_URL/workers" -Headers $headers -Method GET
    Write-Host "📊 Total workers: $($workers.Count)" -ForegroundColor Cyan
    
    if ($workers.Count -eq 0) {
        Write-Host "❌ TIDAK ADA DATA WORKERS!" -ForegroundColor Red
        Write-Host "   Semua data invalid, tidak ada yang tersimpan" -ForegroundColor Yellow
    } else {
        Write-Host "✅ Ada $($workers.Count) data workers" -ForegroundColor Green
    }
} catch {
    Write-Host "❌ Error getting workers: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🔍 Step 3: Kemungkinan Penyebab Semua Data Invalid" -ForegroundColor Yellow
Write-Host "📋 Checklist masalah umum:" -ForegroundColor Cyan

Write-Host "1. ❓ Format Excel tidak sesuai template" -ForegroundColor Gray
Write-Host "   - Pastikan ada 31 kolom (A-AE)" -ForegroundColor Gray
Write-Host "   - Header row harus sesuai template" -ForegroundColor Gray
Write-Host "   - Data dimulai dari row 2" -ForegroundColor Gray

Write-Host "2. ❓ Data kosong di kolom wajib" -ForegroundColor Gray
Write-Host "   - Nama (kolom B) tidak boleh kosong" -ForegroundColor Gray
Write-Host "   - NIK (kolom P) harus 16 digit" -ForegroundColor Gray
Write-Host "   - Email (kolom K) format harus benar" -ForegroundColor Gray

Write-Host "3. ❓ Format data salah" -ForegroundColor Gray
Write-Host "   - NIK harus angka 16 digit" -ForegroundColor Gray
Write-Host "   - NPWP harus angka 15 digit" -ForegroundColor Gray
Write-Host "   - Golongan darah: A, B, AB, O saja" -ForegroundColor Gray
Write-Host "   - Jenis kelamin: Laki-laki/Perempuan" -ForegroundColor Gray

Write-Host "4. ❓ NIK duplikat" -ForegroundColor Gray
Write-Host "   - NIK sudah ada di database" -ForegroundColor Gray
Write-Host "   - Cek data yang sudah ada" -ForegroundColor Gray

Write-Host "`n🔍 Step 4: Debug Log Backend" -ForegroundColor Yellow
Write-Host "📋 Cek log backend untuk error detail:" -ForegroundColor Cyan
Write-Host "   - === UPLOAD TK DEBUG ===" -ForegroundColor Gray
Write-Host "   - Header row content" -ForegroundColor Gray
Write-Host "   - First data row content" -ForegroundColor Gray
Write-Host "   - Processing each row" -ForegroundColor Gray
Write-Host "   - Validation errors per row" -ForegroundColor Gray
Write-Host "   - Row X invalid: [error reason]" -ForegroundColor Gray

Write-Host "`n🛠️ Solusi Cepat:" -ForegroundColor Green
Write-Host "1. Cek log backend untuk error detail" -ForegroundColor White
Write-Host "2. Pastikan format Excel sesuai template 31 kolom" -ForegroundColor White
Write-Host "3. Isi data minimal yang valid:" -ForegroundColor White
Write-Host "   - Nama (kolom B)" -ForegroundColor Gray
Write-Host "   - NIK 16 digit (kolom P)" -ForegroundColor Gray
Write-Host "   - Email valid (kolom K)" -ForegroundColor Gray
Write-Host "4. Test dengan 1-2 data dulu" -ForegroundColor White
Write-Host "5. Cek apakah NIK sudah ada" -ForegroundColor White

Write-Host "`n📋 Template Excel yang Benar:" -ForegroundColor Yellow
Write-Host "A: NO_PEGAWAI" -ForegroundColor Gray
Write-Host "B: NAMA_LENGKAP (WAJIB)" -ForegroundColor Red
Write-Host "C: GELAR" -ForegroundColor Gray
Write-Host "D: TELEPON_AREA_RUMAH" -ForegroundColor Gray
Write-Host "E: TELEPON_RUMAH" -ForegroundColor Gray
Write-Host "F: (kosong)" -ForegroundColor Gray
Write-Host "G: TELEPON_AREA_KANTOR" -ForegroundColor Gray
Write-Host "H: TELEPON_KANTOR" -ForegroundColor Gray
Write-Host "I: TELEPON_EXT_KANTOR" -ForegroundColor Gray
Write-Host "J: HP" -ForegroundColor Gray
Write-Host "K: EMAIL (WAJIB VALID)" -ForegroundColor Red
Write-Host "L: TEMPAT_LAHIR" -ForegroundColor Gray
Write-Host "M: TANGGAL_LAHIR" -ForegroundColor Gray
Write-Host "N: NAMA_IBU_KANDUNG" -ForegroundColor Gray
Write-Host "O: JENIS_IDENTITAS" -ForegroundColor Gray
Write-Host "P: NO_IDENTITAS/NIK (WAJIB 16 DIGIT)" -ForegroundColor Red
Write-Host "Q: MASA_LAKU_IDENTITAS" -ForegroundColor Gray
Write-Host "R: JENIS_KELAMIN" -ForegroundColor Gray
Write-Host "S: SURAT_MENYURAT_KE" -ForegroundColor Gray
Write-Host "T: TANGGAL_KEPESERTAAN" -ForegroundColor Gray
Write-Host "U: STATUS_KAWIN" -ForegroundColor Gray
Write-Host "V: GOLONGAN_DARAH" -ForegroundColor Gray
Write-Host "W: NPWP" -ForegroundColor Gray
Write-Host "X: KODE_NEGARA" -ForegroundColor Gray
Write-Host "Y: UPAH" -ForegroundColor Gray
Write-Host "Z: ALAMAT" -ForegroundColor Gray
Write-Host "AA: KODE_POS" -ForegroundColor Gray
Write-Host "AB: LOKASI_PEKERJAAN" -ForegroundColor Gray
Write-Host "AC: STATUS_PEGAWAI" -ForegroundColor Gray
Write-Host "AD: TGL_AWAL_BEKERJA" -ForegroundColor Gray
Write-Host "AE: TGL_AKHIR_KONTRAK" -ForegroundColor Gray

Write-Host "`n✅ Check completed!" -ForegroundColor Green
Write-Host "Cek log backend untuk error detail yang spesifik" -ForegroundColor Cyan







