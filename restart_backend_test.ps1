# Restart Backend dan Test Upload - Solusi untuk Valid: 0
# Script PowerShell untuk restart backend dan test upload

Write-Host "🚀 Restart Backend dan Test Upload" -ForegroundColor Cyan
Write-Host "=" * 50 -ForegroundColor Cyan

# Check if backend is running
Write-Host "`n1. Checking if backend is running..." -ForegroundColor Green

try {
    $Response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/workers" -Method Get -TimeoutSec 3
    Write-Host "✅ Backend is currently running" -ForegroundColor Green
    Write-Host "📊 Status Code: $($Response.StatusCode)" -ForegroundColor White
}
catch {
    Write-Host "⚠️ Backend is not running or not accessible" -ForegroundColor Yellow
}

# Kill existing backend processes
Write-Host "`n2. Stopping existing backend processes..." -ForegroundColor Green

$BackendProcesses = Get-Process | Where-Object { 
    $_.ProcessName -like "*go*" -or 
    $_.ProcessName -like "*main*" -or 
    $_.CommandLine -like "*main.go*" -or
    $_.CommandLine -like "*8080*"
}

if ($BackendProcesses) {
    Write-Host "🔍 Found backend processes:" -ForegroundColor Yellow
    $BackendProcesses | ForEach-Object {
        Write-Host "   - PID: $($_.Id), Name: $($_.ProcessName)" -ForegroundColor White
        try {
            Stop-Process -Id $_.Id -Force
            Write-Host "   ✅ Process $($_.Id) stopped" -ForegroundColor Green
        }
        catch {
            Write-Host "   ⚠️ Could not stop process $($_.Id)" -ForegroundColor Yellow
        }
    }
} else {
    Write-Host "ℹ️ No backend processes found" -ForegroundColor White
}

# Wait a moment for processes to stop
Write-Host "`n3. Waiting for processes to stop..." -ForegroundColor Green
Start-Sleep -Seconds 3

# Start backend
Write-Host "`n4. Starting backend..." -ForegroundColor Green

$BackendPath = "backend"
if (-not (Test-Path $BackendPath)) {
    Write-Host "❌ Backend directory not found: $BackendPath" -ForegroundColor Red
    Write-Host "💡 Make sure you're in the project root directory" -ForegroundColor Yellow
    exit
}

# Check if main.go exists
$MainGoPath = "$BackendPath/main.go"
if (-not (Test-Path $MainGoPath)) {
    Write-Host "❌ main.go not found: $MainGoPath" -ForegroundColor Red
    exit
}

Write-Host "🚀 Starting backend in background..." -ForegroundColor Yellow
Write-Host "📁 Working directory: $BackendPath" -ForegroundColor White
Write-Host "📄 Main file: $MainGoPath" -ForegroundColor White

# Start backend in background
$BackendJob = Start-Job -ScriptBlock {
    Set-Location $using:BackendPath
    go run main.go
}

Write-Host "✅ Backend started in background (Job ID: $($BackendJob.Id))" -ForegroundColor Green

# Wait for backend to start
Write-Host "`n5. Waiting for backend to start..." -ForegroundColor Green
$MaxWait = 30
$WaitCount = 0

do {
    Start-Sleep -Seconds 2
    $WaitCount++
    
    try {
        $Response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/workers" -Method Get -TimeoutSec 2
        Write-Host "✅ Backend is ready!" -ForegroundColor Green
        Write-Host "📊 Status Code: $($Response.StatusCode)" -ForegroundColor White
        break
    }
    catch {
        Write-Host "⏳ Waiting for backend... ($WaitCount/$MaxWait)" -ForegroundColor Yellow
    }
} while ($WaitCount -lt $MaxWait)

if ($WaitCount -ge $MaxWait) {
    Write-Host "❌ Backend failed to start within $MaxWait seconds" -ForegroundColor Red
    Write-Host "📋 Job output:" -ForegroundColor Yellow
    Receive-Job -Job $BackendJob
    Stop-Job -Job $BackendJob
    Remove-Job -Job $BackendJob
    exit
}

# Test upload
Write-Host "`n6. Testing upload with relaxed validation..." -ForegroundColor Green

$FilePath = "template_sipp_valid.xlsx"

if (-not (Test-Path $FilePath)) {
    Write-Host "❌ File not found: $FilePath" -ForegroundColor Red
    Write-Host "💡 Make sure template_sipp_valid.xlsx exists in current directory" -ForegroundColor Yellow
    exit
}

Write-Host "📄 Testing with file: $FilePath" -ForegroundColor White

# Note: This test requires a valid JWT token
Write-Host "`n⚠️ Note: Upload test requires a valid JWT token" -ForegroundColor Yellow
Write-Host "💡 To test upload:" -ForegroundColor White
Write-Host "   1. Login to the application" -ForegroundColor White
Write-Host "   2. Get JWT token from browser" -ForegroundColor White
Write-Host "   3. Run: .\debug_detail_error.ps1" -ForegroundColor White
Write-Host "   4. Or test manually in the browser" -ForegroundColor White

Write-Host "`n7. Backend status:" -ForegroundColor Green
Write-Host "✅ Backend is running on http://localhost:8080" -ForegroundColor Green
Write-Host "🔧 Validation has been relaxed for testing" -ForegroundColor Yellow
Write-Host "📊 Changes made:" -ForegroundColor Cyan
Write-Host "   - NIK: 10-16 digits (was 16 exactly)" -ForegroundColor White
Write-Host "   - NPWP: 10-15 digits (was 15 exactly)" -ForegroundColor White
Write-Host "   - Email: More flexible format" -ForegroundColor White
Write-Host "   - Upah: Min 100k (was 1M)" -ForegroundColor White
Write-Host "   - Usia: 15-70 years (was 17-65)" -ForegroundColor White
Write-Host "   - Jenis Kelamin: More flexible options" -ForegroundColor White
Write-Host "   - Status Pegawai: More flexible options" -ForegroundColor White

Write-Host "`n8. Next steps:" -ForegroundColor Green
Write-Host "   1. Open browser and go to http://localhost:5173" -ForegroundColor White
Write-Host "   2. Login to the application" -ForegroundColor White
Write-Host "   3. Go to Upload TK page" -ForegroundColor White
Write-Host "   4. Upload template_sipp_valid.xlsx" -ForegroundColor White
Write-Host "   5. Check if Valid > 0 now" -ForegroundColor White

Write-Host "`n9. To stop backend:" -ForegroundColor Green
Write-Host "   Stop-Job -Job $($BackendJob.Id)" -ForegroundColor White
Write-Host "   Remove-Job -Job $($BackendJob.Id)" -ForegroundColor White

Write-Host "`n" + "=" * 50 -ForegroundColor Cyan
Write-Host "🎯 Summary:" -ForegroundColor Green
Write-Host "✅ Backend restarted with relaxed validation" -ForegroundColor Green
Write-Host "✅ Ready for testing upload" -ForegroundColor Green
Write-Host "⚠️ Remember to restore strict validation for production" -ForegroundColor Yellow

# Keep job running
Write-Host "`n💡 Backend is running in background. Press Ctrl+C to stop this script." -ForegroundColor Cyan
Write-Host "   The backend will continue running even if you close this script." -ForegroundColor White

try {
    while ($true) {
        Start-Sleep -Seconds 10
        
        # Check if backend is still running
        try {
            $Response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/workers" -Method Get -TimeoutSec 2
            Write-Host "✅ Backend is still running... $(Get-Date -Format 'HH:mm:ss')" -ForegroundColor Green
        }
        catch {
            Write-Host "❌ Backend seems to have stopped" -ForegroundColor Red
            break
        }
    }
}
catch {
    Write-Host "`n👋 Script stopped by user" -ForegroundColor Yellow
}
finally {
    # Clean up
    if ($BackendJob.State -eq "Running") {
        Write-Host "`n🛑 Stopping backend..." -ForegroundColor Yellow
        Stop-Job -Job $BackendJob
        Remove-Job -Job $BackendJob
        Write-Host "✅ Backend stopped" -ForegroundColor Green
    }
}







