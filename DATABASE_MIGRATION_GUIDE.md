# 🗄️ Database Migration Guide - e-BPJSTK

## 📋 **Overview**
Script ini akan menambahkan field yang hilang di tabel `workers` agar sesuai dengan persyaratan SIPP.

## 🔧 **Field yang Akan Ditambahkan**

### **Contact Info Fields:**
- `telepon_area_rumah` VARCHAR(8)
- `telepon_area_kantor` VARCHAR(8) 
- `telepon_kantor` VARCHAR(32)
- `telepon_ext_kantor` VARCHAR(16)

### **Additional Info Fields:**
- `jenis_identitas` VARCHAR(16)
- `masa_laku_identitas` VARCHAR(16)
- `surat_menyurat_ke` VARCHAR(200)
- `tanggal_kepesertaan` VARCHAR(16)
- `kode_negara` VARCHAR(4)

## 🚀 **Cara Menjalankan Migration**

### **Opsi 1: Menggunakan Go Migration Script**
```bash
# Linux/Mac
chmod +x run_migration.sh
./run_migration.sh

# Windows PowerShell
.\run_migration.ps1
```

### **Opsi 2: Manual SQL Migration**
```sql
-- Jalankan di psql atau database client
\i database_migration.sql
```

### **Opsi 3: Menggunakan GORM AutoMigrate**
```bash
cd backend
go run main.go
# GORM akan otomatis menambahkan field yang hilang
```

## 🔍 **Verifikasi Migration**

### **1. Cek Struktur Tabel**
```sql
-- Cek semua field di tabel workers
\d workers;

-- Atau menggunakan query
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'workers' 
ORDER BY ordinal_position;
```

### **2. Cek Field SIPP**
```sql
-- Jalankan script verifikasi
\i check_database.sql
```

### **3. Test Insert Data**
```sql
-- Test insert data dengan field baru
INSERT INTO workers (
    no_pegawai, nama, nik, telepon_area_rumah, telepon_kantor,
    jenis_identitas, tanggal_kepesertaan, kode_negara
) VALUES (
    'TEST001', 'Test Worker', '3175051405900001', '021', '021-1234',
    'KTP', '2020-07-01', 'ID'
);
```

## 📊 **Expected Results**

### **Sebelum Migration:**
```
Field yang ada: 25 field
Field SIPP: 16 field (kurang 9 field)
```

### **Setelah Migration:**
```
Field yang ada: 34 field
Field SIPP: 25 field (lengkap)
```

## ⚠️ **Troubleshooting**

### **Error: Column already exists**
```sql
-- Field sudah ada, skip migration
-- Cek apakah field sudah benar
SELECT column_name FROM information_schema.columns 
WHERE table_name = 'workers' AND column_name = 'telepon_area_rumah';
```

### **Error: Permission denied**
```sql
-- Pastikan user memiliki ALTER TABLE permission
GRANT ALTER ON TABLE workers TO your_user;
```

### **Error: Database connection failed**
```bash
# Cek environment variables
echo $DB_HOST $DB_PORT $DB_USER $DB_NAME

# Atau set manual
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=e-bpjstk
```

## ✅ **Verification Checklist**

- [ ] Field `telepon_area_rumah` ada
- [ ] Field `telepon_area_kantor` ada  
- [ ] Field `telepon_kantor` ada
- [ ] Field `telepon_ext_kantor` ada
- [ ] Field `jenis_identitas` ada
- [ ] Field `masa_laku_identitas` ada
- [ ] Field `surat_menyurat_ke` ada
- [ ] Field `tanggal_kepesertaan` ada
- [ ] Field `kode_negara` ada
- [ ] Upload Excel berhasil
- [ ] Data tersimpan di field yang benar

## 🎯 **Next Steps**

Setelah migration selesai:
1. **Restart backend** untuk load model baru
2. **Test upload Excel** dengan template 31 kolom
3. **Verifikasi data** tersimpan di field yang benar
4. **Test form manual** untuk tambah individu

## 📞 **Support**

Jika ada masalah dengan migration:
1. Cek log error di terminal
2. Verifikasi koneksi database
3. Pastikan user memiliki permission yang cukup
4. Cek apakah field sudah ada sebelumnya







