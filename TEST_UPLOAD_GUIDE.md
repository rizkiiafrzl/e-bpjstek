# 🧪 **Panduan Test Upload Massal TK**

## 📁 **File Test yang Tersedia:**

### **1. `template_tk_test_valid.xlsx`** ✅
- **Format:** Excel (.xlsx)
- **Data:** 3 baris data valid
- **Kolom:** 31 kolom (A-AE) sesuai template SIPP
- **Status:** Siap untuk upload

### **2. `template_tk_test_valid.csv`** ✅
- **Format:** CSV
- **Data:** 3 baris data valid
- **Kolom:** 31 kolom sesuai template SIPP
- **Status:** Bisa dikonversi ke Excel

## 📊 **Data Test yang Tersedia:**

| No | Nama | NIK | Email | Status |
|----|------|-----|-------|--------|
| 1 | John Doe | 3175051405900001 | john.doe@example.com | Valid |
| 2 | Jane Smith | 3175051405900002 | jane.smith@example.com | Valid |
| 3 | Ahmad Rahman | 3175051405900003 | ahmad.rahman@example.com | Valid |

## 🔍 **Format Data yang Benar:**

### **Kolom Wajib (Tidak boleh kosong):**
- **B: NAMA_LENGKAP** - Nama lengkap
- **K: EMAIL** - Format: user@domain.com
- **P: NO_IDENTITAS/NIK** - 16 digit angka

### **Format Data Valid:**
- **NIK:** 16 digit angka (3175051405900001)
- **Email:** Format email valid (john.doe@example.com)
- **NPWP:** 15 digit angka (012345678901234)
- **Golongan Darah:** A, B, AB, O
- **Jenis Kelamin:** Laki-laki, Perempuan
- **Tanggal:** YYYY-MM-DD (1990-05-14)

## 🚀 **Cara Test Upload:**

### **Step 1: Siapkan File**
```bash
# File sudah tersedia di root directory
template_tk_test_valid.xlsx
```

### **Step 2: Upload via Frontend**
1. Buka aplikasi e-BPJSTK
2. Login dengan akun yang valid
3. Pilih menu "Upload TK" atau "Tambah Massal"
4. Pilih file `template_tk_test_valid.xlsx`
5. Klik "Upload" atau "Submit"

### **Step 3: Cek Hasil**
1. **Jika berhasil:** Data akan muncul di tabel
2. **Jika gagal:** Cek popup error dan log backend

## 🔍 **Debug jika Masih Invalid:**

### **1. Cek Log Backend**
```bash
# Lihat log backend untuk error detail
tail -f backend.log | grep -i "upload\|error\|invalid"
```

### **2. Cek Upload History**
```powershell
# Jalankan script debug
.\check_upload_errors.ps1
```

### **3. Cek Database**
```sql
-- Cek data terbaru
SELECT * FROM workers ORDER BY created_at DESC LIMIT 5;

-- Cek upload history
SELECT * FROM upload_histories ORDER BY created_at DESC LIMIT 3;
```

## ⚠️ **Kemungkinan Masalah:**

### **1. Format Excel Salah**
- ❌ **Tidak ada 31 kolom** (A-AE)
- ❌ **Header row tidak sesuai** template
- ❌ **Data dimulai dari row yang salah**

### **2. Data Kosong di Kolom Wajib**
- ❌ **Nama kosong** (kolom B)
- ❌ **NIK kosong/tidak 16 digit** (kolom P)
- ❌ **Email kosong/format salah** (kolom K)

### **3. Format Data Salah**
- ❌ **NIK bukan 16 digit angka**
- ❌ **Email format salah**
- ❌ **NPWP bukan 15 digit**
- ❌ **Golongan darah bukan A/B/AB/O**

### **4. NIK Duplikat**
- ❌ **NIK sudah ada** di database

## 🛠️ **Solusi Cepat:**

### **1. Restart Backend**
```bash
cd backend
go run main.go
```

### **2. Jalankan Migration**
```bash
# Windows
.\run_migration.ps1

# Linux/Mac
./run_migration.sh
```

### **3. Test dengan Data Minimal**
- Gunakan file test yang sudah disediakan
- Cek log backend untuk error detail
- Pastikan format Excel sesuai template

## 📋 **Template Excel yang Benar:**

### **Header Row (Row 1):**
```
A: NO_PEGAWAI, B: NAMA_LENGKAP, C: GELAR, D: TELEPON_AREA_RUMAH, E: TELEPON_RUMAH, F: (KOSONG), G: TELEPON_AREA_KANTOR, H: TELEPON_KANTOR, I: TELEPON_EXT_KANTOR, J: HP, K: EMAIL, L: TEMPAT_LAHIR, M: TANGGAL_LAHIR, N: NAMA_IBU_KANDUNG, O: JENIS_IDENTITAS, P: NO_IDENTITAS, Q: MASA_LAKU_IDENTITAS, R: JENIS_KELAMIN, S: SURAT_MENYURAT_KE, T: TANGGAL_KEPESERTAAN, U: STATUS_KAWIN, V: GOLONGAN_DARAH, W: NPWP, X: KODE_NEGARA, Y: UPAH, Z: ALAMAT, AA: KODE_POS, AB: LOKASI_PEKERJAAN, AC: STATUS_PEGAWAI, AD: TGL_AWAL_BEKERJA, AE: TGL_AKHIR_KONTRAK
```

### **Data Row (Row 2+):**
```
PEG001,John Doe,S.Kom,021,7654321,,021,021-1234,1234,081234567890,john.doe@example.com,Jakarta,1990-05-14,Siti Aminah,KTP,3175051405900001,2016-05-14,Laki-laki,Jalan Mandaka No. 10,2020-07-01,Belum Kawin,O,012345678901234,ID,5000000,Jl. Merpati 12 Jakarta,10130,Kantor Pusat,Tetap,2020-07-01,
```

## ✅ **Expected Results:**

### **Jika Upload Berhasil:**
- **Total Data:** 3
- **Valid:** 3
- **Invalid:** 0
- **Data muncul di tabel** dengan 3 baris baru

### **Jika Upload Gagal:**
- **Total Data:** 3
- **Valid:** 0
- **Invalid:** 3
- **Cek log backend** untuk error detail

## 🎯 **Next Steps:**

1. **Upload file test** `template_tk_test_valid.xlsx`
2. **Cek hasil** di tabel frontend
3. **Jika berhasil:** Upload massal sudah berfungsi
4. **Jika gagal:** Cek log backend dan debug error

**File test sudah siap untuk digunakan!** 🚀







