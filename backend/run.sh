#!/bin/bash

# Farm Management Backend Runner Script

echo "🚀 Starting Farm Management Backend..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Creating from example..."
    cp env.example .env
    echo "📝 Please edit .env file with your database credentials"
    echo "   Then run this script again"
    exit 1
fi

# Check if go.mod exists
if [ ! -f go.mod ]; then
    echo "❌ go.mod not found. Please run 'go mod init' first"
    exit 1
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy

# Run the application
echo "🏃 Running application..."
go run main.go
