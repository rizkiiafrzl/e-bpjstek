# Farm Management Backend

Backend API untuk aplikasi Farm Management menggunakan Go Fiber dan PostgreSQL.

## 🚀 Fitur

- **Authentication**: Register, Login, JWT Token
- **User Management**: Profile management
- **Farm Management**: CRUD operations untuk farm
- **Database**: PostgreSQL dengan GORM
- **Security**: JWT authentication, password hashing
- **CORS**: Cross-origin resource sharing

## 📋 Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Git

## 🛠️ Installation

1. **Clone repository**

   ```bash
   git clone <repository-url>
   cd farm-management-backend
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Setup database**
   - Install PostgreSQL
   - Create database: `farm_management`
   - Update `.env` file with your database credentials

4. **Environment setup**
   ```bash
   cp env.example .env
   # Edit .env file with your configuration
   ```

## 🏃‍♂️ Running

### Development

```bash
# Using script (Linux/Mac)
./run.sh

# Or manually
go run main.go
```

### Production

```bash
# Build
go build -o farm-management-api main.go

# Run
./farm-management-api
```

## 📡 API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/google` - Google OAuth (coming soon)
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

## 🔧 Configuration

Edit `.env` file:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=farm_management
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRE_HOURS=24

# Server
PORT=8080
CORS_ORIGIN=http://localhost:5173
```

## 🗄️ Database Schema

### Users Table

- `id` (Primary Key)
- `email` (Unique)
- `password` (Hashed)
- `full_name`
- `avatar`
- `is_active`
- `created_at`
- `updated_at`
- `deleted_at`

### Farms Table

- `id` (Primary Key)
- `name`
- `description`
- `location`
- `area` (hectares)
- `user_id` (Foreign Key)
- `created_at`
- `updated_at`
- `deleted_at`

## 🔒 Security Features

- JWT Authentication
- Password hashing with bcrypt
- CORS protection
- Input validation
- SQL injection protection (GORM)

## 📝 API Examples

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login User

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Create Farm (with JWT token)

```bash
curl -X POST http://localhost:8080/api/v1/farms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "My Farm",
    "description": "A beautiful farm",
    "location": "Jakarta, Indonesia",
    "area": 10.5
  }'
```

## 🐛 Troubleshooting

### Database Connection Issues

- Check PostgreSQL is running
- Verify database credentials in `.env`
- Ensure database `farm_management` exists

### Port Already in Use

- Change `PORT` in `.env` file
- Or kill process using the port

### JWT Token Issues

- Check `JWT_SECRET` in `.env`
- Ensure token is not expired
- Verify Authorization header format

## 📚 Dependencies

- **Fiber v2**: Web framework
- **GORM**: ORM for database
- **JWT**: Authentication
- **bcrypt**: Password hashing
- **PostgreSQL**: Database driver
- **Validator**: Input validation

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## 📄 License

MIT License
