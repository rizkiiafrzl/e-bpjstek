#!/usr/bin/env python3
"""
Script untuk membuat file Excel test dengan data valid untuk upload TK
"""

import pandas as pd
from datetime import datetime, timedelta
import random

def create_test_excel():
    """Membuat file Excel test dengan data valid"""
    
    # Data test dengan format yang benar
    test_data = [
        {
            'NO_PEGAWAI': 'PEG001',
            'NAMA_LENGKAP': 'John Doe',
            'GELAR': 'S.Kom',
            'TELEPON_AREA_RUMAH': '021',
            'TELEPON_RUMAH': '7654321',
            '': '',  # Kolom F kosong
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
            'STATUS_PEGAWAI': 'Tetap',
            'TGL_AWAL_BEKERJA': '2020-07-01',
            'TGL_AKHIR_KONTRAK': ''
        },
        {
            'NO_PEGAWAI': 'PEG002',
            'NAMA_LENGKAP': 'Jane Smith',
            'GELAR': 'S.E',
            'TELEPON_AREA_RUMAH': '022',
            'TELEPON_RUMAH': '8765432',
            '': '',  # Kolom F kosong
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
            'STATUS_PEGAWAI': 'Tetap',
            'TGL_AWAL_BEKERJA': '2021-01-15',
            'TGL_AKHIR_KONTRAK': ''
        },
        {
            'NO_PEGAWAI': 'PEG003',
            'NAMA_LENGKAP': 'Ahmad Rahman',
            'GELAR': 'S.T',
            'TELEPON_AREA_RUMAH': '024',
            'TELEPON_RUMAH': '9876543',
            '': '',  # Kolom F kosong
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
            'STATUS_PEGAWAI': 'Kontrak',
            'TGL_AWAL_BEKERJA': '2022-03-10',
            'TGL_AKHIR_KONTRAK': '2024-03-10'
        }
    ]
    
    # Buat DataFrame
    df = pd.DataFrame(test_data)
    
    # Simpan ke Excel dengan format yang benar
    filename = 'template_tk_test_valid.xlsx'
    
    with pd.ExcelWriter(filename, engine='openpyxl') as writer:
        df.to_excel(writer, sheet_name='Sheet1', index=False, header=True)
        
        # Dapatkan worksheet untuk formatting
        worksheet = writer.sheets['Sheet1']
        
        # Set column widths
        column_widths = {
            'A': 12,  # NO_PEGAWAI
            'B': 20,  # NAMA_LENGKAP
            'C': 10,  # GELAR
            'D': 15,  # TELEPON_AREA_RUMAH
            'E': 15,  # TELEPON_RUMAH
            'F': 5,   # (KOSONG)
            'G': 15,  # TELEPON_AREA_KANTOR
            'H': 15,  # TELEPON_KANTOR
            'I': 15,  # TELEPON_EXT_KANTOR
            'J': 15,  # HP
            'K': 25,  # EMAIL
            'L': 15,  # TEMPAT_LAHIR
            'M': 15,  # TANGGAL_LAHIR
            'N': 20,  # NAMA_IBU_KANDUNG
            'O': 15,  # JENIS_IDENTITAS
            'P': 20,  # NO_IDENTITAS
            'Q': 15,  # MASA_LAKU_IDENTITAS
            'R': 15,  # JENIS_KELAMIN
            'S': 25,  # SURAT_MENYURAT_KE
            'T': 15,  # TANGGAL_KEPESERTAAN
            'U': 15,  # STATUS_KAWIN
            'V': 15,  # GOLONGAN_DARAH
            'W': 20,  # NPWP
            'X': 10,  # KODE_NEGARA
            'Y': 15,  # UPAH
            'Z': 25,  # ALAMAT
            'AA': 10, # KODE_POS
            'AB': 20, # LOKASI_PEKERJAAN
            'AC': 15, # STATUS_PEGAWAI
            'AD': 15, # TGL_AWAL_BEKERJA
            'AE': 15  # TGL_AKHIR_KONTRAK
        }
        
        for col, width in column_widths.items():
            worksheet.column_dimensions[col].width = width
    
    print(f"✅ File Excel test berhasil dibuat: {filename}")
    print(f"📊 Total data: {len(test_data)}")
    print(f"📋 Kolom: {len(df.columns)} (A-AE)")
    
    return filename

if __name__ == "__main__":
    create_test_excel()







