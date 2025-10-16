# PowerShell script untuk membuat file Excel test
Write-Host "📊 Creating Test Excel File..." -ForegroundColor Green

# Data test dengan format yang benar
$testData = @(
    @{
        'NO_PEGAWAI' = 'PEG001'
        'NAMA_LENGKAP' = 'John Doe'
        'GELAR' = 'S.Kom'
        'TELEPON_AREA_RUMAH' = '021'
        'TELEPON_RUMAH' = '7654321'
        '' = ''  # Kolom F kosong
        'TELEPON_AREA_KANTOR' = '021'
        'TELEPON_KANTOR' = '021-1234'
        'TELEPON_EXT_KANTOR' = '1234'
        'HP' = '081234567890'
        'EMAIL' = 'john.doe@example.com'
        'TEMPAT_LAHIR' = 'Jakarta'
        'TANGGAL_LAHIR' = '1990-05-14'
        'NAMA_IBU_KANDUNG' = 'Siti Aminah'
        'JENIS_IDENTITAS' = 'KTP'
        'NO_IDENTITAS' = '3175051405900001'
        'MASA_LAKU_IDENTITAS' = '2016-05-14'
        'JENIS_KELAMIN' = 'Laki-laki'
        'SURAT_MENYURAT_KE' = 'Jalan Mandaka No. 10'
        'TANGGAL_KEPESERTAAN' = '2020-07-01'
        'STATUS_KAWIN' = 'Belum Kawin'
        'GOLONGAN_DARAH' = 'O'
        'NPWP' = '012345678901234'
        'KODE_NEGARA' = 'ID'
        'UPAH' = '5000000'
        'ALAMAT' = 'Jl. Merpati 12 Jakarta'
        'KODE_POS' = '10130'
        'LOKASI_PEKERJAAN' = 'Kantor Pusat'
        'STATUS_PEGAWAI' = 'Tetap'
        'TGL_AWAL_BEKERJA' = '2020-07-01'
        'TGL_AKHIR_KONTRAK' = ''
    },
    @{
        'NO_PEGAWAI' = 'PEG002'
        'NAMA_LENGKAP' = 'Jane Smith'
        'GELAR' = 'S.E'
        'TELEPON_AREA_RUMAH' = '022'
        'TELEPON_RUMAH' = '8765432'
        '' = ''  # Kolom F kosong
        'TELEPON_AREA_KANTOR' = '022'
        'TELEPON_KANTOR' = '022-5678'
        'TELEPON_EXT_KANTOR' = '5678'
        'HP' = '082345678901'
        'EMAIL' = 'jane.smith@example.com'
        'TEMPAT_LAHIR' = 'Bandung'
        'TANGGAL_LAHIR' = '1988-03-22'
        'NAMA_IBU_KANDUNG' = 'Maria Sari'
        'JENIS_IDENTITAS' = 'KTP'
        'NO_IDENTITAS' = '3175051405900002'
        'MASA_LAKU_IDENTITAS' = '2018-03-22'
        'JENIS_KELAMIN' = 'Perempuan'
        'SURAT_MENYURAT_KE' = 'Jalan Sudirman No. 5'
        'TANGGAL_KEPESERTAAN' = '2021-01-15'
        'STATUS_KAWIN' = 'Kawin'
        'GOLONGAN_DARAH' = 'A'
        'NPWP' = '012345678901235'
        'KODE_NEGARA' = 'ID'
        'UPAH' = '6000000'
        'ALAMAT' = 'Jl. Gatot Subroto 25 Bandung'
        'KODE_POS' = '40112'
        'LOKASI_PEKERJAAN' = 'Cabang Bandung'
        'STATUS_PEGAWAI' = 'Tetap'
        'TGL_AWAL_BEKERJA' = '2021-01-15'
        'TGL_AKHIR_KONTRAK' = ''
    },
    @{
        'NO_PEGAWAI' = 'PEG003'
        'NAMA_LENGKAP' = 'Ahmad Rahman'
        'GELAR' = 'S.T'
        'TELEPON_AREA_RUMAH' = '024'
        'TELEPON_RUMAH' = '9876543'
        '' = ''  # Kolom F kosong
        'TELEPON_AREA_KANTOR' = '024'
        'TELEPON_KANTOR' = '024-9012'
        'TELEPON_EXT_KANTOR' = '9012'
        'HP' = '083456789012'
        'EMAIL' = 'ahmad.rahman@example.com'
        'TEMPAT_LAHIR' = 'Semarang'
        'TANGGAL_LAHIR' = '1992-11-08'
        'NAMA_IBU_KANDUNG' = 'Fatimah Zahra'
        'JENIS_IDENTITAS' = 'KTP'
        'NO_IDENTITAS' = '3175051405900003'
        'MASA_LAKU_IDENTITAS' = '2020-11-08'
        'JENIS_KELAMIN' = 'Laki-laki'
        'SURAT_MENYURAT_KE' = 'Jalan Diponegoro No. 15'
        'TANGGAL_KEPESERTAAN' = '2022-03-10'
        'STATUS_KAWIN' = 'Belum Kawin'
        'GOLONGAN_DARAH' = 'B'
        'NPWP' = '012345678901236'
        'KODE_NEGARA' = 'ID'
        'UPAH' = '4500000'
        'ALAMAT' = 'Jl. Imam Bonjol 8 Semarang'
        'KODE_POS' = '50111'
        'LOKASI_PEKERJAAN' = 'Cabang Semarang'
        'STATUS_PEGAWAI' = 'Kontrak'
        'TGL_AWAL_BEKERJA' = '2022-03-10'
        'TGL_AKHIR_KONTRAK' = '2024-03-10'
    }
)

# Buat CSV terlebih dahulu
$csvFile = "template_tk_test_valid.csv"
Write-Host "📝 Creating CSV file: $csvFile" -ForegroundColor Yellow

# Header
$headers = @(
    'NO_PEGAWAI', 'NAMA_LENGKAP', 'GELAR', 'TELEPON_AREA_RUMAH', 'TELEPON_RUMAH', '',
    'TELEPON_AREA_KANTOR', 'TELEPON_KANTOR', 'TELEPON_EXT_KANTOR', 'HP', 'EMAIL',
    'TEMPAT_LAHIR', 'TANGGAL_LAHIR', 'NAMA_IBU_KANDUNG', 'JENIS_IDENTITAS', 'NO_IDENTITAS',
    'MASA_LAKU_IDENTITAS', 'JENIS_KELAMIN', 'SURAT_MENYURAT_KE', 'TANGGAL_KEPESERTAAN',
    'STATUS_KAWIN', 'GOLONGAN_DARAH', 'NPWP', 'KODE_NEGARA', 'UPAH', 'ALAMAT',
    'KODE_POS', 'LOKASI_PEKERJAAN', 'STATUS_PEGAWAI', 'TGL_AWAL_BEKERJA', 'TGL_AKHIR_KONTRAK'
)

# Buat CSV content
$csvContent = $headers -join ','
$csvContent += "`n"

foreach ($row in $testData) {
    $rowValues = @()
    foreach ($header in $headers) {
        if ($header -eq '') {
            $rowValues += ''
        } else {
            $rowValues += $row[$header]
        }
    }
    $csvContent += ($rowValues -join ',') + "`n"
}

# Simpan CSV
$csvContent | Out-File -FilePath $csvFile -Encoding UTF8
Write-Host "✅ CSV file created: $csvFile" -ForegroundColor Green

# Cek apakah ada Excel COM object
try {
    Write-Host "📊 Converting CSV to Excel..." -ForegroundColor Yellow
    
    $excel = New-Object -ComObject Excel.Application
    $excel.Visible = $false
    $excel.DisplayAlerts = $false
    
    # Buka CSV
    $workbook = $excel.Workbooks.Open((Resolve-Path $csvFile).Path)
    $worksheet = $workbook.Worksheets.Item(1)
    
    # Set column widths
    $worksheet.Columns.Item('A').ColumnWidth = 12  # NO_PEGAWAI
    $worksheet.Columns.Item('B').ColumnWidth = 20  # NAMA_LENGKAP
    $worksheet.Columns.Item('C').ColumnWidth = 10  # GELAR
    $worksheet.Columns.Item('D').ColumnWidth = 15  # TELEPON_AREA_RUMAH
    $worksheet.Columns.Item('E').ColumnWidth = 15  # TELEPON_RUMAH
    $worksheet.Columns.Item('F').ColumnWidth = 5   # (KOSONG)
    $worksheet.Columns.Item('G').ColumnWidth = 15  # TELEPON_AREA_KANTOR
    $worksheet.Columns.Item('H').ColumnWidth = 15  # TELEPON_KANTOR
    $worksheet.Columns.Item('I').ColumnWidth = 15  # TELEPON_EXT_KANTOR
    $worksheet.Columns.Item('J').ColumnWidth = 15  # HP
    $worksheet.Columns.Item('K').ColumnWidth = 25  # EMAIL
    $worksheet.Columns.Item('L').ColumnWidth = 15  # TEMPAT_LAHIR
    $worksheet.Columns.Item('M').ColumnWidth = 15  # TANGGAL_LAHIR
    $worksheet.Columns.Item('N').ColumnWidth = 20  # NAMA_IBU_KANDUNG
    $worksheet.Columns.Item('O').ColumnWidth = 15  # JENIS_IDENTITAS
    $worksheet.Columns.Item('P').ColumnWidth = 20  # NO_IDENTITAS
    $worksheet.Columns.Item('Q').ColumnWidth = 15  # MASA_LAKU_IDENTITAS
    $worksheet.Columns.Item('R').ColumnWidth = 15  # JENIS_KELAMIN
    $worksheet.Columns.Item('S').ColumnWidth = 25  # SURAT_MENYURAT_KE
    $worksheet.Columns.Item('T').ColumnWidth = 15  # TANGGAL_KEPESERTAAN
    $worksheet.Columns.Item('U').ColumnWidth = 15  # STATUS_KAWIN
    $worksheet.Columns.Item('V').ColumnWidth = 15  # GOLONGAN_DARAH
    $worksheet.Columns.Item('W').ColumnWidth = 20  # NPWP
    $worksheet.Columns.Item('X').ColumnWidth = 10  # KODE_NEGARA
    $worksheet.Columns.Item('Y').ColumnWidth = 15  # UPAH
    $worksheet.Columns.Item('Z').ColumnWidth = 25  # ALAMAT
    $worksheet.Columns.Item('AA').ColumnWidth = 10 # KODE_POS
    $worksheet.Columns.Item('AB').ColumnWidth = 20 # LOKASI_PEKERJAAN
    $worksheet.Columns.Item('AC').ColumnWidth = 15 # STATUS_PEGAWAI
    $worksheet.Columns.Item('AD').ColumnWidth = 15 # TGL_AWAL_BEKERJA
    $worksheet.Columns.Item('AE').ColumnWidth = 15 # TGL_AKHIR_KONTRAK
    
    # Simpan sebagai Excel
    $excelFile = "template_tk_test_valid.xlsx"
    $workbook.SaveAs((Resolve-Path .).Path + "\" + $excelFile, 51) # 51 = xlOpenXMLWorkbook
    $workbook.Close()
    $excel.Quit()
    
    Write-Host "✅ Excel file created: $excelFile" -ForegroundColor Green
    
} catch {
    Write-Host "⚠️  Excel COM object not available, using CSV file" -ForegroundColor Yellow
    Write-Host "   You can open $csvFile in Excel and save as .xlsx" -ForegroundColor Gray
}

Write-Host "`n📋 File Details:" -ForegroundColor Cyan
Write-Host "   - Total data: 3 rows" -ForegroundColor White
Write-Host "   - Columns: 31 (A-AE)" -ForegroundColor White
Write-Host "   - Format: Valid for upload" -ForegroundColor White

Write-Host "`n✅ Test files created successfully!" -ForegroundColor Green
Write-Host "📁 Files:" -ForegroundColor Yellow
Write-Host "   - template_tk_test_valid.csv" -ForegroundColor White
Write-Host "   - template_tk_test_valid.xlsx (if Excel available)" -ForegroundColor White

Write-Host "`n🚀 Next Steps:" -ForegroundColor Green
Write-Host "1. Upload template_tk_test_valid.xlsx" -ForegroundColor White
Write-Host "2. Check if data appears in table" -ForegroundColor White
Write-Host "3. If still invalid, check backend log for errors" -ForegroundColor White







