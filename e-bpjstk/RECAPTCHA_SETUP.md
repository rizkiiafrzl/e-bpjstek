# Setup reCAPTCHA untuk E-BPJS Ketenagakerjaan

## 1. Daftar Google reCAPTCHA

1. Kunjungi [Google reCAPTCHA Admin Console](https://www.google.com/recaptcha/admin)
2. Login dengan akun Google Anda
3. Klik "Create" untuk membuat site baru
4. Pilih reCAPTCHA v2 ("I'm not a robot" Checkbox)
5. Masukkan domain Anda (contoh: `localhost`, `yourdomain.com`)
6. Setelah dibuat, Anda akan mendapatkan:
   - **Site Key** (untuk frontend)
   - **Secret Key** (untuk backend)

## 2. Konfigurasi Environment Variables

Buat file `.env` di root project dengan isi:

```env
# API Configuration
VITE_API_URL=http://localhost:8080/api/v1

# reCAPTCHA Configuration
VITE_RECAPTCHA_SITE_KEY=your_recaptcha_site_key_here
VITE_RECAPTCHA_SECRET_KEY=your_recaptcha_secret_key_here
```

## 3. Backend Implementation

Di backend Anda, tambahkan endpoint untuk validasi reCAPTCHA:

```javascript
// Contoh untuk Node.js/Express
app.post('/api/v1/auth/validate-recaptcha', async (req, res) => {
  const { recaptchaToken } = req.body

  const response = await fetch('https://www.google.com/recaptcha/api/siteverify', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body: `secret=${process.env.RECAPTCHA_SECRET_KEY}&response=${recaptchaToken}`,
  })

  const data = await response.json()

  if (data.success) {
    res.json({ success: true })
  } else {
    res.status(400).json({ error: 'reCAPTCHA validation failed' })
  }
})
```

## 4. Testing

Untuk testing, Anda bisa menggunakan test keys:

- **Site Key**: `6LeIxAcTAAAAAJcZVRqyHh71UMIEGNQ_MXjiZKhI`
- **Secret Key**: `6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe`

## 5. Production Setup

1. Ganti test keys dengan production keys
2. Pastikan domain sudah terdaftar di Google reCAPTCHA
3. Test di production environment

## Fitur yang Sudah Diimplementasi

✅ reCAPTCHA di halaman Login (`LoginView.vue`)
✅ reCAPTCHA di halaman Auth gabungan (`AuthView.vue`)
✅ Validasi reCAPTCHA sebelum submit form
✅ Reset reCAPTCHA setelah error
✅ Error handling untuk reCAPTCHA
✅ API service untuk validasi reCAPTCHA

## Troubleshooting

### reCAPTCHA tidak muncul

- Pastikan `VITE_RECAPTCHA_SITE_KEY` sudah di-set
- Cek console browser untuk error
- Pastikan domain sudah terdaftar di Google reCAPTCHA

### reCAPTCHA selalu gagal

- Pastikan secret key benar di backend
- Cek apakah backend endpoint `/auth/validate-recaptcha` sudah ada
- Pastikan domain sesuai dengan yang terdaftar

### Error "reCAPTCHA error"

- Refresh halaman
- Cek koneksi internet
- Pastikan script reCAPTCHA bisa di-load



