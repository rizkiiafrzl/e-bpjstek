// Tambahkan route endpoint baru untuk kalkulasi iuran dan program BPJS
// File: backend/routes/worker.go (atau di main.go)

// Tambahkan route berikut ke dalam fungsi setupRoutes atau main:

// Endpoint untuk kalkulasi iuran BPJS
app.Post("/api/v1/workers/calculate-iuran", handlers.KalkulasiIuranEndpoint(db))

// Endpoint untuk mendapatkan daftar program BPJS
app.Get("/api/v1/workers/programs", handlers.GetProgramBPJS(db))

// Contoh penggunaan endpoint kalkulasi iuran:
/*
POST /api/v1/workers/calculate-iuran
Content-Type: application/json
Authorization: Bearer <token>

{
  "upah": 5000000,
  "program": {
    "jkk": true,
    "jkm": true,
    "jht": true,
    "jp": true
  },
  "risikoUsaha": 0.54
}

Response:
{
  "upah": 5000000,
  "program": {
    "jkk": true,
    "jkm": true,
    "jht": true,
    "jp": true
  },
  "risikoUsaha": 0.54,
  "iuran": {
    "upah": 5000000,
    "risikoUsaha": 0.54,
    "program": {...},
    "iuranJKK": 27000,
    "iuranJKM": 15000,
    "iuranJHTPerusahaan": 185000,
    "iuranJHTTK": 100000,
    "iuranJPPerusahaan": 100000,
    "iuranJPTK": 50000,
    "totalDitanggungPerusahaan": 327000,
    "totalDitanggungTK": 150000,
    "totalIuran": 477000
  }
}
*/

// Contoh penggunaan endpoint program BPJS:
/*
GET /api/v1/workers/programs
Authorization: Bearer <token>

Response:
{
  "programs": {
    "jkk": {
      "name": "Jaminan Kecelakaan Kerja",
      "rate": "0.24% - 1.74% (tergantung risiko usaha)",
      "company": "100%",
      "worker": "0%"
    },
    "jkm": {
      "name": "Jaminan Kematian",
      "rate": "0.3%",
      "company": "100%",
      "worker": "0%"
    },
    "jht": {
      "name": "Jaminan Hari Tua",
      "rate": "5.7%",
      "company": "3.7%",
      "worker": "2%"
    },
    "jp": {
      "name": "Jaminan Pensiun",
      "rate": "3%",
      "company": "2%",
      "worker": "1%"
    }
  },
  "description": "Program kepesertaan BPJS Ketenagakerjaan sesuai UU No. 24 Tahun 2011"
}
*/







