# PowerShell script untuk validasi format Excel
Write-Host "🔍 Validating Excel Format..." -ForegroundColor Green

Write-Host "`n📋 Format Excel yang Benar untuk Upload TK:" -ForegroundColor Yellow
Write-Host "===============================================" -ForegroundColor Cyan

Write-Host "`n📊 Header Row (Row 1):" -ForegroundColor Yellow
Write-Host "A:  NO_PEGAWAI" -ForegroundColor White
Write-Host "B:  NAMA_LENGKAP" -ForegroundColor Red
Write-Host "C:  GELAR" -ForegroundColor White
Write-Host "D:  TELEPON_AREA_RUMAH" -ForegroundColor White
Write-Host "E:  TELEPON_RUMAH" -ForegroundColor White
Write-Host "F:  (KOSONG)" -ForegroundColor Gray
Write-Host "G:  TELEPON_AREA_KANTOR" -ForegroundColor White
Write-Host "H:  TELEPON_KANTOR" -ForegroundColor White
Write-Host "I:  TELEPON_EXT_KANTOR" -ForegroundColor White
Write-Host "J:  HP" -ForegroundColor White
Write-Host "K:  EMAIL" -ForegroundColor Red
Write-Host "L:  TEMPAT_LAHIR" -ForegroundColor White
Write-Host "M:  TANGGAL_LAHIR" -ForegroundColor White
Write-Host "N:  NAMA_IBU_KANDUNG" -ForegroundColor White
Write-Host "O:  JENIS_IDENTITAS" -ForegroundColor White
Write-Host "P:  NO_IDENTITAS" -ForegroundColor Red
Write-Host "Q:  MASA_LAKU_IDENTITAS" -ForegroundColor White
Write-Host "R:  JENIS_KELAMIN" -ForegroundColor White
Write-Host "S:  SURAT_MENYURAT_KE" -ForegroundColor White
Write-Host "T:  TANGGAL_KEPESERTAAN" -ForegroundColor White
Write-Host "U:  STATUS_KAWIN" -ForegroundColor White
Write-Host "V:  GOLONGAN_DARAH" -ForegroundColor White
Write-Host "W:  NPWP" -ForegroundColor White
Write-Host "X:  KODE_NEGARA" -ForegroundColor White
Write-Host "Y:  UPAH" -ForegroundColor White
Write-Host "Z:  ALAMAT" -ForegroundColor White
Write-Host "AA: KODE_POS" -ForegroundColor White
Write-Host "AB: LOKASI_PEKERJAAN" -ForegroundColor White
Write-Host "AC: STATUS_PEGAWAI" -ForegroundColor White
Write-Host "AD: TGL_AWAL_BEKERJA" -ForegroundColor White
Write-Host "AE: TGL_AKHIR_KONTRAK" -ForegroundColor White

Write-Host "`n⚠️  Kolom WAJIB (Tidak boleh kosong):" -ForegroundColor Red
Write-Host "   B:  NAMA_LENGKAP" -ForegroundColor Red
Write-Host "   K:  EMAIL (format: user@domain.com)" -ForegroundColor Red
Write-Host "   P:  NO_IDENTITAS/NIK (16 digit angka)" -ForegroundColor Red

Write-Host "`n📝 Format Data yang Benar:" -ForegroundColor Yellow
Write-Host "   NIK: 3175051405900001 (16 digit angka)" -ForegroundColor White
Write-Host "   Email: john.doe@example.com" -ForegroundColor White
Write-Host "   NPWP: 012345678901234 (15 digit angka)" -ForegroundColor White
Write-Host "   Golongan Darah: A, B, AB, O" -ForegroundColor White
Write-Host "   Jenis Kelamin: Laki-laki, Perempuan" -ForegroundColor White
Write-Host "   Tanggal: 2024-01-01 (YYYY-MM-DD)" -ForegroundColor White

Write-Host "`n🔍 Contoh Data Row yang Valid:" -ForegroundColor Yellow
Write-Host "A:  PEG001" -ForegroundColor White
Write-Host "B:  John Doe" -ForegroundColor Red
Write-Host "C:  S.Kom" -ForegroundColor White
Write-Host "D:  021" -ForegroundColor White
Write-Host "E:  7654321" -ForegroundColor White
Write-Host "F:  (kosong)" -ForegroundColor Gray
Write-Host "G:  021" -ForegroundColor White
Write-Host "H:  021-1234" -ForegroundColor White
Write-Host "I:  1234" -ForegroundColor White
Write-Host "J:  081234567890" -ForegroundColor White
Write-Host "K:  john.doe@example.com" -ForegroundColor Red
Write-Host "L:  Jakarta" -ForegroundColor White
Write-Host "M:  1990-05-14" -ForegroundColor White
Write-Host "N:  Siti Aminah" -ForegroundColor White
Write-Host "O:  KTP" -ForegroundColor White
Write-Host "P:  3175051405900001" -ForegroundColor Red
Write-Host "Q:  2016-05-14" -ForegroundColor White
Write-Host "R:  Laki-laki" -ForegroundColor White
Write-Host "S:  Jalan Mandaka No. 10" -ForegroundColor White
Write-Host "T:  2020-07-01" -ForegroundColor White
Write-Host "U:  Belum Kawin" -ForegroundColor White
Write-Host "V:  O" -ForegroundColor White
Write-Host "W:  012345678901234" -ForegroundColor White
Write-Host "X:  ID" -ForegroundColor White
Write-Host "Y:  5000000" -ForegroundColor White
Write-Host "Z:  Jl. Merpati 12 Jakarta" -ForegroundColor White
Write-Host "AA: 10130" -ForegroundColor White
Write-Host "AB: Kantor Pusat" -ForegroundColor White
Write-Host "AC: Tetap" -ForegroundColor White
Write-Host "AD: 2020-07-01" -ForegroundColor White
Write-Host "AE: (kosong untuk tetap)" -ForegroundColor Gray

Write-Host "`n❌ Kesalahan Umum yang Menyebabkan Invalid:" -ForegroundColor Red
Write-Host "1. Nama kosong (kolom B)" -ForegroundColor Gray
Write-Host "2. NIK bukan 16 digit (kolom P)" -ForegroundColor Gray
Write-Host "3. Email format salah (kolom K)" -ForegroundColor Gray
Write-Host "4. NPWP bukan 15 digit (kolom W)" -ForegroundColor Gray
Write-Host "5. Golongan darah bukan A/B/AB/O (kolom V)" -ForegroundColor Gray
Write-Host "6. Jenis kelamin bukan Laki-laki/Perempuan (kolom R)" -ForegroundColor Gray
Write-Host "7. NIK duplikat (sudah ada di database)" -ForegroundColor Gray
Write-Host "8. Format tanggal salah (YYYY-MM-DD)" -ForegroundColor Gray

Write-Host "`n🛠️ Cara Memperbaiki:" -ForegroundColor Green
Write-Host "1. Pastikan ada 31 kolom (A-AE)" -ForegroundColor White
Write-Host "2. Isi kolom wajib: Nama, Email, NIK" -ForegroundColor White
Write-Host "3. Format NIK: 16 digit angka" -ForegroundColor White
Write-Host "4. Format Email: user@domain.com" -ForegroundColor White
Write-Host "5. Cek NIK tidak duplikat" -ForegroundColor White
Write-Host "6. Test dengan 1-2 data dulu" -ForegroundColor White

Write-Host "`n✅ Validasi selesai!" -ForegroundColor Green
Write-Host "Pastikan format Excel sesuai template di atas" -ForegroundColor Cyan







