#!/bin/bash

# Script untuk mengecek data yang sudah diupload massal
# Pastikan data Excel sudah tersimpan di database dengan benar

echo "🔍 Checking uploaded data from Excel..."

# Set your API details
API_URL="http://localhost:8080/api/v1"
TOKEN="YOUR_TOKEN_HERE"  # Ganti dengan token yang valid

# 1. Cek data workers terbaru
echo ""
echo "=== DATA WORKERS TERBARU ==="
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$API_URL/workers" | jq '.[0:5] | .[] | {
        no_pegawai: .noPegawai,
        nama: .nama,
        nik: .nik,
        telepon_area_rumah: .teleponAreaRumah,
        telepon_kantor: .teleponKantor,
        email: .email,
        jenis_identitas: .jenisIdentitas,
        tanggal_kepesertaan: .tanggalKepesertaan,
        kode_negara: .kodeNegara,
        upah: .upah,
        created_at: .created_at
    }'

# 2. Cek upload history
echo ""
echo "=== UPLOAD HISTORY ==="
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$API_URL/workers/upload-history" | jq '.[0:5] | .[] | {
        file_name: .fileName,
        type: .type,
        total_data: .totalData,
        total_valid: .totalValid,
        total_invalid: .totalInvalid,
        status: .validationStatus,
        created_at: .created_at
    }'

# 3. Cek field mapping (contoh data)
echo ""
echo "=== FIELD MAPPING CHECK ==="
echo "Excel A (NO_PEGAWAI) → Database no_pegawai:"
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$API_URL/workers" | jq '.[0:3] | .[] | .noPegawai'

echo "Excel D (TELEPON_AREA_RUMAH) → Database telepon_area_rumah:"
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$API_URL/workers" | jq '.[0:3] | .[] | .teleponAreaRumah'

echo "Excel K (EMAIL) → Database email:"
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$API_URL/workers" | jq '.[0:3] | .[] | .email'

echo "Excel P (NO_IDENTITAS) → Database nik:"
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$API_URL/workers" | jq '.[0:3] | .[] | .nik'

echo "Excel Y (UPAH) → Database upah:"
curl -s -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     "$API_URL/workers" | jq '.[0:3] | .[] | .upah'

echo ""
echo "✅ Data check completed!"
echo ""
echo "📋 Jika data kosong atau tidak sesuai:"
echo "1. Pastikan backend sudah di-restart"
echo "2. Pastikan database migration sudah dijalankan"
echo "3. Cek log backend untuk error"
echo "4. Test upload Excel lagi"







