# Debug Detail Error Upload - Script untuk melihat error spesifik
# Script PowerShell untuk mengidentifikasi masalah validasi

Write-Host "🔍 Debug Detail Error Upload" -ForegroundColor Cyan
Write-Host "=" * 50 -ForegroundColor Cyan

# Configuration
$ApiUrl = "http://localhost:8080/api/v1/workers"
$UploadUrl = "http://localhost:8080/api/v1/workers/upload-tk"
$AuthToken = "YOUR_JWT_TOKEN_HERE" # !!! IMPORTANT: Replace with valid JWT token !!!

if ($AuthToken -eq "YOUR_JWT_TOKEN_HERE") {
    Write-Host "❌ ERROR: Please replace 'YOUR_JWT_TOKEN_HERE' with a valid JWT token." -ForegroundColor Red
    Write-Host "💡 How to get token:" -ForegroundColor Yellow
    Write-Host "   1. Login to the application" -ForegroundColor White
    Write-Host "   2. Open Developer Tools (F12)" -ForegroundColor White
    Write-Host "   3. Go to Network tab" -ForegroundColor White
    Write-Host "   4. Look for login request" -ForegroundColor White
    Write-Host "   5. Copy token from Authorization header" -ForegroundColor White
    Write-Host "`n   Example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." -ForegroundColor Gray
    exit
}

$Headers = @{
    "Authorization" = "Bearer $AuthToken"
    "Content-Type"  = "application/json"
}

# Test 1: Check if backend is running
Write-Host "`n1. Checking backend status..." -ForegroundColor Green

try {
    $Response = Invoke-RestMethod -Uri $ApiUrl -Method Get -Headers $Headers
    Write-Host "✅ Backend is running and accessible" -ForegroundColor Green
    Write-Host "👥 Current workers count: $($Response.Count)" -ForegroundColor White
}
catch {
    Write-Host "❌ Backend is not accessible" -ForegroundColor Red
    Write-Host "📋 Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "💡 Make sure backend is running on port 8080" -ForegroundColor Yellow
    exit
}

# Test 2: Test individual worker creation with detailed validation
Write-Host "`n2. Testing individual worker creation..." -ForegroundColor Green

# Data test yang sesuai dengan template Excel
$TestData = @{
    noPegawai = "PEG001"
    nama = "John Doe"
    gelar = "S.Kom"
    teleponAreaRumah = "021"
    teleponRumah = "7654321"
    teleponAreaKantor = "021"
    teleponKantor = "021-1234"
    teleponExtKantor = "1234"
    handphone = "081234567890"
    email = "john.doe@example.com"
    tempatLahir = "Jakarta"
    dateOfBirth = "1990-05-14"
    ibuKandung = "Siti Aminah"
    jenisIdentitas = "KTP"
    nik = "3175051405900001"
    masaLakuIdentitas = "2016-05-14"
    jenisKelamin = "Laki-laki"
    suratMenyuratKe = "Jalan Mandaka No. 10"
    tanggalKepesertaan = "2020-07-01"
    statusKawin = "Belum Kawin"
    golDarah = "O"
    npwp = "012345678901234"
    kodeNegara = "ID"
    upah = 5000000
    alamat = "Jl. Merpati 12 Jakarta"
    kodePos = "10130"
    lokasiPekerjaan = "Kantor Pusat"
    statusPegawai = "PKWTT"
    tanggalAwalBekerja = "2020-07-01"
    tanggalAkhirKontrak = ""
    rapel = 0
    nationality = "WNI"
} | ConvertTo-Json

try {
    Write-Host "🚀 Creating individual worker..." -ForegroundColor Yellow
    $Response = Invoke-RestMethod -Uri $ApiUrl -Method Post -Headers $Headers -Body $TestData
    
    Write-Host "✅ Individual worker created successfully!" -ForegroundColor Green
    Write-Host "👤 Worker ID: $($Response.worker.id)" -ForegroundColor White
    Write-Host "💰 Total Iuran: Rp $($Response.iuran.totalIuran.ToString('N0'))" -ForegroundColor White
}
catch {
    Write-Host "❌ Individual worker creation failed!" -ForegroundColor Red
    Write-Host "📋 Error: $($_.Exception.Message)" -ForegroundColor Red
    
    if ($_.Exception.Response) {
        $ErrorResponse = $_.Exception.Response.GetResponseStream()
        $Reader = New-Object System.IO.StreamReader($ErrorResponse)
        $ResponseBody = $Reader.ReadToEnd()
        Write-Host "📋 Response Body: $ResponseBody" -ForegroundColor Red
        
        # Parse error details
        try {
            $ErrorJson = $ResponseBody | ConvertFrom-Json
            if ($ErrorJson.details) {
                Write-Host "🔍 Validation Errors:" -ForegroundColor Yellow
                foreach ($Error in $ErrorJson.details) {
                    Write-Host "   - $Error" -ForegroundColor Red
                }
            }
        }
        catch {
            Write-Host "⚠️ Could not parse error details" -ForegroundColor Yellow
        }
    }
}

# Test 3: Mass upload with detailed error checking
Write-Host "`n3. Testing mass upload..." -ForegroundColor Green

$FilePath = "template_sipp_valid.xlsx"

if (-not (Test-Path $FilePath)) {
    Write-Host "❌ File not found: $FilePath" -ForegroundColor Red
    Write-Host "💡 Make sure template_sipp_valid.xlsx exists in current directory" -ForegroundColor Yellow
    exit
}

try {
    Write-Host "🚀 Uploading file: $FilePath" -ForegroundColor Yellow
    
    # Create multipart form data
    $FileBytes = [System.IO.File]::ReadAllBytes($FilePath)
    $Boundary = [System.Guid]::NewGuid().ToString()
    $NewLine = "`r`n"
    
    $Body = "--$Boundary$NewLine"
    $Body += "Content-Disposition: form-data; name=`"file`"; filename=`"$([System.IO.Path]::GetFileName($FilePath))`"$NewLine"
    $Body += "Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet$NewLine"
    $Body += $NewLine
    $Body += [System.Text.Encoding]::GetEncoding("iso-8859-1").GetString($FileBytes)
    $Body += "$NewLine--$Boundary--$NewLine"
    
    $UploadHeaders = @{
        "Authorization" = "Bearer $AuthToken"
        "Content-Type"  = "multipart/form-data; boundary=$Boundary"
    }
    
    $UploadResponse = Invoke-RestMethod -Uri $UploadUrl -Method Post -Headers $UploadHeaders -Body $Body -ContentType "multipart/form-data; boundary=$Boundary"
    
    Write-Host "✅ Upload completed!" -ForegroundColor Green
    Write-Host "📊 Total Data: $($UploadResponse.total)" -ForegroundColor White
    Write-Host "✅ Valid: $($UploadResponse.valid)" -ForegroundColor Green
    Write-Host "❌ Invalid: $($UploadResponse.invalid)" -ForegroundColor Red
    
    if ($UploadResponse.errors -and $UploadResponse.errors.Count -gt 0) {
        Write-Host "🔍 Upload Error Details:" -ForegroundColor Yellow
        foreach ($Error in $UploadResponse.errors) {
            Write-Host "   - $Error" -ForegroundColor Red
        }
    }
    
    # Check upload history for more details
    if ($UploadResponse.history_id) {
        Write-Host "📝 History ID: $($UploadResponse.history_id)" -ForegroundColor White
        
        $HistoryUrl = "http://localhost:8080/api/v1/upload-histories/$($UploadResponse.history_id)"
        try {
            $HistoryResponse = Invoke-RestMethod -Uri $HistoryUrl -Method Get -Headers $Headers
            
            Write-Host "📊 Upload History Details:" -ForegroundColor Cyan
            Write-Host "   - File Name: $($HistoryResponse.file_name)" -ForegroundColor White
            Write-Host "   - Total Data: $($HistoryResponse.total_data)" -ForegroundColor White
            Write-Host "   - Valid: $($HistoryResponse.total_valid)" -ForegroundColor Green
            Write-Host "   - Invalid: $($HistoryResponse.total_invalid)" -ForegroundColor Red
            Write-Host "   - Status: $($HistoryResponse.validation_status)" -ForegroundColor White
            
            if ($HistoryResponse.errors -and $HistoryResponse.errors.Count -gt 0) {
                Write-Host "🔍 Detailed Validation Errors:" -ForegroundColor Yellow
                foreach ($Error in $HistoryResponse.errors) {
                    Write-Host "   - $Error" -ForegroundColor Red
                }
            }
        }
        catch {
            Write-Host "⚠️ Could not fetch upload history details" -ForegroundColor Yellow
        }
    }
}
catch {
    Write-Host "❌ Upload failed!" -ForegroundColor Red
    Write-Host "📋 Error: $($_.Exception.Message)" -ForegroundColor Red
    
    if ($_.Exception.Response) {
        $ErrorResponse = $_.Exception.Response.GetResponseStream()
        $Reader = New-Object System.IO.StreamReader($ErrorResponse)
        $ResponseBody = $Reader.ReadToEnd()
        Write-Host "📋 Response Body: $ResponseBody" -ForegroundColor Red
    }
}

# Test 4: Check database schema
Write-Host "`n4. Checking database schema..." -ForegroundColor Green

try {
    # Try to get a worker to check if new fields exist
    $WorkersResponse = Invoke-RestMethod -Uri $ApiUrl -Method Get -Headers $Headers
    
    if ($WorkersResponse.Count -gt 0) {
        $FirstWorker = $WorkersResponse[0]
        Write-Host "✅ Database connection working" -ForegroundColor Green
        
        # Check if new fields exist
        $NewFields = @(
            "telepon_area_rumah",
            "telepon_area_kantor", 
            "telepon_kantor",
            "telepon_ext_kantor",
            "jenis_identitas",
            "masa_laku_identitas",
            "surat_menyurat_ke",
            "tanggal_kepesertaan",
            "kode_negara"
        )
        
        Write-Host "🔍 Checking new SIPP fields:" -ForegroundColor Cyan
        foreach ($Field in $NewFields) {
            if ($FirstWorker.PSObject.Properties.Name -contains $Field) {
                Write-Host "   ✅ $Field exists" -ForegroundColor Green
            } else {
                Write-Host "   ❌ $Field missing" -ForegroundColor Red
            }
        }
    }
}
catch {
    Write-Host "⚠️ Could not check database schema" -ForegroundColor Yellow
}

Write-Host "`n" + "=" * 50 -ForegroundColor Cyan
Write-Host "🎯 Summary & Recommendations:" -ForegroundColor Green

Write-Host "`n💡 Jika masih 'Valid: 0':" -ForegroundColor Yellow
Write-Host "   1. Jalankan: .\fix_validation_quick.ps1" -ForegroundColor White
Write-Host "   2. Restart backend" -ForegroundColor White
Write-Host "   3. Coba upload lagi" -ForegroundColor White

Write-Host "`n🔧 Jika ada field yang missing:" -ForegroundColor Yellow
Write-Host "   1. Jalankan: .\run_migration.ps1" -ForegroundColor White
Write-Host "   2. Restart backend" -ForegroundColor White

Write-Host "`n📋 Jika ada error validasi spesifik:" -ForegroundColor Yellow
Write-Host "   1. Perbaiki data di Excel sesuai error" -ForegroundColor White
Write-Host "   2. Upload ulang" -ForegroundColor White

Write-Host "`n🚀 Quick Fix Commands:" -ForegroundColor Cyan
Write-Host "   # Reduce validation" -ForegroundColor White
Write-Host "   .\fix_validation_quick.ps1" -ForegroundColor White
Write-Host "   " -ForegroundColor White
Write-Host "   # Run migration" -ForegroundColor White
Write-Host "   .\run_migration.ps1" -ForegroundColor White
Write-Host "   " -ForegroundColor White
Write-Host "   # Restart backend" -ForegroundColor White
Write-Host "   cd backend && go run main.go" -ForegroundColor White







