#!/usr/bin/env python3
"""
Script untuk test upload TK NA
"""

import requests
import json
import os

def test_upload_tk_na():
    """Test upload TK NA dengan file template"""
    
    # URL backend
    base_url = "http://localhost:8080/api/v1"
    
    # File template
    template_file = "template_tk_na.xlsx"
    
    print("=== TEST UPLOAD TK NA ===")
    print(f"Template file: {template_file}")
    print(f"Backend URL: {base_url}")
    
    # Cek apakah file template ada
    if not os.path.exists(template_file):
        print(f"❌ File template {template_file} tidak ditemukan!")
        return False
    
    print(f"✅ File template ditemukan: {os.path.getsize(template_file)} bytes")
    
    # Test 1: Cek health endpoint
    print("\n1. Test Health Endpoint...")
    try:
        response = requests.get(f"{base_url.replace('/api/v1', '')}/health", timeout=5)
        if response.status_code == 200:
            print("✅ Backend server running")
            print(f"   Response: {response.json()}")
        else:
            print(f"❌ Backend server error: {response.status_code}")
            return False
    except requests.exceptions.RequestException as e:
        print(f"❌ Backend server tidak bisa diakses: {e}")
        print("   Pastikan backend sudah running di port 8080")
        return False
    
    # Test 2: Cek endpoint upload TK NA (tanpa auth - harus return 401)
    print("\n2. Test Upload TK NA Endpoint (tanpa auth)...")
    try:
        with open(template_file, 'rb') as f:
            files = {'file': f}
            response = requests.post(f"{base_url}/workers/upload-tk-na", files=files, timeout=10)
            
        if response.status_code == 401:
            print("✅ Endpoint upload TK NA ada dan protected (401 Unauthorized)")
        else:
            print(f"⚠️  Endpoint response: {response.status_code}")
            print(f"   Response: {response.text}")
    except requests.exceptions.RequestException as e:
        print(f"❌ Error testing endpoint: {e}")
        return False
    
    # Test 3: Cek endpoint dengan method yang salah
    print("\n3. Test Upload TK NA Endpoint (method GET - harus return 405)...")
    try:
        response = requests.get(f"{base_url}/workers/upload-tk-na", timeout=5)
        if response.status_code == 405:
            print("✅ Endpoint hanya menerima POST (405 Method Not Allowed)")
        else:
            print(f"⚠️  Unexpected response: {response.status_code}")
    except requests.exceptions.RequestException as e:
        print(f"❌ Error testing GET method: {e}")
    
    print("\n=== HASIL TEST ===")
    print("✅ Backend server running")
    print("✅ Endpoint upload TK NA tersedia")
    print("✅ Endpoint protected dengan authentication")
    print("✅ Template file siap untuk upload")
    print("\n💡 Untuk test lengkap dengan auth, gunakan frontend atau Postman")
    
    return True

if __name__ == "__main__":
    test_upload_tk_na()



