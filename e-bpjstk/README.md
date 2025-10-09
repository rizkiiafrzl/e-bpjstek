# Farm Management System

Sistem manajemen pertanian modern dengan frontend Vue.js dan backend Go Fiber + PostgreSQL.

## 🚀 Fitur

### Frontend (Vue.js 3)

- **Authentication**: Login, Register, JWT Token
- **Responsive Design**: Layout dua kolom dengan gambar landscape
- **Modern UI**: Glass morphism, gradient backgrounds
- **Form Validation**: Real-time validation
- **API Integration**: Terintegrasi dengan backend Go

### Backend (Go Fiber + PostgreSQL)

- **RESTful API**: Endpoint lengkap untuk CRUD operations
- **Authentication**: JWT-based authentication
- **Database**: PostgreSQL dengan GORM ORM
- **Security**: Password hashing, CORS protection
- **Farm Management**: CRUD operations untuk farm data

## 📁 Struktur Proyek

```
farm-management/
├── src/                          # Frontend Vue.js
│   ├── assets/
│   │   ├── auth.css             # Styling terpisah untuk auth pages
│   │   ├── base.css             # Global styles
│   │   └── main.css              # Main styles
│   ├── components/
│   │   └── icons/                # SVG icon components
│   ├── services/
│   │   └── api.js                # API service untuk backend calls
│   ├── views/
│   │   ├── AuthView.vue          # Combined login/register
│   │   ├── LoginView.vue         # Standalone login
│   │   ├── RegisterView.vue      # Standalone register
│   │   └── DashboardView.vue     # Dashboard setelah login
│   ├── router/
│   │   └── index.js              # Vue Router configuration
│   └── main.js                   # Vue app entry point
├── backend/                      # Backend Go Fiber
│   ├── config/                   # Configuration
│   ├── database/                 # Database connection & migration
│   ├── handlers/                 # API handlers
│   ├── middleware/               # Authentication middleware
│   ├── models/                   # Database models
│   ├── main.go                   # Backend entry point
│   ├── go.mod                    # Go dependencies
│   └── README.md                 # Backend documentation
└── README.md                     # This file
```

## 🛠️ Installation & Setup

### Prerequisites

- Node.js 18+
- Go 1.21+
- PostgreSQL 12+
- Git

### 1. Frontend Setup

```bash
# Install dependencies
npm install

# Copy environment file
cp env.example .env

# Start development server
npm run dev
```

### 2. Backend Setup

```bash
# Navigate to backend directory
cd backend

# Install Go dependencies
go mod tidy

# Copy environment file
cp env.example .env

# Edit .env file with your database credentials
# Then run the backend
go run main.go
```

### 3. Database Setup

```sql
-- Create database
CREATE DATABASE farm_management;

-- Create user (optional)
CREATE USER farm_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE farm_management TO farm_user;
```

## 🎨 Design Features

### Layout Dua Kolom

- **Kolom Kiri**: Background landscape dengan matahari dan bukit hijau
- **Kolom Kanan**: Form login/register dalam card putih
- **Responsive**: Layout vertikal di mobile

### Styling Terpisah

- **`src/assets/auth.css`**: Semua styling untuk halaman authentication
- **Import CSS**: Menggunakan `@import '@/assets/auth.css'`
- **Consistent**: Styling yang sama di semua auth pages

### Visual Elements

- **Gradient Background**: Biru langit ke hijau sawah
- **Sun**: Bulat emas dengan shadow
- **Hills**: 3 lapis bukit hijau dengan gradient
- **Glass Morphism**: Card transparan dengan backdrop blur

## 🔧 API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/logout` - Logout user

### User Profile

- `GET /api/v1/profile` - Get user profile
- `PUT /api/v1/profile` - Update user profile

### Farm Management

- `GET /api/v1/farms` - Get all farms
- `POST /api/v1/farms` - Create farm
- `GET /api/v1/farms/:id` - Get specific farm
- `PUT /api/v1/farms/:id` - Update farm
- `DELETE /api/v1/farms/:id` - Delete farm

## 🚀 Running the Application

### Development Mode

1. **Start Backend** (Terminal 1):

   ```bash
   cd backend
   go run main.go
   ```

2. **Start Frontend** (Terminal 2):

   ```bash
   npm run dev
   ```

3. **Access Application**:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080

### Production Mode

1. **Build Frontend**:

   ```bash
   npm run build
   ```

2. **Build Backend**:
   ```bash
   cd backend
   go build -o farm-management-api main.go
   ```

## 📱 Responsive Design

### Desktop (>768px)

- Layout dua kolom side-by-side
- Gambar landscape di kiri
- Form di kanan dengan card styling

### Mobile (≤768px)

- Layout vertikal
- Gambar di atas (height: 200px)
- Form di bawah dengan full width

## 🔒 Security Features

- **JWT Authentication**: Secure token-based auth
- **Password Hashing**: bcrypt encryption
- **CORS Protection**: Cross-origin security
- **Input Validation**: Server-side validation
- **SQL Injection Protection**: GORM ORM

## 🎯 Key Features Implemented

### ✅ Styling Separation

- Semua CSS dipindah ke `src/assets/auth.css`
- Import menggunakan `@import '@/assets/auth.css'`
- Styling konsisten di semua auth pages

### ✅ Backend Go Fiber

- RESTful API dengan Go Fiber
- PostgreSQL database dengan GORM
- JWT authentication
- CRUD operations untuk farm management

### ✅ Frontend Integration

- API service untuk backend calls
- Real authentication dengan backend
- Error handling yang proper
- Token management

### ✅ Modern UI/UX

- Layout dua kolom yang responsive
- Glass morphism design
- Gradient backgrounds
- Smooth animations dan transitions

## 🐛 Troubleshooting

### Frontend Issues

- **Styling tidak muncul**: Pastikan `@import '@/assets/auth.css'` ada di semua auth pages
- **API calls gagal**: Check `VITE_API_URL` di `.env` file
- **CORS errors**: Pastikan backend CORS configuration benar

### Backend Issues

- **Database connection**: Check PostgreSQL credentials di `.env`
- **Port already in use**: Change `PORT` di `.env` file
- **JWT errors**: Verify `JWT_SECRET` configuration

## 📚 Dependencies

### Frontend

- Vue.js 3
- Vue Router
- Vite
- CSS3 (Flexbox, Grid, Animations)

### Backend

- Go Fiber v2
- GORM
- PostgreSQL driver
- JWT authentication
- bcrypt password hashing

## 🤝 Contributing

1. Fork repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## 📄 License

MIT License - feel free to use for your projects!

---

**Farm Management System** - Modern web application untuk manajemen pertanian dengan teknologi terdepan! 🌾🚀
