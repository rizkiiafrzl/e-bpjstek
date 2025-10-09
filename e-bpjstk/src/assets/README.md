# Assets Documentation

## 🖼️ Gambar yang Digunakan

### 1. Background Image

- **File**: `src/components/icons/background.jpg`
- **Penggunaan**: Background utama untuk halaman authentication
- **CSS**: `.auth-container { background: url('@/components/icons/background.jpg') center/cover no-repeat; }`

### 2. Logo BPJS

- **File**: `src/components/icons/Log BPJS Ketenagakerjaan.png`
- **Penggunaan**: Logo di header form login/register
- **CSS**: `.logo-icon { background: url('@/components/icons/Log BPJS Ketenagakerjaan.png') center/contain no-repeat; }`

### 3. Side Image

- **File**: `src/components/icons/samping login.jpg`
- **Penggunaan**: Gambar di samping form (kolom kiri)
- **CSS**: `.auth-image { background: url('@/components/icons/samping login.jpg') center/cover no-repeat; }`

## 🎨 Layout Structure

```
┌─────────────────────────────────────────────────────────┐
│                    Background Image                     │
│  ┌─────────────────────────────────────────────────┐   │
│  │              Auth Layout                        │   │
│  │  ┌─────────────┐  ┌─────────────────────────┐   │   │
│  │  │             │  │                         │   │   │
│  │  │ Side Image  │  │     Form Container      │   │   │
│  │  │             │  │  ┌─────────────────┐    │   │   │
│  │  │             │  │  │   Logo BPJS     │    │   │   │
│  │  │             │  │  │                 │    │   │   │
│  │  │             │  │  │   Login Form    │    │   │   │
│  │  │             │  │  │                 │    │   │   │
│  │  │             │  │  └─────────────────┘    │   │   │
│  │  └─────────────┘  └─────────────────────────┘   │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

## 📱 Responsive Behavior

### Desktop (>768px)

- **Layout**: Dua kolom side-by-side
- **Left Column**: Side image (50% width)
- **Right Column**: Form container (50% width)

### Mobile (≤768px)

- **Layout**: Vertikal
- **Top**: Side image (height: 200px)
- **Bottom**: Form container (full width)

## 🔧 CSS Configuration

### Background Container

```css
.auth-container {
  background: url('@/components/icons/background.jpg') center/cover no-repeat;
}
```

### Side Image

```css
.auth-image {
  background: url('@/components/icons/samping login.jpg') center/cover no-repeat;
}
```

### Logo

```css
.logo-icon {
  background: url('@/components/icons/Log BPJS Ketenagakerjaan.png') center/contain no-repeat;
  background-size: 100% 100%;
}
```

## 📝 Notes

- Semua gambar menggunakan path `@/components/icons/` untuk konsistensi
- Background menggunakan `center/cover` untuk fill penuh
- Logo menggunakan `center/contain` untuk proporsi yang tepat
- Side image menggunakan `center/cover` untuk fill area
- Responsive design memastikan gambar terlihat baik di semua ukuran layar





