#!/usr/bin/env python3
"""
Script untuk debug upload TK dan melihat error validasi
"""

import requests
import json

def debug_upload():
    """Debug upload TK untuk melihat error detail"""
    
    # URL API
    api_url = "http://localhost:8080/api/v1/workers/upload-tk"
    
    # Token JWT (ganti dengan token yang valid)
    token = "YOUR_JWT_TOKEN_HERE"
    
    if token == "YOUR_JWT_TOKEN_HERE":
        print("❌ ERROR: Ganti YOUR_JWT_TOKEN_HERE dengan token JWT yang valid")
        print("💡 Cara mendapatkan token:")
        print("   1. Login ke aplikasi")
        print("   2. Buka Developer Tools (F12)")
        print("   3. Lihat di Network tab saat login")
        print("   4. Copy token dari header Authorization")
        return
    
    # Headers
    headers = {
        "Authorization": f"Bearer {token}"
    }
    
    # File yang akan diupload
    file_path = "template_sipp_valid.xlsx"
    
    try:
        # Upload file
        with open(file_path, 'rb') as f:
            files = {'file': (file_path, f, 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet')}
            
            print(f"🚀 Uploading file: {file_path}")
            response = requests.post(api_url, headers=headers, files=files)
        
        print(f"📊 Status Code: {response.status_code}")
        print(f"📋 Response Headers: {dict(response.headers)}")
        
        if response.status_code == 200:
            result = response.json()
            print("✅ Upload berhasil!")
            print(f"📈 Total Data: {result.get('total', 'N/A')}")
            print(f"✅ Valid: {result.get('valid', 'N/A')}")
            print(f"❌ Invalid: {result.get('invalid', 'N/A')}")
            
            if 'errors' in result:
                print(f"🔍 Error Details:")
                for error in result['errors']:
                    print(f"   - {error}")
            
            if 'history_id' in result:
                print(f"📝 History ID: {result['history_id']}")
                
                # Cek detail upload history
                history_url = f"http://localhost:8080/api/v1/upload-histories/{result['history_id']}"
                history_response = requests.get(history_url, headers=headers)
                
                if history_response.status_code == 200:
                    history = history_response.json()
                    print(f"📊 Upload History Details:")
                    print(f"   - File Name: {history.get('file_name', 'N/A')}")
                    print(f"   - Total Data: {history.get('total_data', 'N/A')}")
                    print(f"   - Valid: {history.get('total_valid', 'N/A')}")
                    print(f"   - Invalid: {history.get('total_invalid', 'N/A')}")
                    print(f"   - Status: {history.get('validation_status', 'N/A')}")
                    
                    if 'errors' in history and history['errors']:
                        print(f"🔍 Validation Errors:")
                        for error in history['errors']:
                            print(f"   - {error}")
        else:
            print(f"❌ Upload gagal!")
            print(f"📋 Response: {response.text}")
            
    except FileNotFoundError:
        print(f"❌ File tidak ditemukan: {file_path}")
        print("💡 Pastikan file template_sipp_valid.xlsx ada di direktori ini")
    except requests.exceptions.ConnectionError:
        print("❌ Tidak bisa connect ke backend")
        print("💡 Pastikan backend berjalan di http://localhost:8080")
    except Exception as e:
        print(f"❌ Error: {e}")

def test_individual_worker():
    """Test tambah TK individual untuk debug validasi"""
    
    api_url = "http://localhost:8080/api/v1/workers"
    token = "YOUR_JWT_TOKEN_HERE"
    
    if token == "YOUR_JWT_TOKEN_HERE":
        print("❌ ERROR: Ganti YOUR_JWT_TOKEN_HERE dengan token JWT yang valid")
        return
    
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    
    # Data test sesuai template SIPP
    test_data = {
        "noPegawai": "PEG001",
        "nama": "John Doe",
        "gelar": "S.Kom",
        "teleponAreaRumah": "021",
        "teleponRumah": "7654321",
        "teleponAreaKantor": "021",
        "teleponKantor": "021-1234",
        "teleponExtKantor": "1234",
        "handphone": "081234567890",
        "email": "john.doe@example.com",
        "tempatLahir": "Jakarta",
        "dateOfBirth": "1990-05-14",
        "ibuKandung": "Siti Aminah",
        "jenisIdentitas": "KTP",
        "nik": "3175051405900001",
        "masaLakuIdentitas": "2016-05-14",
        "jenisKelamin": "Laki-laki",
        "suratMenyuratKe": "Jalan Mandaka No. 10",
        "tanggalKepesertaan": "2020-07-01",
        "statusKawin": "Belum Kawin",
        "golDarah": "O",
        "npwp": "012345678901234",
        "kodeNegara": "ID",
        "upah": 5000000,
        "alamat": "Jl. Merpati 12 Jakarta",
        "kodePos": "10130",
        "lokasiPekerjaan": "Kantor Pusat",
        "statusPegawai": "PKWTT",
        "tanggalAwalBekerja": "2020-07-01",
        "tanggalAkhirKontrak": "",
        "rapel": 0,
        "nationality": "WNI"
    }
    
    try:
        print("🚀 Testing individual worker creation...")
        response = requests.post(api_url, headers=headers, json=test_data)
        
        print(f"📊 Status Code: {response.status_code}")
        
        if response.status_code == 201:
            result = response.json()
            print("✅ Worker berhasil dibuat!")
            print(f"👤 Worker ID: {result['worker']['id']}")
            print(f"💰 Iuran Total: Rp {result['iuran']['totalIuran']:,.0f}")
        else:
            print(f"❌ Worker creation gagal!")
            print(f"📋 Response: {response.text}")
            
    except Exception as e:
        print(f"❌ Error: {e}")

if __name__ == "__main__":
    print("🔍 Debug Upload TK - SIPP Compliance")
    print("=" * 50)
    
    print("\n1. Testing Individual Worker Creation:")
    test_individual_worker()
    
    print("\n2. Testing Mass Upload:")
    debug_upload()
    
    print("\n" + "=" * 50)
    print("💡 Tips:")
    print("   - Pastikan backend berjalan di port 8080")
    print("   - Ganti YOUR_JWT_TOKEN_HERE dengan token yang valid")
    print("   - File template_sipp_valid.xlsx harus ada di direktori ini")







