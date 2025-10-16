-- Database Migration Script untuk e-BPJSTK
-- Menambahkan field yang hilang di tabel workers sesuai SIPP

-- Cek struktur tabel workers saat ini
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns 
WHERE table_name = 'workers' 
ORDER BY ordinal_position;

-- Tambahkan field yang hilang jika belum ada
-- Field Contact Info
ALTER TABLE workers 
ADD COLUMN IF NOT EXISTS telepon_area_rumah VARCHAR(8);

ALTER TABLE workers 
ADD COLUMN IF NOT EXISTS telepon_area_kantor VARCHAR(8);

ALTER TABLE workers 
ADD COLUMN IF NOT EXISTS telepon_kantor VARCHAR(32);

ALTER TABLE workers 
ADD COLUMN IF NOT EXISTS telepon_ext_kantor VARCHAR(16);

-- Field Additional Info
ALTER TABLE workers 
ADD COLUMN IF NOT EXISTS jenis_identitas VARCHAR(16);

ALTER TABLE workers 
ADD COLUMN IF NOT EXISTS masa_laku_identitas VARCHAR(16);

ALTER TABLE workers 
ADD COLUMN IF NOT EXISTS surat_menyurat_ke VARCHAR(200);

ALTER TABLE workers 
ADD COLUMN IF NOT EXISTS tanggal_kepesertaan VARCHAR(16);

ALTER TABLE workers 
ADD COLUMN IF NOT EXISTS kode_negara VARCHAR(4);

-- Cek struktur tabel setelah migration
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns 
WHERE table_name = 'workers' 
ORDER BY ordinal_position;

-- Verifikasi field yang diperlukan untuk SIPP
SELECT 
    'Field SIPP Check' as check_type,
    column_name,
    CASE 
        WHEN column_name IN (
            'no_pegawai', 'nama', 'nik', 'tempat_lahir', 'date_of_birth',
            'ibu_kandung', 'jenis_kelamin', 'gol_darah', 'status_kawin',
            'status_pegawai', 'tanggal_awal_bekerja', 'tanggal_akhir_kontrak',
            'lokasi_pekerjaan', 'telepon_area_rumah', 'telepon_rumah',
            'telepon_area_kantor', 'telepon_kantor', 'telepon_ext_kantor',
            'handphone', 'email', 'alamat', 'kabupaten', 'kode_pos',
            'npwp', 'jenis_identitas', 'masa_laku_identitas',
            'surat_menyurat_ke', 'tanggal_kepesertaan', 'kode_negara',
            'upah', 'rapel', 'nationality', 'passport_no', 'passport_valid_until'
        ) THEN '✅ Required Field'
        ELSE '⚠️ Additional Field'
    END as status
FROM information_schema.columns 
WHERE table_name = 'workers' 
ORDER BY column_name;







