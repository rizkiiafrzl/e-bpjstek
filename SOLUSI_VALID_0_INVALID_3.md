# 🚨 **SOLUSI MASALAH: Valid: 0, Invalid: 3**

## 📋 **Masalah yang Terjadi**

Dari screenshot yang Anda tunjukkan, upload file Excel berhasil dilakukan namun hasilnya:
- **Total Data:** 3
- **Valid:** 0  
- **Invalid:** 3

Artinya semua 3 data dianggap tidak valid oleh sistem. Ini menunjukkan ada masalah dengan validasi backend.

---

## 🔍 **Kemungkinan Penyebab**

### **1. Validasi Terlalu Ketat**
Backend menggunakan validasi SIPP yang ketat, kemungkinan ada field yang tidak sesuai:

- **NIK:** Harus tepat 16 digit angka
- **Email:** Format email yang valid
- **NPWP:** Harus tepat 15 digit angka
- **Upah:** Minimal Rp 1.000.000
- **Status Pegawai:** PKWT memerlukan tanggal akhir kontrak
- **Jenis Kelamin:** Harus "Laki-laki" atau "Perempuan"
- **Golongan Darah:** Harus A, B, AB, atau O

### **2. Database Schema**
Field baru mungkin belum ada di database:
- `telepon_area_rumah`
- `telepon_area_kantor`
- `telepon_kantor`
- `telepon_ext_kantor`
- `jenis_identitas`
- `masa_laku_identitas`
- `surat_menyurat_ke`
- `tanggal_kepesertaan`
- `kode_negara`

### **3. Mapping Kolom Excel**
Mapping dari Excel ke database mungkin tidak sesuai dengan template SIPP.

---

## 🛠️ **Langkah-langkah Debugging**

### **Step 1: Cek Log Backend**
```bash
# Lihat console backend untuk error detail
# Backend akan menampilkan log seperti:
# "Validation failed for row 2: [NIK harus 16 digit angka]"
```

### **Step 2: Jalankan Script Debug**
```powershell
# Jalankan script debug yang sudah dibuat
.\debug_upload_sipp.ps1
```

**Sebelum menjalankan, ganti `YOUR_JWT_TOKEN_HERE` dengan token JWT yang valid.**

### **Step 3: Cek Database Schema**
```sql
-- Jalankan query ini untuk cek field yang ada
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'workers' 
ORDER BY column_name;
```

### **Step 4: Jalankan Migration**
```bash
# Jika field baru belum ada, jalankan migration
.\run_migration.ps1
```

---

## 🔧 **Solusi Cepat**

### **A. Perbaiki Validasi (Temporary)**

Jika ingin mengurangi validasi sementara untuk testing, edit `backend/handlers/worker.go`:

```go
// Di fungsi ValidasiSIPP, comment beberapa validasi yang ketat
func ValidasiSIPP(req *CreateWorkerRequest) []string {
    var errors []string

    // Hanya validasi yang paling penting
    if req.Nama == "" {
        errors = append(errors, "Nama lengkap wajib diisi")
    }
    if req.NIK == "" {
        errors = append(errors, "NIK wajib diisi")
    }
    if req.Email == "" {
        errors = append(errors, "Email wajib diisi")
    }
    
    // Comment validasi yang ketat untuk testing
    /*
    if req.Upah < 1000000 {
        errors = append(errors, "Upah minimal Rp 1.000.000")
    }
    */
    
    return errors
}
```

### **B. Perbaiki Data Excel**

Pastikan data di Excel sesuai format:

| Kolom | Field | Format | Contoh |
|-------|-------|--------|--------|
| A | NO_PEGAWAI | Text | PEG001 |
| B | NAMA_LENGKAP | Text | John Doe |
| P | NO_IDENTITAS | 16 digit | 3175051405900001 |
| K | EMAIL | Email format | john.doe@example.com |
| Y | UPAH | Number | 5000000 |
| AC | STATUS_PEGAWAI | PKWTT/PKWT | PKWTT |
| R | JENIS_KELAMIN | Laki-laki/Perempuan | Laki-laki |
| V | GOLONGAN_DARAH | A/B/AB/O | O |

### **C. Cek Mapping Kolom**

Pastikan mapping di backend sesuai dengan template Excel:

```go
// Di fungsi UploadTK, pastikan mapping benar
req := CreateWorkerRequest{
    NoPegawai:           get(0),   // A: NO_PEGAWAI
    Nama:                get(1),   // B: NAMA_LENGKAP
    NIK:                 get(15),  // P: NO_IDENTITAS
    Email:               get(10),  // K: EMAIL
    Upah:                parseUpah(get(24)), // Y: UPAH
    StatusPegawai:       get(28),  // AC: STATUS_PEGAWAI
    JenisKelamin:        get(17),  // R: JENIS_KELAMIN
    GolDarah:            get(21),  // V: GOLONGAN_DARAH
    // ... dst
}
```

---

## 🚀 **Langkah-langkah Perbaikan**

### **1. Restart Backend dengan Logging**
```bash
cd backend
go run main.go
```

### **2. Upload File Excel**
Gunakan file `template_sipp_valid.xlsx` yang sudah dibuat.

### **3. Lihat Log Backend**
Backend akan menampilkan error detail seperti:
```
Validation failed for row 2: [NIK harus 16 digit angka, Email wajib diisi]
Validation failed for row 3: [Tanggal akhir kontrak wajib diisi untuk PKWT]
```

### **4. Perbaiki Data Excel**
Berdasarkan error yang muncul, perbaiki data di Excel.

### **5. Upload Ulang**
Upload file yang sudah diperbaiki.

---

## 📊 **Expected Results Setelah Perbaikan**

### **Jika Berhasil:**
```
✅ Upload berhasil!
📊 Total Data: 3
✅ Valid: 3
❌ Invalid: 0
```

### **Jika Masih Ada Error:**
```
❌ Upload selesai dengan error
📊 Total Data: 3
✅ Valid: 2
❌ Invalid: 1
🔍 Error Details:
   - Row 3: NIK harus 16 digit angka
```

---

## 🔍 **Script Debug yang Tersedia**

### **1. `debug_upload_sipp.ps1`**
- Test individual worker creation
- Test mass upload
- Tampilkan error detail
- Cek upload history

### **2. `check_backend_logs.ps1`**
- Cek status backend
- Lihat log files
- Cek database connection
- Cek upload history

### **3. `debug_upload_sipp.py`**
- Script Python untuk debug
- Test API endpoints
- Validasi data

---

## 💡 **Tips Debugging**

### **1. Gunakan Token JWT yang Valid**
```javascript
// Cara mendapatkan token:
// 1. Login ke aplikasi
// 2. Buka Developer Tools (F12)
// 3. Lihat Network tab
// 4. Cari request login
// 5. Copy token dari header Authorization
```

### **2. Cek Backend Console**
Backend akan menampilkan log detail untuk setiap validasi yang gagal.

### **3. Gunakan File Excel yang Benar**
Pastikan menggunakan `template_sipp_valid.xlsx` yang sudah dibuat dengan format yang benar.

### **4. Cek Database Schema**
Pastikan semua field baru sudah ada di database.

---

## 🎯 **Quick Fix**

Jika ingin solusi cepat untuk testing:

1. **Comment validasi ketat** di `ValidasiSIPP`
2. **Restart backend**
3. **Upload file Excel**
4. **Lihat hasil**

Setelah berhasil, kembalikan validasi untuk production.

---

## 📞 **Jika Masih Bermasalah**

1. **Jalankan script debug** dengan token yang valid
2. **Cek log backend** untuk error detail
3. **Verifikasi data Excel** sesuai template SIPP
4. **Pastikan database migration** sudah dijalankan
5. **Cek mapping kolom** Excel ke database

**File Excel `template_sipp_valid.xlsx` sudah dibuat dengan data yang benar, jadi masalahnya kemungkinan di validasi backend atau database schema.** 🚀







