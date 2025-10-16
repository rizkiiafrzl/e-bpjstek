# PowerShell script untuk test upload setelah perbaikan
Write-Host "🧪 Testing Upload Fix..." -ForegroundColor Green

# Set your API details
$API_URL = "http://localhost:8080/api/v1"
$TOKEN = "YOUR_TOKEN_HERE"  # Ganti dengan token yang valid

$headers = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

Write-Host "`n🔍 Step 1: Cek Data Sebelum Upload" -ForegroundColor Yellow
try {
    $workersBefore = Invoke-RestMethod -Uri "$API_URL/workers" -Headers $headers -Method GET
    Write-Host "📊 Workers sebelum upload: $($workersBefore.Count)" -ForegroundColor Cyan
} catch {
    Write-Host "❌ Error getting workers: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🔍 Step 2: Cek Upload History Sebelum" -ForegroundColor Yellow
try {
    $historyBefore = Invoke-RestMethod -Uri "$API_URL/workers/upload-history" -Headers $headers -Method GET
    Write-Host "📊 Upload history sebelum: $($historyBefore.Count)" -ForegroundColor Cyan
} catch {
    Write-Host "❌ Error getting upload history: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🔍 Step 3: Test Upload Excel" -ForegroundColor Yellow
Write-Host "📋 Untuk test upload:" -ForegroundColor Cyan
Write-Host "1. Pastikan backend sudah running dengan patch baru" -ForegroundColor Gray
Write-Host "2. Siapkan file Excel dengan format yang benar" -ForegroundColor Gray
Write-Host "3. Upload melalui frontend atau API" -ForegroundColor Gray
Write-Host "4. Cek log backend untuk debug info" -ForegroundColor Gray

Write-Host "`n🔍 Step 4: Cek Data Setelah Upload" -ForegroundColor Yellow
Write-Host "📋 Setelah upload, jalankan script ini lagi untuk cek:" -ForegroundColor Cyan
Write-Host "   - Jumlah workers bertambah" -ForegroundColor Gray
Write-Host "   - Upload history bertambah" -ForegroundColor Gray
Write-Host "   - Data muncul di tabel frontend" -ForegroundColor Gray

Write-Host "`n🔍 Step 5: Debug Log Backend" -ForegroundColor Yellow
Write-Host "📋 Cek log backend untuk debug info:" -ForegroundColor Cyan
Write-Host "   - === UPLOAD TK DEBUG ===" -ForegroundColor Gray
Write-Host "   - Total rows in Excel" -ForegroundColor Gray
Write-Host "   - Header row content" -ForegroundColor Gray
Write-Host "   - Processing each row" -ForegroundColor Gray
Write-Host "   - Validation errors" -ForegroundColor Gray
Write-Host "   - Success/failure messages" -ForegroundColor Gray

Write-Host "`n🛠️ Troubleshooting:" -ForegroundColor Yellow
Write-Host "Jika upload masih gagal:" -ForegroundColor Cyan
Write-Host "1. Cek format Excel sesuai template 31 kolom" -ForegroundColor Gray
Write-Host "2. Pastikan data valid (NIK 16 digit, email format benar)" -ForegroundColor Gray
Write-Host "3. Cek log backend untuk error detail" -ForegroundColor Gray
Write-Host "4. Test dengan data minimal (2-3 rows)" -ForegroundColor Gray
Write-Host "5. Pastikan database migration sudah dijalankan" -ForegroundColor Gray

Write-Host "`n📋 Template Excel yang Benar:" -ForegroundColor Yellow
Write-Host "Kolom A: NO_PEGAWAI" -ForegroundColor Gray
Write-Host "Kolom B: NAMA_LENGKAP" -ForegroundColor Gray
Write-Host "Kolom D: TELEPON_AREA_RUMAH" -ForegroundColor Gray
Write-Host "Kolom E: TELEPON_RUMAH" -ForegroundColor Gray
Write-Host "Kolom K: EMAIL" -ForegroundColor Gray
Write-Host "Kolom P: NO_IDENTITAS (NIK 16 digit)" -ForegroundColor Gray
Write-Host "Kolom Y: UPAH" -ForegroundColor Gray
Write-Host "... dan seterusnya sesuai template SIPP" -ForegroundColor Gray

Write-Host "`n✅ Test script completed!" -ForegroundColor Green
Write-Host "Jalankan script ini lagi setelah upload untuk verifikasi hasil" -ForegroundColor Cyan







