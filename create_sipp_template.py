#!/usr/bin/env python3
"""
Script untuk membuat file Excel template SIPP yang benar-benar valid
Berdasarkan spesifikasi lengkap dari dokumentasi SIPP BPJamsostek
"""

import pandas as pd
from datetime import datetime, timedelta
import random

def create_sipp_template():
    """Membuat template Excel SIPP yang benar-benar valid"""
    
    # Data test dengan format yang benar sesuai SIPP
    test_data = [
        {
            'NO_PEGAWAI': 'PEG001',
            'NAMA_LENGKAP': 'John Doe',
            'GELAR': 'S.Kom',
            'TELEPON_AREA_RUMAH': '021',
            'TELEPON_RUMAH': '7654321',
            'TELEPON_AREA_KANTOR': '021',
            'TELEPON_KANTOR': '021-1234',
            'TELEPON_EXT_KANTOR': '1234',
            'HP': '081234567890',
            'EMAIL': 'john.doe@example.com',
            'TEMPAT_LAHIR': 'Jakarta',
            'TANGGAL_LAHIR': '1990-05-14',
            'NAMA_IBU_KANDUNG': 'Siti Aminah',
            'JENIS_IDENTITAS': 'KTP',
            'NO_IDENTITAS': '3175051405900001',
            'MASA_LAKU_IDENTITAS': '2016-05-14',
            'JENIS_KELAMIN': 'Laki-laki',
            'SURAT_MENYURAT_KE': 'Jalan Mandaka No. 10',
            'TANGGAL_KEPESERTAAN': '2020-07-01',
            'STATUS_KAWIN': 'Belum Kawin',
            'GOLONGAN_DARAH': 'O',
            'NPWP': '012345678901234',
            'KODE_NEGARA': 'ID',
            'UPAH': '5000000',
            'ALAMAT': 'Jl. Merpati 12 Jakarta',
            'KODE_POS': '10130',
            'LOKASI_PEKERJAAN': 'Kantor Pusat',
            'STATUS_PEGAWAI': 'PKWTT',
            'TGL_AWAL_BEKERJA': '2020-07-01',
            'TGL_AKHIR_KONTRAK': '',
            'RAPEL': '0',
            'NATIONALITY': 'WNI'
        },
        {
            'NO_PEGAWAI': 'PEG002',
            'NAMA_LENGKAP': 'Jane Smith',
            'GELAR': 'S.E',
            'TELEPON_AREA_RUMAH': '022',
            'TELEPON_RUMAH': '8765432',
            'TELEPON_AREA_KANTOR': '022',
            'TELEPON_KANTOR': '022-5678',
            'TELEPON_EXT_KANTOR': '5678',
            'HP': '082345678901',
            'EMAIL': 'jane.smith@example.com',
            'TEMPAT_LAHIR': 'Bandung',
            'TANGGAL_LAHIR': '1988-03-22',
            'NAMA_IBU_KANDUNG': 'Maria Sari',
            'JENIS_IDENTITAS': 'KTP',
            'NO_IDENTITAS': '3175051405900002',
            'MASA_LAKU_IDENTITAS': '2018-03-22',
            'JENIS_KELAMIN': 'Perempuan',
            'SURAT_MENYURAT_KE': 'Jalan Sudirman No. 5',
            'TANGGAL_KEPESERTAAN': '2021-01-15',
            'STATUS_KAWIN': 'Kawin',
            'GOLONGAN_DARAH': 'A',
            'NPWP': '012345678901235',
            'KODE_NEGARA': 'ID',
            'UPAH': '6000000',
            'ALAMAT': 'Jl. Gatot Subroto 25 Bandung',
            'KODE_POS': '40112',
            'LOKASI_PEKERJAAN': 'Cabang Bandung',
            'STATUS_PEGAWAI': 'PKWTT',
            'TGL_AWAL_BEKERJA': '2021-01-15',
            'TGL_AKHIR_KONTRAK': '',
            'RAPEL': '0',
            'NATIONALITY': 'WNI'
        },
        {
            'NO_PEGAWAI': 'PEG003',
            'NAMA_LENGKAP': 'Ahmad Rahman',
            'GELAR': 'S.T',
            'TELEPON_AREA_RUMAH': '024',
            'TELEPON_RUMAH': '9876543',
            'TELEPON_AREA_KANTOR': '024',
            'TELEPON_KANTOR': '024-9012',
            'TELEPON_EXT_KANTOR': '9012',
            'HP': '083456789012',
            'EMAIL': 'ahmad.rahman@example.com',
            'TEMPAT_LAHIR': 'Semarang',
            'TANGGAL_LAHIR': '1992-11-08',
            'NAMA_IBU_KANDUNG': 'Fatimah Zahra',
            'JENIS_IDENTITAS': 'KTP',
            'NO_IDENTITAS': '3175051405900003',
            'MASA_LAKU_IDENTITAS': '2020-11-08',
            'JENIS_KELAMIN': 'Laki-laki',
            'SURAT_MENYURAT_KE': 'Jalan Diponegoro No. 15',
            'TANGGAL_KEPESERTAAN': '2022-03-10',
            'STATUS_KAWIN': 'Belum Kawin',
            'GOLONGAN_DARAH': 'B',
            'NPWP': '012345678901236',
            'KODE_NEGARA': 'ID',
            'UPAH': '4500000',
            'ALAMAT': 'Jl. Imam Bonjol 8 Semarang',
            'KODE_POS': '50111',
            'LOKASI_PEKERJAAN': 'Cabang Semarang',
            'STATUS_PEGAWAI': 'PKWT',
            'TGL_AWAL_BEKERJA': '2022-03-10',
            'TGL_AKHIR_KONTRAK': '2024-03-10',
            'RAPEL': '0',
            'NATIONALITY': 'WNI'
        }
    ]
    
    # Buat DataFrame dengan urutan kolom yang benar
    df = pd.DataFrame(test_data)
    
    # Reorder kolom sesuai template SIPP yang benar
    column_order = [
        'NO_PEGAWAI',           # A
        'NAMA_LENGKAP',         # B
        'GELAR',                # C
        'TELEPON_AREA_RUMAH',   # D
        'TELEPON_RUMAH',        # E
        'TELEPON_AREA_KANTOR',  # G
        'TELEPON_KANTOR',       # H
        'TELEPON_EXT_KANTOR',   # I
        'HP',                   # J
        'EMAIL',                # K
        'TEMPAT_LAHIR',         # L
        'TANGGAL_LAHIR',        # M
        'NAMA_IBU_KANDUNG',     # N
        'JENIS_IDENTITAS',      # O
        'NO_IDENTITAS',         # P
        'MASA_LAKU_IDENTITAS',  # Q
        'JENIS_KELAMIN',        # R
        'SURAT_MENYURAT_KE',    # S
        'TANGGAL_KEPESERTAAN',  # T
        'STATUS_KAWIN',         # U
        'GOLONGAN_DARAH',       # V
        'NPWP',                 # W
        'KODE_NEGARA',          # X
        'UPAH',                 # Y
        'ALAMAT',               # Z
        'KODE_POS',             # AA
        'LOKASI_PEKERJAAN',     # AB
        'STATUS_PEGAWAI',       # AC
        'TGL_AWAL_BEKERJA',     # AD
        'TGL_AKHIR_KONTRAK',    # AE
        'RAPEL',                # AF
        'NATIONALITY'           # AG
    ]
    
    # Reorder DataFrame
    df = df[column_order]
    
    # Simpan ke Excel dengan format yang benar
    filename = 'template_sipp_valid.xlsx'
    
    with pd.ExcelWriter(filename, engine='openpyxl') as writer:
        df.to_excel(writer, sheet_name='Sheet1', index=False, header=True)
        
        # Dapatkan worksheet untuk formatting
        worksheet = writer.sheets['Sheet1']
        
        # Set column widths untuk readability
        column_widths = {
            'A': 12,   # NO_PEGAWAI
            'B': 20,   # NAMA_LENGKAP
            'C': 10,   # GELAR
            'D': 15,   # TELEPON_AREA_RUMAH
            'E': 15,   # TELEPON_RUMAH
            'F': 5,    # (KOSONG - sesuai template)
            'G': 15,   # TELEPON_AREA_KANTOR
            'H': 15,   # TELEPON_KANTOR
            'I': 15,   # TELEPON_EXT_KANTOR
            'J': 15,   # HP
            'K': 25,   # EMAIL
            'L': 15,   # TEMPAT_LAHIR
            'M': 15,   # TANGGAL_LAHIR
            'N': 20,   # NAMA_IBU_KANDUNG
            'O': 15,   # JENIS_IDENTITAS
            'P': 20,   # NO_IDENTITAS
            'Q': 15,   # MASA_LAKU_IDENTITAS
            'R': 15,   # JENIS_KELAMIN
            'S': 25,   # SURAT_MENYURAT_KE
            'T': 15,   # TANGGAL_KEPESERTAAN
            'U': 15,   # STATUS_KAWIN
            'V': 15,   # GOLONGAN_DARAH
            'W': 20,   # NPWP
            'X': 10,   # KODE_NEGARA
            'Y': 15,   # UPAH
            'Z': 25,   # ALAMAT
            'AA': 10,  # KODE_POS
            'AB': 20,  # LOKASI_PEKERJAAN
            'AC': 15,  # STATUS_PEGAWAI
            'AD': 15,  # TGL_AWAL_BEKERJA
            'AE': 15,  # TGL_AKHIR_KONTRAK
            'AF': 10,  # RAPEL
            'AG': 10   # NATIONALITY
        }
        
        for col, width in column_widths.items():
            worksheet.column_dimensions[col].width = width
        
        # Tambahkan kolom F yang kosong (sesuai template SIPP)
        worksheet.insert_cols(6)  # Insert kolom F
        worksheet.cell(row=1, column=6).value = ''  # Header kosong
        
        # Format header row
        from openpyxl.styles import Font, PatternFill, Alignment
        
        header_font = Font(bold=True, color='FFFFFF')
        header_fill = PatternFill(start_color='366092', end_color='366092', fill_type='solid')
        header_alignment = Alignment(horizontal='center', vertical='center')
        
        for col in range(1, len(column_order) + 2):  # +2 karena ada kolom F kosong
            cell = worksheet.cell(row=1, column=col)
            cell.font = header_font
            cell.fill = header_fill
            cell.alignment = header_alignment
    
    print(f"✅ Template SIPP Excel berhasil dibuat: {filename}")
    print(f"📊 Total data: {len(test_data)}")
    print(f"📋 Kolom: {len(column_order) + 1} (A-AG, dengan kolom F kosong)")
    print(f"📝 Format: Sesuai spesifikasi SIPP BPJamsostek")
    
    return filename

if __name__ == "__main__":
    create_sipp_template()







