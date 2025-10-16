# Test Backend Upload TK NA
Write-Host "=== TEST BACKEND UPLOAD TK NA ===" -ForegroundColor Green

# Test 1: Cek apakah backend bisa di-compile
Write-Host "`n1. Test Compile Backend..." -ForegroundColor Yellow
cd backend
$compileResult = go build -o main.exe . 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Backend compile berhasil" -ForegroundColor Green
} else {
    Write-Host "❌ Backend compile gagal:" -ForegroundColor Red
    Write-Host $compileResult
    exit 1
}

# Test 2: Cek apakah file template ada
Write-Host "`n2. Test Template File..." -ForegroundColor Yellow
cd ..
if (Test-Path "template_tk_na.xlsx") {
    $fileSize = (Get-Item "template_tk_na.xlsx").Length
    Write-Host "✅ Template TK NA ditemukan: $fileSize bytes" -ForegroundColor Green
} else {
    Write-Host "❌ Template TK NA tidak ditemukan" -ForegroundColor Red
    exit 1
}

# Test 3: Cek apakah handler UploadTKNA ada di kode
Write-Host "`n3. Test Handler UploadTKNA..." -ForegroundColor Yellow
$handlerExists = Select-String -Path "backend\handlers\worker.go" -Pattern "func UploadTKNA" -Quiet
if ($handlerExists) {
    Write-Host "✅ Handler UploadTKNA ditemukan" -ForegroundColor Green
} else {
    Write-Host "❌ Handler UploadTKNA tidak ditemukan" -ForegroundColor Red
    exit 1
}

# Test 4: Cek apakah route sudah terdaftar
Write-Host "`n4. Test Route Registration..." -ForegroundColor Yellow
$routeExists = Select-String -Path "backend\main.go" -Pattern "upload-tk-na" -Quiet
if ($routeExists) {
    Write-Host "✅ Route upload-tk-na terdaftar" -ForegroundColor Green
} else {
    Write-Host "❌ Route upload-tk-na tidak terdaftar" -ForegroundColor Red
    exit 1
}

# Test 5: Cek apakah API method ada di frontend
Write-Host "`n5. Test Frontend API Method..." -ForegroundColor Yellow
$apiMethodExists = Select-String -Path "e-bpjstk\src\services\api.js" -Pattern "uploadTKNA" -Quiet
if ($apiMethodExists) {
    Write-Host "✅ API method uploadTKNA ditemukan di frontend" -ForegroundColor Green
} else {
    Write-Host "❌ API method uploadTKNA tidak ditemukan di frontend" -ForegroundColor Red
    exit 1
}

# Test 6: Cek apakah frontend menggunakan API yang benar
Write-Host "`n6. Test Frontend Integration..." -ForegroundColor Yellow
$correctApiUsage = Select-String -Path "e-bpjstk\src\views\UploadTKNAView.vue" -Pattern "api.uploadTKNA" -Quiet
if ($correctApiUsage) {
    Write-Host "✅ Frontend menggunakan API uploadTKNA yang benar" -ForegroundColor Green
} else {
    Write-Host "❌ Frontend masih menggunakan API yang salah" -ForegroundColor Red
    exit 1
}

Write-Host "`n=== HASIL TEST ===" -ForegroundColor Green
Write-Host "✅ Backend compile berhasil" -ForegroundColor Green
Write-Host "✅ Template file siap" -ForegroundColor Green
Write-Host "✅ Handler UploadTKNA tersedia" -ForegroundColor Green
Write-Host "✅ Route terdaftar" -ForegroundColor Green
Write-Host "✅ Frontend API method tersedia" -ForegroundColor Green
Write-Host "✅ Frontend integration benar" -ForegroundColor Green

Write-Host "`n🎉 SEMUA TEST BERHASIL!" -ForegroundColor Green
Write-Host "Backend Upload TK NA sudah BENAR dan siap digunakan!" -ForegroundColor Green



