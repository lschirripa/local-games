#!/bin/bash

# Local Games Setup Script
echo "🎮 Setting up Local Games - Multiplayer Gaming Platform"
echo "========================================================"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

echo "✅ Docker and Docker Compose are installed"

# Create necessary directories
echo "📁 Creating project directories..."
mkdir -p frontend/src/lib/components
mkdir -p frontend/src/lib/stores
mkdir -p frontend/src/lib/services
mkdir -p frontend/src/lib/utils
mkdir -p backend/internal/api
mkdir -p backend/internal/game
mkdir -p backend/internal/models
mkdir -p backend/internal/services
mkdir -p backend/internal/websocket
mkdir -p docker/nginx

echo "✅ Directories created"

# Check if .env files exist, create if not
if [ ! -f "backend/.env" ]; then
    echo "📝 Creating backend .env file..."
    cat > backend/.env << EOF
DB_HOST=postgres
DB_PORT=5432
DB_NAME=localgames
DB_USER=localgames
DB_PASSWORD=localgames123
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=your-secret-key-here
EOF
    echo "✅ Backend .env created"
fi

if [ ! -f "frontend/.env" ]; then
    echo "📝 Creating frontend .env file..."
    cat > frontend/.env << EOF
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws
EOF
    echo "✅ Frontend .env created"
fi

# Start the services
echo "🚀 Starting services with Docker Compose..."
docker-compose up -d

echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are running
echo "🔍 Checking service status..."

# Check PostgreSQL
if docker-compose ps postgres | grep -q "Up"; then
    echo "✅ PostgreSQL is running"
else
    echo "❌ PostgreSQL failed to start"
fi

# Check Redis
if docker-compose ps redis | grep -q "Up"; then
    echo "✅ Redis is running"
else
    echo "❌ Redis failed to start"
fi

# Check Backend
if docker-compose ps backend | grep -q "Up"; then
    echo "✅ Backend is running"
else
    echo "❌ Backend failed to start"
fi

# Check Frontend
if docker-compose ps frontend | grep -q "Up"; then
    echo "✅ Frontend is running"
else
    echo "❌ Frontend failed to start"
fi

echo ""
echo "🎉 Setup complete!"
echo "=================="
echo "🌐 Frontend: http://localhost:3000"
echo "🔧 Backend API: http://localhost:8080"
echo "🗄️  Database: localhost:5432"
echo "⚡ Redis: localhost:6379"
echo ""
echo "📖 For more information, check the README.md file"
echo "🐛 For issues, check the logs with: docker-compose logs" 