# Check Backend Logs untuk Debug Validasi
# Script PowerShell untuk melihat log backend

Write-Host "🔍 Checking Backend Logs for Validation Issues" -ForegroundColor Cyan
Write-Host "=" * 60 -ForegroundColor Cyan

# Check if backend is running
Write-Host "`n1. Checking if backend is running..." -ForegroundColor Green

try {
    $Response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/workers" -Method Get -TimeoutSec 5
    Write-Host "✅ Backend is running on port 8080" -ForegroundColor Green
    Write-Host "📊 Status Code: $($Response.StatusCode)" -ForegroundColor White
}
catch {
    Write-Host "❌ Backend is not running or not accessible" -ForegroundColor Red
    Write-Host "💡 Please start the backend first:" -ForegroundColor Yellow
    Write-Host "   cd backend" -ForegroundColor White
    Write-Host "   go run main.go" -ForegroundColor White
    exit
}

# Check backend process
Write-Host "`n2. Checking backend process..." -ForegroundColor Green

$BackendProcess = Get-Process | Where-Object { $_.ProcessName -like "*go*" -or $_.ProcessName -like "*main*" -or $_.CommandLine -like "*main.go*" }

if ($BackendProcess) {
    Write-Host "✅ Found backend process:" -ForegroundColor Green
    $BackendProcess | ForEach-Object {
        Write-Host "   - PID: $($_.Id), Name: $($_.ProcessName)" -ForegroundColor White
    }
} else {
    Write-Host "⚠️ Backend process not found in process list" -ForegroundColor Yellow
}

# Check if log files exist
Write-Host "`n3. Checking for log files..." -ForegroundColor Green

$LogFiles = @(
    "backend.log",
    "app.log",
    "error.log",
    "debug.log"
)

$FoundLogs = @()
foreach ($LogFile in $LogFiles) {
    if (Test-Path $LogFile) {
        $FoundLogs += $LogFile
        Write-Host "✅ Found log file: $LogFile" -ForegroundColor Green
    }
}

if ($FoundLogs.Count -eq 0) {
    Write-Host "⚠️ No log files found in current directory" -ForegroundColor Yellow
    Write-Host "💡 Backend might be logging to console only" -ForegroundColor White
}

# Check recent log entries
if ($FoundLogs.Count -gt 0) {
    Write-Host "`n4. Recent log entries (last 20 lines):" -ForegroundColor Green
    
    foreach ($LogFile in $FoundLogs) {
        Write-Host "📋 Log file: $LogFile" -ForegroundColor Cyan
        try {
            $LogContent = Get-Content $LogFile -Tail 20
            $LogContent | ForEach-Object {
                if ($_ -match "ERROR|error|Error") {
                    Write-Host "   ❌ $_" -ForegroundColor Red
                } elseif ($_ -match "WARN|warn|Warn") {
                    Write-Host "   ⚠️ $_" -ForegroundColor Yellow
                } elseif ($_ -match "DEBUG|debug|Debug") {
                    Write-Host "   🔍 $_" -ForegroundColor Cyan
                } else {
                    Write-Host "   ℹ️ $_" -ForegroundColor White
                }
            }
        }
        catch {
            Write-Host "   ❌ Could not read log file: $LogFile" -ForegroundColor Red
        }
        Write-Host ""
    }
}

# Check database connection
Write-Host "`n5. Testing database connection..." -ForegroundColor Green

try {
    # Try to get workers list to test database
    $Headers = @{
        "Authorization" = "Bearer YOUR_JWT_TOKEN_HERE"
        "Content-Type"  = "application/json"
    }
    
    $Response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/workers" -Method Get -Headers $Headers
    Write-Host "✅ Database connection is working" -ForegroundColor Green
    Write-Host "📊 Found $($Response.Count) workers in database" -ForegroundColor White
}
catch {
    Write-Host "⚠️ Could not test database connection (might need valid token)" -ForegroundColor Yellow
    Write-Host "📋 Error: $($_.Exception.Message)" -ForegroundColor Red
}

# Check upload history
Write-Host "`n6. Checking recent upload history..." -ForegroundColor Green

try {
    $Headers = @{
        "Authorization" = "Bearer YOUR_JWT_TOKEN_HERE"
        "Content-Type"  = "application/json"
    }
    
    $Response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/upload-histories" -Method Get -Headers $Headers
    
    if ($Response.data -and $Response.data.Count -gt 0) {
        Write-Host "✅ Found $($Response.data.Count) upload histories" -ForegroundColor Green
        
        $Response.data | Select-Object -First 3 | ForEach-Object {
            Write-Host "📋 Upload: $($_.file_name)" -ForegroundColor White
            Write-Host "   - Total: $($_.total_data), Valid: $($_.total_valid), Invalid: $($_.total_invalid)" -ForegroundColor White
            Write-Host "   - Status: $($_.validation_status)" -ForegroundColor White
            Write-Host "   - Date: $($_.created_at)" -ForegroundColor White
            
            if ($_.errors -and $_.errors.Count -gt 0) {
                Write-Host "   - Errors:" -ForegroundColor Red
                $_.errors | ForEach-Object {
                    Write-Host "     * $_" -ForegroundColor Red
                }
            }
            Write-Host ""
        }
    } else {
        Write-Host "⚠️ No upload histories found" -ForegroundColor Yellow
    }
}
catch {
    Write-Host "⚠️ Could not fetch upload history (might need valid token)" -ForegroundColor Yellow
}

# Recommendations
Write-Host "`n7. Recommendations for debugging:" -ForegroundColor Green

Write-Host "💡 To debug validation issues:" -ForegroundColor Yellow
Write-Host "   1. Check backend console output for validation errors" -ForegroundColor White
Write-Host "   2. Use debug_upload_sipp.ps1 with valid JWT token" -ForegroundColor White
Write-Host "   3. Check if database migration was run successfully" -ForegroundColor White
Write-Host "   4. Verify Excel file format matches SIPP template exactly" -ForegroundColor White
Write-Host "   5. Check if all required fields are filled in Excel" -ForegroundColor White

Write-Host "`n🔧 Common validation issues:" -ForegroundColor Yellow
Write-Host "   - NIK must be exactly 16 digits" -ForegroundColor White
Write-Host "   - Email must be valid format" -ForegroundColor White
Write-Host "   - NPWP must be exactly 15 digits" -ForegroundColor White
Write-Host "   - Upah must be >= 1,000,000" -ForegroundColor White
Write-Host "   - PKWT requires tanggal akhir kontrak" -ForegroundColor White
Write-Host "   - Jenis kelamin must be 'Laki-laki' or 'Perempuan'" -ForegroundColor White
Write-Host "   - Golongan darah must be A, B, AB, or O" -ForegroundColor White

Write-Host "`n" + "=" * 60 -ForegroundColor Cyan
Write-Host "🎯 Next Steps:" -ForegroundColor Green
Write-Host "   1. Run debug_upload_sipp.ps1 with valid token" -ForegroundColor White
Write-Host "   2. Check backend console for detailed error messages" -ForegroundColor White
Write-Host "   3. Verify Excel file format and data" -ForegroundColor White
Write-Host "   4. Run database migration if needed" -ForegroundColor White







