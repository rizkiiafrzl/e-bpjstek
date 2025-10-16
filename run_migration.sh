#!/bin/bash

# Script untuk menjalankan migration database e-BPJSTK
echo "🚀 Starting Database Migration for e-BPJSTK..."

# Set environment variables (sesuaikan dengan konfigurasi Anda)
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=e-bpjstk
export DB_SSLMODE=disable

# Navigate to backend directory
cd backend

echo "📁 Changed to backend directory"

# Run the migration script
echo "🔄 Running database migration..."
go run ../migrate_database.go

if [ $? -eq 0 ]; then
    echo "✅ Migration completed successfully!"
else
    echo "❌ Migration failed!"
    exit 1
fi

# Optional: Run database check
echo "🔍 Running database structure check..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f ../check_database.sql

echo "🎉 Database migration process completed!"







