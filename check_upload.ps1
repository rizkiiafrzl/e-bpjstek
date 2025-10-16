# PowerShell script untuk mengecek data yang sudah diupload massal
# Pastikan data Excel sudah tersimpan di database dengan benar

Write-Host "🔍 Checking uploaded data from Excel..." -ForegroundColor Green

# Set your API details
$API_URL = "http://localhost:8080/api/v1"
$TOKEN = "YOUR_TOKEN_HERE"  # Ganti dengan token yang valid

# Headers
$headers = @{
    "Authorization" = "Bearer $TOKEN"
    "Content-Type" = "application/json"
}

try {
    # 1. Cek data workers terbaru
    Write-Host "`n=== DATA WORKERS TERBARU ===" -ForegroundColor Yellow
    $workersResponse = Invoke-RestMethod -Uri "$API_URL/workers" -Headers $headers -Method GET
    
    Write-Host "📊 Total workers: $($workersResponse.Count)" -ForegroundColor Cyan
    
    Write-Host "`n📋 Data terbaru (5 record):" -ForegroundColor Cyan
    for ($i = 0; $i -lt [Math]::Min(5, $workersResponse.Count); $i++) {
        $worker = $workersResponse[$i]
        Write-Host "`n$($i+1). $($worker.nama) ($($worker.noPegawai))" -ForegroundColor White
        Write-Host "   NIK: $($worker.nik)" -ForegroundColor Gray
        Write-Host "   Email: $($worker.email)" -ForegroundColor Gray
        Write-Host "   Telepon Area Rumah: $($worker.teleponAreaRumah)" -ForegroundColor Gray
        Write-Host "   Telepon Kantor: $($worker.teleponKantor)" -ForegroundColor Gray
        Write-Host "   Jenis Identitas: $($worker.jenisIdentitas)" -ForegroundColor Gray
        Write-Host "   Tanggal Kepesertaan: $($worker.tanggalKepesertaan)" -ForegroundColor Gray
        Write-Host "   Kode Negara: $($worker.kodeNegara)" -ForegroundColor Gray
        Write-Host "   Upah: $($worker.upah)" -ForegroundColor Gray
        Write-Host "   Created: $($worker.created_at)" -ForegroundColor Gray
    }

    # 2. Cek upload history
    Write-Host "`n=== UPLOAD HISTORY ===" -ForegroundColor Yellow
    $historyResponse = Invoke-RestMethod -Uri "$API_URL/workers/upload-history" -Headers $headers -Method GET
    
    Write-Host "📊 Total uploads: $($historyResponse.Count)" -ForegroundColor Cyan
    
    Write-Host "`n📋 Upload history terbaru:" -ForegroundColor Cyan
    for ($i = 0; $i -lt [Math]::Min(5, $historyResponse.Count); $i++) {
        $history = $historyResponse[$i]
        Write-Host "`n$($i+1). $($history.fileName)" -ForegroundColor White
        Write-Host "   Type: $($history.type)" -ForegroundColor Gray
        Write-Host "   Total Data: $($history.totalData)" -ForegroundColor Gray
        Write-Host "   Valid: $($history.totalValid)" -ForegroundColor Gray
        Write-Host "   Invalid: $($history.totalInvalid)" -ForegroundColor Gray
        Write-Host "   Status: $($history.validationStatus)" -ForegroundColor Gray
        Write-Host "   Uploaded: $($history.created_at)" -ForegroundColor Gray
    }

    # 3. Cek field mapping
    Write-Host "`n=== FIELD MAPPING CHECK ===" -ForegroundColor Yellow
    Write-Host "🔍 Checking field mapping from Excel to Database:" -ForegroundColor Cyan
    
    $fieldChecks = @{
        "TeleponAreaRumah" = 0
        "TeleponAreaKantor" = 0
        "TeleponKantor" = 0
        "TeleponExtKantor" = 0
        "JenisIdentitas" = 0
        "TanggalKepesertaan" = 0
        "KodeNegara" = 0
    }
    
    foreach ($worker in $workersResponse) {
        if ($worker.teleponAreaRumah) { $fieldChecks["TeleponAreaRumah"]++ }
        if ($worker.teleponAreaKantor) { $fieldChecks["TeleponAreaKantor"]++ }
        if ($worker.teleponKantor) { $fieldChecks["TeleponKantor"]++ }
        if ($worker.teleponExtKantor) { $fieldChecks["TeleponExtKantor"]++ }
        if ($worker.jenisIdentitas) { $fieldChecks["JenisIdentitas"]++ }
        if ($worker.tanggalKepesertaan) { $fieldChecks["TanggalKepesertaan"]++ }
        if ($worker.kodeNegara) { $fieldChecks["KodeNegara"]++ }
    }
    
    foreach ($field in $fieldChecks.Keys) {
        $count = $fieldChecks[$field]
        if ($count -gt 0) {
            Write-Host "✅ $field`: $count records" -ForegroundColor Green
        } else {
            Write-Host "❌ $field`: No data found" -ForegroundColor Red
        }
    }

    Write-Host "`n✅ Data check completed!" -ForegroundColor Green

} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "`n📋 Troubleshooting:" -ForegroundColor Yellow
    Write-Host "1. Pastikan backend sudah running" -ForegroundColor Gray
    Write-Host "2. Pastikan token valid" -ForegroundColor Gray
    Write-Host "3. Cek koneksi ke API" -ForegroundColor Gray
}

Write-Host "`n📋 Jika data kosong atau tidak sesuai:" -ForegroundColor Yellow
Write-Host "1. Pastikan backend sudah di-restart" -ForegroundColor Gray
Write-Host "2. Pastikan database migration sudah dijalankan" -ForegroundColor Gray
Write-Host "3. Cek log backend untuk error" -ForegroundColor Gray
Write-Host "4. Test upload Excel lagi" -ForegroundColor Gray







