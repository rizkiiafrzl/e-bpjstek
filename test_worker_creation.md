# Test Pembuatan Tenaga Kerja Baru

## 🧪 **Test Cases untuk Validasi Data**

### **1. Test Data Valid (WNI)**
```json
{
  "nik": "3175051405900001",
  "kpj": "12345678901234567890",
  "noPegawai": "PEG001",
  "nama": "John Doe",
  "dateOfBirth": "1990-05-14",
  "upah": 5000000,
  "rapel": 0,
  "nationality": "WNI",
  "tempatLahir": "Jakarta",
  "ibuKandung": "Siti Aminah",
  "jenisKelamin": "Laki-laki",
  "golDarah": "O",
  "statusKawin": "Belum Kawin",
  "statusPegawai": "Tetap",
  "tanggalAwalBekerja": "2020-07-01",
  "tanggalAkhirKontrak": "",
  "lokasiPekerjaan": "Kantor Pusat",
  "teleponAreaRumah": "021",
  "teleponRumah": "7654321",
  "teleponAreaKantor": "021",
  "teleponKantor": "021-1234",
  "teleponExtKantor": "1234",
  "handphone": "081234567890",
  "email": "john.doe@example.com",
  "alamat": "Jl. Merpati 12 Jakarta",
  "kabupaten": "Jakarta Selatan",
  "kodePos": "10130",
  "npwp": "012345678901234",
  "jenisIdentitas": "KTP",
  "masaLakuIdentitas": "2016-05-14",
  "suratMenyuratKe": "Jalan Mandaka No. 10",
  "tanggalKepesertaan": "2020-07-01",
  "kodeNegara": "ID"
}
```

### **2. Test Data Valid (WNA)**
```json
{
  "nik": "",
  "kpj": "",
  "noPegawai": "PEG002",
  "nama": "Jane Smith",
  "dateOfBirth": "1985-03-20",
  "upah": 7500000,
  "rapel": 0,
  "nationality": "WNA",
  "passportNo": "A1234567",
  "passportValidUntil": "2025-03-20",
  "tempatLahir": "New York",
  "ibuKandung": "Mary Smith",
  "jenisKelamin": "Perempuan",
  "golDarah": "A",
  "statusKawin": "Kawin",
  "statusPegawai": "Kontrak",
  "tanggalAwalBekerja": "2021-01-15",
  "tanggalAkhirKontrak": "2024-01-15",
  "lokasiPekerjaan": "Cabang Bandung",
  "teleponAreaRumah": "022",
  "teleponRumah": "3344556",
  "teleponAreaKantor": "022",
  "teleponKantor": "022-2345",
  "teleponExtKantor": "2345",
  "handphone": "081216789012",
  "email": "jane.smith@example.com",
  "alamat": "Jl. Anggrek 8 Bandung",
  "kabupaten": "Bandung",
  "kodePos": "40212",
  "npwp": "987654321098765",
  "jenisIdentitas": "Passport",
  "masaLakuIdentitas": "2025-03-20",
  "suratMenyuratKe": "Kepala HR",
  "tanggalKepesertaan": "2021-01-15",
  "kodeNegara": "US"
}
```

### **3. Test Data Invalid - NIK Salah Format**
```json
{
  "nik": "317505140590001",
  "nama": "Invalid NIK",
  "nationality": "WNI"
}
```
**Expected Error:** "Format NIK tidak valid (harus 16 digit)"

### **4. Test Data Invalid - Email Salah Format**
```json
{
  "nik": "3175051405900001",
  "nama": "Invalid Email",
  "email": "invalid-email-format",
  "nationality": "WNI"
}
```
**Expected Error:** "Format email tidak valid"

### **5. Test Data Invalid - NPWP Salah Format**
```json
{
  "nik": "3175051405900001",
  "nama": "Invalid NPWP",
  "npwp": "12345678901234",
  "nationality": "WNI"
}
```
**Expected Error:** "Format NPWP tidak valid (harus 15 digit)"

### **6. Test Data Invalid - Golongan Darah Salah**
```json
{
  "nik": "3175051405900001",
  "nama": "Invalid Golongan Darah",
  "golDarah": "X",
  "nationality": "WNI"
}
```
**Expected Error:** "Golongan darah tidak valid (A, B, AB, O)"

### **7. Test Data Invalid - WNA Tanpa Passport**
```json
{
  "nama": "WNA Without Passport",
  "nationality": "WNA"
}
```
**Expected Error:** "Nomor paspor wajib untuk WNA"

## 🚀 **Cara Test**

### **1. Test via API Endpoint**
```bash
# Test create worker
curl -X POST http://localhost:8080/api/v1/workers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d @test_data_valid_wni.json
```

### **2. Test via Frontend**
1. Login ke aplikasi
2. Buka halaman "Tambah Tenaga Kerja"
3. Isi form dengan data test di atas
4. Submit dan lihat hasilnya

### **3. Test Upload Excel**
1. Download template Excel
2. Isi dengan data test
3. Upload file melalui halaman "Upload TK"
4. Lihat hasil validasi

## ✅ **Expected Results**

### **Valid Data:**
- Status: 201 Created
- Response: Data worker yang baru dibuat
- Database: Record tersimpan dengan semua field

### **Invalid Data:**
- Status: 400 Bad Request
- Response: Error message yang spesifik
- Database: Tidak ada record baru

## 🔍 **Field Validation Summary**

| Field | Validation | Error Message |
|-------|------------|---------------|
| NIK | 16 digit angka | "Format NIK tidak valid (harus 16 digit)" |
| Email | Format email valid | "Format email tidak valid" |
| NPWP | 15 digit angka | "Format NPWP tidak valid (harus 15 digit)" |
| GolDarah | A, B, AB, O | "Golongan darah tidak valid (A, B, AB, O)" |
| JenisKelamin | L, P, Laki-laki, Perempuan | "Jenis kelamin tidak valid" |
| WNA | Wajib PassportNo | "Nomor paspor wajib untuk WNA" |
| Nama | Tidak boleh kosong | "Nama wajib diisi" |

## 📊 **Database Schema Update**

Setelah update, tabel `workers` akan memiliki field tambahan:
- `telepon_area_rumah`
- `telepon_area_kantor` 
- `telepon_kantor`
- `telepon_ext_kantor`
- `jenis_identitas`
- `masa_laku_identitas`
- `surat_menyurat_ke`
- `tanggal_kepesertaan`
- `kode_negara`







