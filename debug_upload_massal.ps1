# PowerShell script untuk debug upload massal tidak muncul
Write-Host "🐛 Debugging Upload Massal Issue..." -ForegroundColor Red

# Set your API details
$API_URL = "http://localhost:8080/api/v1"
$TOKEN = "YOUR_TOKEN_HERE"  # Ganti dengan token yang valid

$headers = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

Write-Host "`n🔍 Step 1: Cek Data Workers Saat Ini" -ForegroundColor Yellow
try {
    $workers = Invoke-RestMethod -Uri "$API_URL/workers" -Headers $headers -Method GET
    Write-Host "📊 Total workers: $($workers.Count)" -ForegroundColor Cyan
    
    if ($workers.Count -eq 0) {
        Write-Host "❌ TIDAK ADA DATA WORKERS!" -ForegroundColor Red
    } else {
        Write-Host "✅ Ada $($workers.Count) data workers" -ForegroundColor Green
        Write-Host "`n📋 Data terbaru:" -ForegroundColor Cyan
        for ($i = 0; $i -lt [Math]::Min(3, $workers.Count); $i++) {
            $worker = $workers[$i]
            Write-Host "  $($i+1). $($worker.nama) ($($worker.noPegawai)) - $($worker.created_at)" -ForegroundColor White
        }
    }
} catch {
    Write-Host "❌ Error getting workers: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🔍 Step 2: Cek Upload History" -ForegroundColor Yellow
try {
    $history = Invoke-RestMethod -Uri "$API_URL/workers/upload-history" -Headers $headers -Method GET
    Write-Host "📊 Total uploads: $($history.Count)" -ForegroundColor Cyan
    
    if ($history.Count -eq 0) {
        Write-Host "❌ TIDAK ADA UPLOAD HISTORY!" -ForegroundColor Red
        Write-Host "   Kemungkinan upload gagal atau belum pernah upload" -ForegroundColor Yellow
    } else {
        Write-Host "✅ Ada $($history.Count) upload history" -ForegroundColor Green
        Write-Host "`n📋 Upload terbaru:" -ForegroundColor Cyan
        for ($i = 0; $i -lt [Math]::Min(3, $history.Count); $i++) {
            $h = $history[$i]
            Write-Host "  $($i+1). $($h.fileName)" -ForegroundColor White
            Write-Host "     Type: $($h.type), Valid: $($h.totalValid), Invalid: $($h.totalInvalid)" -ForegroundColor Gray
            Write-Host "     Status: $($h.validationStatus), Date: $($h.created_at)" -ForegroundColor Gray
        }
    }
} catch {
    Write-Host "❌ Error getting upload history: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🔍 Step 3: Cek Database Schema" -ForegroundColor Yellow
Write-Host "📋 Field yang diperlukan untuk SIPP:" -ForegroundColor Cyan
$requiredFields = @(
    "telepon_area_rumah", "telepon_area_kantor", "telepon_kantor", 
    "telepon_ext_kantor", "jenis_identitas", "masa_laku_identitas",
    "surat_menyurat_ke", "tanggal_kepesertaan", "kode_negara"
)

Write-Host "   Field ini harus ada di database untuk upload massal berfungsi" -ForegroundColor Gray
Write-Host "   Jalankan migration script jika field belum ada" -ForegroundColor Gray

Write-Host "`n🔍 Step 4: Test Upload Kecil" -ForegroundColor Yellow
Write-Host "📋 Rekomendasi untuk test:" -ForegroundColor Cyan
Write-Host "1. Buat Excel dengan 2-3 data saja" -ForegroundColor Gray
Write-Host "2. Pastikan format sesuai template" -ForegroundColor Gray
Write-Host "3. Upload dan lihat response" -ForegroundColor Gray
Write-Host "4. Cek log backend saat upload" -ForegroundColor Gray

Write-Host "`n🔍 Step 5: Cek Kemungkinan Masalah" -ForegroundColor Yellow
Write-Host "📋 Checklist masalah umum:" -ForegroundColor Cyan

# Cek apakah ada data tapi tidak tampil
if ($workers -and $workers.Count -gt 0) {
    Write-Host "✅ Data ada di database" -ForegroundColor Green
    Write-Host "❓ Kemungkinan masalah frontend:" -ForegroundColor Yellow
    Write-Host "   - Filter aktif (cari di UI)" -ForegroundColor Gray
    Write-Host "   - Pagination issue" -ForegroundColor Gray
    Write-Host "   - User ID mismatch" -ForegroundColor Gray
    Write-Host "   - Search query salah" -ForegroundColor Gray
} else {
    Write-Host "❌ Tidak ada data di database" -ForegroundColor Red
    Write-Host "❓ Kemungkinan masalah:" -ForegroundColor Yellow
    Write-Host "   - Upload gagal (cek upload history)" -ForegroundColor Gray
    Write-Host "   - Validasi data gagal" -ForegroundColor Gray
    Write-Host "   - Database schema tidak lengkap" -ForegroundColor Gray
    Write-Host "   - Backend error saat insert" -ForegroundColor Gray
}

Write-Host "`n🛠️ SOLUSI CEPAT:" -ForegroundColor Green
Write-Host "1. Restart backend: go run main.go" -ForegroundColor White
Write-Host "2. Jalankan migration: .\run_migration.ps1" -ForegroundColor White
Write-Host "3. Test upload dengan data minimal" -ForegroundColor White
Write-Host "4. Cek log backend saat upload" -ForegroundColor White
Write-Host "5. Refresh halaman frontend" -ForegroundColor White

Write-Host "`n📞 Jika masih bermasalah:" -ForegroundColor Yellow
Write-Host "- Cek log backend untuk error detail" -ForegroundColor Gray
Write-Host "- Verifikasi data di database langsung" -ForegroundColor Gray
Write-Host "- Test dengan data yang lebih sederhana" -ForegroundColor Gray







