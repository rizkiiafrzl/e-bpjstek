-- Script untuk mengecek data yang sudah diupload massal
-- Pastikan data Excel sudah tersimpan di database dengan benar

-- 1. Cek data terbaru yang diupload
\echo '=== DATA TERBARU YANG DIUPLOAD ==='
SELECT 
    id,
    no_pegawai,
    nama,
    nik,
    telepon_area_rumah,
    telepon_rumah,
    telepon_area_kantor,
    telepon_kantor,
    telepon_ext_kantor,
    handphone,
    email,
    tempat_lahir,
    jenis_kelamin,
    gol_darah,
    status_pegawai,
    upah,
    npwp,
    jenis_identitas,
    tanggal_kepesertaan,
    kode_negara,
    created_at
FROM workers 
ORDER BY created_at DESC 
LIMIT 10;

-- 2. Cek field yang baru ditambahkan (harus ada data)
\echo '=== CEK FIELD BARU ADA DATANYA ==='
SELECT 
    COUNT(*) as total_workers,
    COUNT(telepon_area_rumah) as ada_telepon_area_rumah,
    COUNT(telepon_area_kantor) as ada_telepon_area_kantor,
    COUNT(telepon_kantor) as ada_telepon_kantor,
    COUNT(telepon_ext_kantor) as ada_telepon_ext_kantor,
    COUNT(jenis_identitas) as ada_jenis_identitas,
    COUNT(masa_laku_identitas) as ada_masa_laku_identitas,
    COUNT(surat_menyurat_ke) as ada_surat_menyurat_ke,
    COUNT(tanggal_kepesertaan) as ada_tanggal_kepesertaan,
    COUNT(kode_negara) as ada_kode_negara
FROM workers;

-- 3. Cek data yang diupload hari ini
\echo '=== DATA UPLOAD HARI INI ==='
SELECT 
    COUNT(*) as total_upload_hari_ini,
    COUNT(CASE WHEN telepon_area_rumah IS NOT NULL THEN 1 END) as dengan_telepon_area_rumah,
    COUNT(CASE WHEN jenis_identitas IS NOT NULL THEN 1 END) as dengan_jenis_identitas,
    COUNT(CASE WHEN tanggal_kepesertaan IS NOT NULL THEN 1 END) as dengan_tanggal_kepesertaan
FROM workers 
WHERE DATE(created_at) = CURRENT_DATE;

-- 4. Cek mapping Excel yang benar (contoh data)
\echo '=== CONTOH MAPPING EXCEL ==='
SELECT 
    'Excel A' as kolom_excel, 
    no_pegawai as data_database, 
    'NO_PEGAWAI' as field_name
FROM workers 
WHERE no_pegawai IS NOT NULL 
LIMIT 3

UNION ALL

SELECT 
    'Excel B', nama, 'NAMA_LENGKAP'
FROM workers 
WHERE nama IS NOT NULL 
LIMIT 3

UNION ALL

SELECT 
    'Excel D', telepon_area_rumah, 'TELEPON_AREA_RUMAH'
FROM workers 
WHERE telepon_area_rumah IS NOT NULL 
LIMIT 3

UNION ALL

SELECT 
    'Excel K', email, 'EMAIL'
FROM workers 
WHERE email IS NOT NULL 
LIMIT 3

UNION ALL

SELECT 
    'Excel P', nik, 'NO_IDENTITAS'
FROM workers 
WHERE nik IS NOT NULL 
LIMIT 3

UNION ALL

SELECT 
    'Excel Y', upah::text, 'UPAH'
FROM workers 
WHERE upah > 0 
LIMIT 3;

-- 5. Cek upload history
\echo '=== UPLOAD HISTORY ==='
SELECT 
    id,
    file_name,
    total_data,
    total_valid,
    total_invalid,
    validation_status,
    type,
    created_at
FROM upload_histories 
ORDER BY created_at DESC 
LIMIT 5;

-- 6. Cek data yang mungkin error (untuk debugging)
\echo '=== DATA YANG MUNGKIN ERROR ==='
SELECT 
    id,
    no_pegawai,
    nama,
    nik,
    CASE 
        WHEN LENGTH(nik) != 16 THEN 'NIK tidak 16 digit'
        WHEN nik IS NULL THEN 'NIK kosong'
        ELSE 'OK'
    END as nik_status,
    CASE 
        WHEN email NOT LIKE '%@%' THEN 'Email format salah'
        WHEN email IS NULL THEN 'Email kosong'
        ELSE 'OK'
    END as email_status
FROM workers 
WHERE created_at >= CURRENT_DATE - INTERVAL '1 day'
ORDER BY created_at DESC;

-- 7. Cek apakah data Excel sudah masuk dengan benar
\echo '=== VERIFIKASI DATA EXCEL ==='
SELECT 
    'Total Records' as check_type,
    COUNT(*) as count
FROM workers

UNION ALL

SELECT 
    'Records with Excel Fields',
    COUNT(*)
FROM workers 
WHERE telepon_area_rumah IS NOT NULL 
   OR telepon_area_kantor IS NOT NULL 
   OR jenis_identitas IS NOT NULL 
   OR tanggal_kepesertaan IS NOT NULL

UNION ALL

SELECT 
    'Records Today',
    COUNT(*)
FROM workers 
WHERE DATE(created_at) = CURRENT_DATE

UNION ALL

SELECT 
    'Records This Week',
    COUNT(*)
FROM workers 
WHERE created_at >= CURRENT_DATE - INTERVAL '7 days';







