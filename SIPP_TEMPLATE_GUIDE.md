# 🎯 **File Excel SIPP Valid Sudah Siap!**

## 📁 **File yang Tersedia:**

### **1. `template_sipp_valid.xlsx`** ✅ **RECOMMENDED**
- **Format:** Excel (.xlsx) yang benar-benar valid
- **Data:** 3 baris data dengan format SIPP lengkap
- **Kolom:** 33 kolom sesuai template SIPP BPJamsostek
- **Status:** Siap untuk upload massal
- **Size:** 6.3 KB (format Excel yang proper)

### **2. `template_tk_test_valid.xlsx`** ⚠️ **OLD VERSION**
- **Format:** File CSV yang disimpan sebagai .xlsx
- **Status:** Tidak valid (menyebabkan error "File Excel tidak valid atau rusak")

## 📊 **Data Test yang Tersedia:**

| No | Nama | NIK | Email | Status Pegawai | Upah | Status |
|----|------|-----|-------|----------------|------|--------|
| 1 | **John Doe** | 3175051405900001 | john.doe@example.com | PKWTT | 5.000.000 | ✅ Valid |
| 2 | **Jane Smith** | 3175051405900002 | jane.smith@example.com | PKWTT | 6.000.000 | ✅ Valid |
| 3 | **Ahmad Rahman** | 3175051405900003 | ahmad.rahman@example.com | PKWT | 4.500.000 | ✅ Valid |

## 🔍 **Format Data yang Benar (Sesuai SIPP):**

### **Kolom Wajib (Tidak boleh kosong):**
- **A: NO_PEGAWAI** - Kode pegawai (PEG001, PEG002, dll)
- **B: NAMA_LENGKAP** - Nama lengkap sesuai KTP
- **K: EMAIL** - Format: user@domain.com
- **P: NO_IDENTITAS** - NIK 16 digit angka
- **Y: UPAH** - Gaji dalam angka (5000000)

### **Format Data Valid:**
- **NIK:** 16 digit angka (3175051405900001)
- **Email:** Format email valid (john.doe@example.com)
- **NPWP:** 15 digit angka (012345678901234)
- **Golongan Darah:** A, B, AB, O
- **Jenis Kelamin:** Laki-laki, Perempuan
- **Status Pegawai:** PKWTT (Tetap), PKWT (Kontrak)
- **Tanggal:** YYYY-MM-DD (1990-05-14)

### **Mapping Kolom Lengkap (A-AG):**
```
A:  NO_PEGAWAI           - Kode pegawai
B:  NAMA_LENGKAP         - Nama lengkap
C:  GELAR                - Gelar pendidikan
D:  TELEPON_AREA_RUMAH   - Kode area telepon rumah
E:  TELEPON_RUMAH        - Nomor telepon rumah
F:  (KOSONG)             - Kolom kosong sesuai template
G:  TELEPON_AREA_KANTOR  - Kode area telepon kantor
H:  TELEPON_KANTOR       - Nomor telepon kantor
I:  TELEPON_EXT_KANTOR  - Ekstensi telepon kantor
J:  HP                   - Nomor handphone
K:  EMAIL                - Alamat email
L:  TEMPAT_LAHIR         - Tempat lahir
M:  TANGGAL_LAHIR        - Tanggal lahir
N:  NAMA_IBU_KANDUNG     - Nama ibu kandung
O:  JENIS_IDENTITAS      - Jenis identitas (KTP)
P:  NO_IDENTITAS         - Nomor identitas (NIK)
Q:  MASA_LAKU_IDENTITAS  - Masa berlaku identitas
R:  JENIS_KELAMIN        - Jenis kelamin
S:  SURAT_MENYURAT_KE    - Alamat surat menyurat
T:  TANGGAL_KEPESERTAAN  - Tanggal kepesertaan
U:  STATUS_KAWIN         - Status kawin
V:  GOLONGAN_DARAH       - Golongan darah
W:  NPWP                 - Nomor NPWP
X:  KODE_NEGARA          - Kode negara
Y:  UPAH                 - Upah/gaji
Z:  ALAMAT               - Alamat lengkap
AA: KODE_POS             - Kode pos
AB: LOKASI_PEKERJAAN     - Lokasi pekerjaan
AC: STATUS_PEGAWAI       - Status pegawai (PKWTT/PKWT)
AD: TGL_AWAL_BEKERJA     - Tanggal awal bekerja
AE: TGL_AKHIR_KONTRAK    - Tanggal akhir kontrak (kosong untuk PKWTT)
AF: RAPEL                - Rapel
AG: NATIONALITY          - Kewarganegaraan
```

## 🚀 **Cara Test Upload:**

### **Step 1: Gunakan File yang Benar**
```
✅ GUNAKAN: template_sipp_valid.xlsx
❌ JANGAN: template_tk_test_valid.xlsx (tidak valid)
```

### **Step 2: Upload via Frontend**
1. **Buka aplikasi** e-BPJSTK
2. **Login** dengan akun yang valid
3. **Pilih menu** "Upload TK" atau "Tambah Massal"
4. **Pilih file** `template_sipp_valid.xlsx`
5. **Klik "Upload"** atau "Submit"

### **Step 3: Expected Results**

#### **Jika Upload Berhasil:**
- **Total Data:** 3
- **Valid:** 3
- **Invalid:** 0
- **Data muncul di tabel** dengan 3 baris baru
- **Pop-up:** "File berhasil diupload!"

#### **Jika Upload Gagal:**
- **Total Data:** 3
- **Valid:** 0
- **Invalid:** 3
- **Cek log backend** untuk error detail

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

### **3. Kemungkinan Masalah:**

#### **A. Database Schema**
- ❌ **Field baru belum ada** di database
- ❌ **Migration belum dijalankan**

**Solusi:**
```bash
# Jalankan migration
.\run_migration.ps1
```

#### **B. Backend Logic**
- ❌ **Validasi terlalu ketat**
- ❌ **Mapping kolom salah**

**Solusi:**
```bash
# Restart backend dengan patch baru
cd backend
go run main.go
```

#### **C. Format Data**
- ❌ **NIK duplikat** (sudah ada di database)
- ❌ **Email format salah**
- ❌ **Upah < UMR**

**Solusi:**
- Cek NIK tidak duplikat
- Pastikan email format benar
- Sesuaikan upah dengan UMR

## 🛠️ **Troubleshooting Lengkap:**

### **1. File Excel Tidak Valid**
```
Error: "File Excel tidak valid atau rusak"
Solusi: Gunakan template_sipp_valid.xlsx (bukan template_tk_test_valid.xlsx)
```

### **2. Semua Data Invalid**
```
Error: Valid: 0, Invalid: 3
Solusi: 
1. Cek log backend untuk error detail
2. Pastikan database migration sudah dijalankan
3. Cek NIK tidak duplikat
4. Pastikan format data sesuai template
```

### **3. Data Tidak Muncul di Tabel**
```
Error: Upload berhasil tapi data tidak tampil
Solusi:
1. Refresh halaman (F5)
2. Cek filter/search aktif
3. Cek pagination
4. Cek user ID mismatch
```

## 📋 **Checklist Upload Massal:**

### **Sebelum Upload:**
- [ ] File Excel format .xlsx (bukan CSV)
- [ ] Ada 33 kolom (A-AG)
- [ ] Header sesuai template SIPP
- [ ] Data minimal 1 baris (selain header)
- [ ] NIK 16 digit angka
- [ ] Email format benar
- [ ] Upah >= UMR

### **Saat Upload:**
- [ ] Pilih file yang benar
- [ ] Tunggu proses upload selesai
- [ ] Cek hasil validasi
- [ ] Lihat pop-up hasil

### **Setelah Upload:**
- [ ] Data muncul di tabel
- [ ] Jumlah TK bertambah
- [ ] Upload history tercatat
- [ ] Tidak ada error di log

## ✅ **File Test Siap Digunakan!**

**Gunakan file `template_sipp_valid.xlsx` untuk test upload massal. File ini sudah dibuat sesuai spesifikasi SIPP BPJamsostek yang lengkap dan seharusnya bisa diupload tanpa error.**

**Jika masih ada masalah, cek log backend untuk error detail yang spesifik.** 🚀







