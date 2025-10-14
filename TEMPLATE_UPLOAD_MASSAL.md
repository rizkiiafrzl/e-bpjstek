# 📋 Template Upload Massal - Panduan Lengkap

## 🎯 **Template Upload TK (Tenaga Kerja) - 31 Kolom**

Template ini mendukung upload massal dengan **semua field** yang sama seperti form individu.

### **Format Template Excel:**

| Kolom | Field | Format | Contoh | Keterangan |
|-------|-------|--------|--------|------------|
| A | NO_PEGAWAI | Text | PEG001 | Nomor pegawai |
| B | NAMA_LENGKAP | Text | John Doe | **WAJIB** - Nama lengkap |
| C | GELAR | Text | S.Kom | Gelar akademik |
| D | TELEPON_AREA_RUMAH | Text | 021 | Kode area telepon rumah |
| E | TELEPON_RUMAH | Text | 7654321 | Nomor telepon rumah |
| F | EMAIL | Text | john@example.com | Email (kolom kosong) |
| G | TELEPON_AREA_KANTOR | Text | 021 | Kode area telepon kantor |
| H | TELEPON_KANTOR | Text | 021-1234 | Nomor telepon kantor |
| I | TELEPON_EXT_KANTOR | Text | 1234 | Ekstensi telepon kantor |
| J | HP | Text | 081234567890 | Nomor handphone |
| K | EMAIL | Text | john@example.com | **Email sebenarnya** |
| L | TEMPAT_LAHIR | Text | Jakarta | Tempat lahir |
| M | TANGGAL_LAHIR | YYYY-MM-DD | 1990-05-14 | **PENTING: Format tanggal** |
| N | NAMA_IBU_KANDUNG | Text | Siti Aminah | Nama ibu kandung |
| O | JENIS_IDENTITAS | Text | KTP | Jenis identitas |
| P | NO_IDENTITAS | 16 digit | 3175051405900001 | **WAJIB** - NIK 16 digit |
| Q | MASA_LAKU_IDENTITAS | YYYY-MM-DD | 2016-05-14 | Masa berlaku identitas |
| R | JENIS_KELAMIN | L/P | Laki-laki | L=Laki, P=Perempuan |
| S | SURAT_MENYURAT_KE | Text | Jalan Mandaka No. 10 | Alamat surat menyurat |
| T | TANGGAL_KEPESERTAAN | YYYY-MM-DD | 2020-07-01 | Tanggal kepesertaan |
| U | STATUS_KAWIN | Text | Belum Kawin | Status perkawinan |
| V | GOLONGAN_DARAH | A/B/AB/O | O | Golongan darah |
| W | NPWP | Text | 012345678901234 | Nomor NPWP |
| X | KODE_NEGARA | Text | ID | Kode negara |
| Y | UPAH | Numeric | 5000000.00 | **PENTING: Titik (.) bukan koma** |
| Z | ALAMAT | Text | Jl. Merpati 12 Jakarta | Alamat lengkap |
| AA | KODE_POS | 5 digit | 10130 | Kode pos |
| AB | LOKASI_PEKERJAAN | Text | Kantor Pusat | Lokasi pekerjaan |
| AC | STATUS_PEGAWAI | Text | Tetap | Status pegawai |
| AD | TGL_AWAL_BEKERJA | YYYY-MM-DD | 2020-07-01 | Tanggal mulai bekerja |
| AE | TGL_AKHIR_KONTRAK | YYYY-MM-DD | 2023-12-15 | Tanggal akhir kontrak (kosong untuk tetap) |

## ⚠️ **PENTING - Format Data:**

### **1. Format Angka:**
- ✅ **Benar**: `5000000.00`
- ❌ **Salah**: `5.000.000,00`

### **2. Format Tanggal:**
- ✅ **Benar**: `2020-07-01`
- ❌ **Salah**: `01/07/2020` atau `1 Juli 2020`

### **3. Format NIK:**
- ✅ **Benar**: `3175051405900001` (16 digit)
- ❌ **Salah**: `317505140590000` (kurang digit)

## 🚨 **Validasi Data:**

1. **NIK harus 16 digit dan unik**
2. **Nama tidak boleh kosong**
3. **Format tanggal: YYYY-MM-DD**
4. **Format angka: gunakan titik (.) untuk desimal**
5. **Email harus valid (jika diisi)**

## 📝 **Contoh Data yang Benar:**

```
PEG001,John Doe,S.Kom,021,7654321,,021,021-1234,1234,081234567890,john.doe@example.com,Jakarta,1990-05-14,Siti Aminah,KTP,3175051405900001,2016-05-14,Laki-laki,Jalan Mandaka No. 10,2020-07-01,Belum Kawin,O,012345678901234,ID,5000000.00,Jl. Merpati 12 Jakarta,10130,Kantor Pusat,Tetap,2020-07-01,
```

## 🔧 **Langkah Upload:**

1. **Download template** dari sistem
2. **Isi data** sesuai format di atas
3. **Validasi** sebelum upload
4. **Upload file** melalui sistem
5. **Cek hasil** upload

## ✅ **Keuntungan Template Lengkap:**

- ✅ **Data lengkap** seperti form individu
- ✅ **Upload massal** efisien
- ✅ **Validasi otomatis** di backend
- ✅ **Mapping field** yang tepat
- ✅ **Format data** yang konsisten

## 🆘 **Troubleshooting:**

### **Data tidak masuk:**
- Cek format NIK (16 digit)
- Cek format tanggal (YYYY-MM-DD)
- Cek format angka (titik desimal)

### **Data duplikat:**
- NIK sudah ada di sistem
- Gunakan NIK yang unik

### **Format salah:**
- Pastikan format sesuai panduan
- Gunakan template yang benar
