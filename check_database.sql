-- Script untuk memeriksa struktur database e-BPJSTK
-- Pastikan semua field SIPP sudah ada di tabel workers

-- 1. Cek struktur tabel workers
\echo '=== STRUKTUR TABEL WORKERS ==='
\d workers;

-- 2. Cek field yang diperlukan untuk SIPP
\echo '=== FIELD SIPP CHECK ==='
SELECT 
    column_name as "Field Name",
    data_type as "Data Type",
    character_maximum_length as "Max Length",
    is_nullable as "Nullable",
    column_default as "Default Value"
FROM information_schema.columns 
WHERE table_name = 'workers' 
ORDER BY ordinal_position;

-- 3. Cek field yang hilang untuk SIPP
\echo '=== MISSING FIELDS CHECK ==='
WITH required_fields AS (
    SELECT unnest(ARRAY[
        'telepon_area_rumah', 'telepon_area_kantor', 'telepon_kantor', 
        'telepon_ext_kantor', 'jenis_identitas', 'masa_laku_identitas',
        'surat_menyurat_ke', 'tanggal_kepesertaan', 'kode_negara'
    ]) as field_name
),
existing_fields AS (
    SELECT column_name 
    FROM information_schema.columns 
    WHERE table_name = 'workers'
)
SELECT 
    rf.field_name,
    CASE 
        WHEN ef.column_name IS NULL THEN '❌ MISSING'
        ELSE '✅ EXISTS'
    END as status
FROM required_fields rf
LEFT JOIN existing_fields ef ON rf.field_name = ef.column_name
ORDER BY rf.field_name;

-- 4. Cek data sample (jika ada)
\echo '=== SAMPLE DATA CHECK ==='
SELECT 
    id, no_pegawai, nama, nik, telepon_area_rumah, telepon_kantor,
    jenis_identitas, tanggal_kepesertaan, kode_negara
FROM workers 
LIMIT 5;

-- 5. Cek constraint dan index
\echo '=== CONSTRAINTS & INDEXES ==='
SELECT 
    tc.constraint_name,
    tc.constraint_type,
    kcu.column_name
FROM information_schema.table_constraints tc
JOIN information_schema.key_column_usage kcu 
    ON tc.constraint_name = kcu.constraint_name
WHERE tc.table_name = 'workers'
ORDER BY tc.constraint_type, kcu.column_name;







