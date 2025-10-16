# 🐛 Debug Upload Massal Tidak Muncul di Tabel

## 🔍 **Analisis Masalah dari Screenshot:**

Dari screenshot terlihat:
- ✅ **Tabel ada** - Struktur tabel sudah benar
- ✅ **1 data tampil** - "devi" dengan NIK 1234567890123456
- ❌ **Data upload massal tidak muncul** - Seharusnya ada lebih banyak data

## 🔧 **Kemungkinan Penyebab:**

### **1. Data Upload Gagal Disimpan**
- Validasi data gagal
- Error saat insert ke database
- Field mapping salah

### **2. Data Tersimpan Tapi Tidak Tampil**
- Filter/query salah di frontend
- Pagination issue
- User ID tidak sesuai

### **3. Database Schema Issue**
- Field baru belum ada di database
- Migration belum dijalankan

## 🚀 **Langkah Debugging:**

### **Step 1: Cek Upload History**
```bash
# Cek via API
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://localhost:8080/api/v1/workers/upload-history
```

### **Step 2: Cek Data di Database**
```sql
-- Cek data terbaru
SELECT id, no_pegawai, nama, nik, created_at, user_id 
FROM workers 
ORDER BY created_at DESC 
LIMIT 10;

-- Cek data hari ini
SELECT COUNT(*) as total_today
FROM workers 
WHERE DATE(created_at) = CURRENT_DATE;
```

### **Step 3: Cek Log Backend**
```bash
# Lihat log saat upload
tail -f backend.log | grep -i "upload\|error\|worker"
```

### **Step 4: Test Upload Kecil**
Upload Excel dengan 2-3 data saja untuk test

## 🔍 **Script Debug Cepat:**

### **1. Cek Data via API**
```bash
# Windows PowerShell
$headers = @{"Authorization" = "Bearer YOUR_TOKEN"}
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/workers" -Headers $headers
```

### **2. Cek Upload History**
```bash
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/workers/upload-history" -Headers $headers
```

### **3. Cek Database Langsung**
```sql
-- Cek semua data
SELECT * FROM workers ORDER BY created_at DESC;

-- Cek data per user
SELECT user_id, COUNT(*) as total 
FROM workers 
GROUP BY user_id;
```

## ⚠️ **Kemungkinan Masalah Spesifik:**

### **1. User ID Mismatch**
- Data tersimpan dengan user_id berbeda
- Frontend filter berdasarkan user yang salah

### **2. Validasi Data Gagal**
- NIK duplikat
- Format data salah
- Field required kosong

### **3. Database Field Missing**
- Field baru belum ada di database
- Migration belum dijalankan

### **4. Frontend Filter Issue**
- Query parameter salah
- Pagination bug
- Search filter aktif

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

### **3. Test Upload Ulang**
- Upload Excel dengan data minimal
- Cek response API
- Lihat log backend

### **4. Cek Frontend**
- Refresh halaman
- Cek filter/search
- Lihat network tab di browser

## 📋 **Checklist Debug:**

- [ ] Backend running dan log tidak ada error
- [ ] Database migration sudah dijalankan
- [ ] Upload history menampilkan data
- [ ] Data ada di database dengan user_id yang benar
- [ ] Frontend tidak ada filter yang aktif
- [ ] Network request berhasil (status 200)
- [ ] Response API mengembalikan data

## 🎯 **Next Steps:**

1. **Jalankan script debug** di atas
2. **Cek log backend** saat upload
3. **Verifikasi data** di database
4. **Test upload** dengan data minimal
5. **Report hasil** untuk analisis lebih lanjut







