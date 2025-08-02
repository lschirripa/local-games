# Local Games - Multiplayer Gaming Platform

A modern, responsive multiplayer gaming platform built with **Svelte** (frontend) and **Go** (backend) that allows players to create and join various games in real-time. The platform features the "Impostor" game as the first implementation, with room for expansion to many other game types.

## 🎮 Features

- **Real-time Multiplayer**: WebSocket-based real-time communication
- **Responsive Design**: Works seamlessly on desktop, tablet, and mobile
- **Cross-platform**: Accessible from any modern browser
- **Game Types**: Extensible system for multiple game types
- **Player Management**: Unique player identification and session management
- **Room System**: Create and join games with room codes
- **Impostor Game**: Find the impostor among your friends!

## 🏗️ Architecture

### Technology Stack

**Frontend (Svelte)**
- Svelte 5 - Modern reactive framework
- SvelteKit - Full-stack framework
- TypeScript - Type safety
- Tailwind CSS - Utility-first CSS framework
- Socket.IO Client - Real-time communication
- Vite - Fast build tool

**Backend (Go)**
- Go 1.21+ - High-performance server language
- Gin - HTTP web framework
- Gorilla WebSocket - WebSocket implementation
- GORM - ORM for database operations
- PostgreSQL - Primary database
- Redis - Caching and session management

**DevOps & Infrastructure**
- Docker Compose - Local development environment
- Nginx - Reverse proxy and load balancer
- GitHub Actions - CI/CD pipeline (planned)

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for local development)
- Go 1.21+ (for local development)

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd local-games
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432
   - Redis: localhost:6379

### Local Development

1. **Backend Setup**
   ```bash
   cd backend
   go mod download
   go run cmd/server/main.go
   ```

2. **Frontend Setup**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

3. **Database Setup**
   ```bash
   # Start PostgreSQL and Redis
   docker-compose up postgres redis -d
   ```

## 📁 Project Structure

```
local-games/
├── frontend/                 # Svelte application
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/   # Reusable UI components
│   │   │   ├── stores/       # Svelte stores
│   │   │   ├── services/     # API and WebSocket services
│   │   │   └── utils/        # Utility functions
│   │   ├── routes/           # SvelteKit routes
│   │   └── app.html
│   ├── static/               # Static assets
│   └── package.json
├── backend/                  # Go application
│   ├── cmd/
│   │   └── server/           # Main application entry
│   ├── internal/
│   │   ├── api/             # HTTP handlers
│   │   ├── game/            # Game logic
│   │   ├── models/          # Data models
│   │   ├── services/        # Business logic
│   │   └── websocket/       # WebSocket handlers
│   ├── pkg/                 # Public packages
│   └── go.mod
├── docker-compose.yml        # Development environment
├── docker/                  # Docker configurations
├── docs/                    # Documentation
└── scripts/                 # Build and deployment scripts
```

## 🎯 Game Types

### Impostor Game

The first implemented game type where:
- Players receive cards with the same answer (e.g., "football players")
- One player (the impostor) receives a different card
- Players must identify the impostor through discussion
- Configurable categories and game settings

### Future Game Types

- Word Association
- Drawing Challenge
- Trivia Battle
- And many more...

## 🔧 Configuration

### Environment Variables

**Backend (.env)**
```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=localgames
DB_USER=localgames
DB_PASSWORD=localgames123
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-secret-key-here
```

**Frontend (.env)**
```env
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080/ws
```

## 🛠️ Development

### API Endpoints

**Authentication**
- `POST /api/auth/session` - Create session
- `DELETE /api/auth/session` - Delete session
- `GET /api/auth/me` - Get current player

**Games**
- `GET /api/games` - List games
- `POST /api/games` - Create game
- `GET /api/games/{id}` - Get game
- `PUT /api/games/{id}` - Update game
- `DELETE /api/games/{id}` - Delete game
- `POST /api/games/{id}/join` - Join game
- `POST /api/games/{id}/leave` - Leave game

**WebSocket Events**
- `join_game` - Join game room
- `game_action` - Game-specific actions
- `leave_game` - Leave game room

### Database Schema

**Players Table**
```sql
CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    last_seen TIMESTAMP DEFAULT NOW()
);
```

**Games Table**
```sql
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_type VARCHAR(50) NOT NULL,
    room_code VARCHAR(10) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'waiting',
    max_players INTEGER DEFAULT 8,
    current_players INTEGER DEFAULT 0,
    game_config JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    started_at TIMESTAMP,
    ended_at TIMESTAMP
);
```

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm run test
```

## 🚀 Deployment

### Production Build

1. **Build Frontend**
   ```bash
   cd frontend
   npm run build
   ```

2. **Build Backend**
   ```bash
   cd backend
   go build -o bin/server cmd/server/main.go
   ```

3. **Docker Production**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

## 📊 Monitoring

- **Health Check**: `GET /health`
- **Metrics**: Prometheus endpoints (planned)
- **Logging**: Structured logging with correlation IDs
- **Error Tracking**: Comprehensive error handling

## 🔒 Security

- **HTTPS**: Enforced in production
- **CORS**: Configured for cross-origin requests
- **Input Validation**: Server-side validation
- **Rate Limiting**: API endpoint protection
- **SQL Injection**: Parameterized queries
- **XSS Protection**: Content encoding

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with modern web technologies
- Inspired by the need for accessible multiplayer gaming
- Community-driven development approach

## 📞 Contact

- Project Link: [https://github.com/yourusername/local-games](https://github.com/yourusername/local-games)
- Issues: [GitHub Issues](https://github.com/yourusername/local-games/issues)

---

**Note:** This project is under active development. Features and documentation may change frequently.