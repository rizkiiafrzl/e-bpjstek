# PowerShell script untuk memperbaiki fungsi UploadTK
Write-Host "🔧 Patching UploadTK function..." -ForegroundColor Green

# Backup file asli
Copy-Item "backend/handlers/worker.go" "backend/handlers/worker.go.backup"
Write-Host "✅ Backup created: worker.go.backup" -ForegroundColor Green

# Tambahkan import fmt jika belum ada
$workerFile = "backend/handlers/worker.go"
$content = Get-Content $workerFile -Raw

if ($content -notmatch '"fmt"') {
    Write-Host "📝 Adding fmt import..." -ForegroundColor Yellow
    $content = $content -replace 'import \(', 'import (`n`t"fmt"'
    $content = $content -replace '"log"', '"log"`n`t'
    Set-Content $workerFile $content
    Write-Host "✅ fmt import added" -ForegroundColor Green
}

# Tambahkan fungsi min helper jika belum ada
if ($content -notmatch 'func min\(') {
    Write-Host "📝 Adding min helper function..." -ForegroundColor Yellow
    $minFunction = @"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
"@
    $content = $content + $minFunction
    Set-Content $workerFile $content
    Write-Host "✅ min function added" -ForegroundColor Green
}

Write-Host "🎉 Patch completed!" -ForegroundColor Green
Write-Host ""
Write-Host "📋 Next steps:" -ForegroundColor Yellow
Write-Host "1. Restart backend: cd backend; go run main.go" -ForegroundColor White
Write-Host "2. Test upload Excel dengan data valid" -ForegroundColor White
Write-Host "3. Cek log backend untuk debug info" -ForegroundColor White
Write-Host "4. Verifikasi data muncul di tabel" -ForegroundColor White
Write-Host ""
Write-Host "🔍 Debug info akan muncul di log backend:" -ForegroundColor Cyan
Write-Host "   - Total rows in Excel" -ForegroundColor Gray
Write-Host "   - Header row content" -ForegroundColor Gray
Write-Host "   - First data row content" -ForegroundColor Gray
Write-Host "   - Processing each row" -ForegroundColor Gray
Write-Host "   - Validation errors" -ForegroundColor Gray
Write-Host "   - Success/failure for each row" -ForegroundColor Gray







