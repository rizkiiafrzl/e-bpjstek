# PowerShell script untuk menjalankan migration database e-BPJSTK
Write-Host "🚀 Starting Database Migration for e-BPJSTK..." -ForegroundColor Green

# Set environment variables (sesuaikan dengan konfigurasi Anda)
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "password"
$env:DB_NAME = "e-bpjstk"
$env:DB_SSLMODE = "disable"

# Navigate to backend directory
Set-Location backend
Write-Host "📁 Changed to backend directory" -ForegroundColor Yellow

# Run the migration script
Write-Host "🔄 Running database migration..." -ForegroundColor Cyan
go run ../migrate_database.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Migration completed successfully!" -ForegroundColor Green
} else {
    Write-Host "❌ Migration failed!" -ForegroundColor Red
    exit 1
}

# Optional: Run database check (jika psql tersedia)
Write-Host "🔍 Running database structure check..." -ForegroundColor Cyan
try {
    psql -h $env:DB_HOST -p $env:DB_PORT -U $env:DB_USER -d $env:DB_NAME -f ../check_database.sql
} catch {
    Write-Host "⚠️ psql not found, skipping database check" -ForegroundColor Yellow
    Write-Host "You can manually run the check_database.sql file" -ForegroundColor Yellow
}

Write-Host "🎉 Database migration process completed!" -ForegroundColor Green







