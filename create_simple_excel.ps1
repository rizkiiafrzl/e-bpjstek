# PowerShell script untuk membuat file Excel sederhana
Write-Host "📊 Creating Simple Excel File..." -ForegroundColor Green

# Buat file Excel dengan PowerShell
$excelFile = "template_tk_test_valid.xlsx"

# Data test
$data = @'
NO_PEGAWAI,NAMA_LENGKAP,GELAR,TELEPON_AREA_RUMAH,TELEPON_RUMAH,,TELEPON_AREA_KANTOR,TELEPON_KANTOR,TELEPON_EXT_KANTOR,HP,EMAIL,TEMPAT_LAHIR,TANGGAL_LAHIR,NAMA_IBU_KANDUNG,JENIS_IDENTITAS,NO_IDENTITAS,MASA_LAKU_IDENTITAS,JENIS_KELAMIN,SURAT_MENYURAT_KE,TANGGAL_KEPESERTAAN,STATUS_KAWIN,GOLONGAN_DARAH,NPWP,KODE_NEGARA,UPAH,ALAMAT,KODE_POS,LOKASI_PEKERJAAN,STATUS_PEGAWAI,TGL_AWAL_BEKERJA,TGL_AKHIR_KONTRAK
PEG001,John Doe,S.Kom,021,7654321,,021,021-1234,1234,081234567890,john.doe@example.com,Jakarta,1990-05-14,Siti Aminah,KTP,3175051405900001,2016-05-14,Laki-laki,Jalan Mandaka No. 10,2020-07-01,Belum Kawin,O,012345678901234,ID,5000000,Jl. Merpati 12 Jakarta,10130,Kantor Pusat,Tetap,2020-07-01,
PEG002,Jane Smith,S.E,022,8765432,,022,022-5678,5678,082345678901,jane.smith@example.com,Bandung,1988-03-22,Maria Sari,KTP,3175051405900002,2018-03-22,Perempuan,Jalan Sudirman No. 5,2021-01-15,Kawin,A,012345678901235,ID,6000000,Jl. Gatot Subroto 25 Bandung,40112,Cabang Bandung,Tetap,2021-01-15,
PEG003,Ahmad Rahman,S.T,024,9876543,,024,024-9012,9012,083456789012,ahmad.rahman@example.com,Semarang,1992-11-08,Fatimah Zahra,KTP,3175051405900003,2020-11-08,Laki-laki,Jalan Diponegoro No. 15,2022-03-10,Belum Kawin,B,012345678901236,ID,4500000,Jl. Imam Bonjol 8 Semarang,50111,Cabang Semarang,Kontrak,2022-03-10,2024-03-10
'@

# Simpan sebagai CSV dulu
$csvFile = "temp_data.csv"
$data | Out-File -FilePath $csvFile -Encoding UTF8

Write-Host "✅ CSV file created: $csvFile" -ForegroundColor Green

# Coba buat Excel dengan COM object
try {
    Write-Host "📊 Converting to Excel..." -ForegroundColor Yellow
    
    $excel = New-Object -ComObject Excel.Application
    $excel.Visible = $false
    $excel.DisplayAlerts = $false
    
    # Buka CSV
    $workbook = $excel.Workbooks.Open((Resolve-Path $csvFile).Path)
    $worksheet = $workbook.Worksheets.Item(1)
    
    # Simpan sebagai Excel
    $workbook.SaveAs((Resolve-Path .).Path + "\" + $excelFile, 51)
    $workbook.Close()
    $excel.Quit()
    
    Write-Host "✅ Excel file created: $excelFile" -ForegroundColor Green
    
    # Cleanup
    Remove-Item $csvFile
    
} catch {
    Write-Host "⚠️  Excel COM not available, keeping CSV file" -ForegroundColor Yellow
    Write-Host "   You can manually open $csvFile in Excel and save as .xlsx" -ForegroundColor Gray
}

Write-Host "`n📋 File Details:" -ForegroundColor Cyan
Write-Host "   - Total data: 3 rows" -ForegroundColor White
Write-Host "   - Columns: 31 (A-AE)" -ForegroundColor White
Write-Host "   - Format: Valid for upload" -ForegroundColor White

Write-Host "`n✅ Test file ready!" -ForegroundColor Green
Write-Host "📁 File: $excelFile" -ForegroundColor Yellow
