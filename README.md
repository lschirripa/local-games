# Local Games - Multiplayer Gaming Platform

A real-time multiplayer gaming platform built with **SvelteKit** (frontend) and **Golang** (backend). Players can create and join games, with the first implemented game being "Impostor" - a social deduction game where players must identify the impostor among them.

## 🏗️ Architecture Overview

### Technology Stack

**Frontend:**
- **SvelteKit** - Modern web framework for building reactive UIs
- **TypeScript** - Type-safe JavaScript
- **Tailwind CSS** - Utility-first CSS framework for responsive design
- **Socket.IO Client** - Real-time WebSocket communication

**Backend:**
- **Golang** - High-performance server language
- **Gin** - HTTP web framework
- **Gorilla WebSocket** - WebSocket implementation
- **Redis** - In-memory data store for sessions and caching
- **PostgreSQL** - Relational database for persistent data

**Infrastructure:**
- **Docker** - Containerization (optional)
- **WebRTC** - Peer-to-peer communication (future enhancement)

### System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   Database      │
│   (SvelteKit)   │◄──►│   (Golang)      │◄──►│   (PostgreSQL)  │
│                 │    │                 │    │                 │
│ • Responsive UI │    │ • REST API      │    │ • Game Data     │
│ • Real-time     │    │ • WebSocket     │    │ • Player Data   │
│ • Game Logic    │    │ • Game Engine   │    │ • Session Data  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   Redis Cache   │
                       │                 │
                       │ • Sessions      │
                       │ • Game State    │
                       │ • Pub/Sub       │
                       └─────────────────┘
```

## 🎮 Game: Impostor

The "Impostor" game is a social deduction game where:

1. **Setup**: Players join a game and one is randomly selected as the impostor
2. **Categories**: Games can be configured with different word categories (Football Players, Movies, Animals, etc.)
3. **Gameplay**: 
   - All players receive the same word from a category
   - The impostor receives no word (or a different word)
   - Players discuss and try to identify the impostor
   - Voting phase to eliminate suspected impostors
4. **Scoring**: Points awarded for correct identification or successful deception

### Game Flow

```
1. Create/Join Game
   ↓
2. Wait for Players (Min: 3, Max: 20)
   ↓
3. Start Game
   ↓
4. Round Begins
   ├── Category Selection
   ├── Word Distribution
   └── Impostor Selection
   ↓
5. Discussion Phase
   ↓
6. Voting Phase
   ↓
7. Results & Scoring
   ↓
8. Next Round or Game End
```

## 🚀 Getting Started

### Prerequisites

- **Go** 1.21 or higher
- **Node.js** 18 or higher
- **PostgreSQL** 13 or higher
- **Redis** 6 or higher

### Backend Setup

1. **Navigate to backend directory:**
   ```bash
   cd backend
   ```

2. **Install Go dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables:**
   ```bash
   cp env.example .env
   # Edit .env with your database and Redis credentials
   ```

4. **Initialize database:**
   ```sql
   CREATE DATABASE local_games;
   ```

5. **Run the backend:**
   ```bash
   go run cmd/api/main.go
   ```
   The backend will start on `http://localhost:8080`

### Frontend Setup

1. **Navigate to frontend directory:**
   ```bash
   cd front
   ```

2. **Install dependencies:**
   ```bash
   npm install
   ```

3. **Run the development server:**
   ```bash
   npm run dev
   ```
   The frontend will be available at `http://localhost:5173`

## 📁 Project Structure

```
local-games/
├── backend/
│   ├── cmd/api/main.go          # Application entry point
│   ├── internal/
│   │   ├── config/              # Configuration management
│   │   ├── database/            # Database connection & migrations
│   │   ├── games/               # Game logic (Impostor game)
│   │   ├── handlers/            # HTTP & WebSocket handlers
│   │   ├── middleware/          # HTTP middleware
│   │   ├── models/              # Data models & structs
│   │   ├── redis/               # Redis client & operations
│   │   ├── server/              # Server setup
│   │   └── services/            # Business logic layer
│   ├── go.mod                   # Go module dependencies
│   └── env.example              # Environment variables template
├── front/
│   ├── src/
│   │   ├── routes/              # SvelteKit pages
│   │   ├── lib/                 # Shared components & utilities
│   │   └── app.css              # Global styles (Tailwind)
│   ├── package.json             # Node.js dependencies
│   ├── tailwind.config.js       # Tailwind CSS configuration
│   └── postcss.config.js        # PostCSS configuration
└── README.md                    # This file
```

## 🔧 API Endpoints

### Players
- `POST /api/v1/players` - Create a new player
- `GET /api/v1/players/:id` - Get player information
- `GET /api/v1/players/:id/session` - Get player session
- `GET /api/v1/players/:id/games` - Get player's games
- `GET /api/v1/players/:id/stats` - Get player statistics

### Games
- `POST /api/v1/games` - Create a new game
- `GET /api/v1/games` - List available games
- `GET /api/v1/games/:id` - Get game details
- `POST /api/v1/games/:id/join` - Join a game
- `POST /api/v1/games/:id/leave` - Leave a game
- `POST /api/v1/games/:id/start` - Start a game
- `POST /api/v1/games/:id/end` - End a game
- `GET /api/v1/games/:id/state` - Get game state

### WebSocket
- `GET /ws?player_id=:id` - WebSocket connection for real-time updates

## 🎯 Features

### Core Features
- ✅ **Real-time multiplayer gaming**
- ✅ **Responsive design** (mobile & desktop)
- ✅ **Player authentication** (session-based)
- ✅ **Game lobby system**
- ✅ **Impostor game implementation**
- ✅ **WebSocket communication**
- ✅ **Redis caching & pub/sub**

### Game Features
- ✅ **Multiple word categories** (Football Players, Movies, Animals, etc.)
- ✅ **Configurable game settings** (rounds, time limits, player counts)
- ✅ **Real-time game state updates**
- ✅ **Voting system**
- ✅ **Score tracking**

### Technical Features
- ✅ **TypeScript** for type safety
- ✅ **Tailwind CSS** for styling
- ✅ **PostgreSQL** for data persistence
- ✅ **Redis** for caching and real-time features
- ✅ **WebSocket** for real-time communication
- ✅ **RESTful API** design
- ✅ **Error handling** and validation

## 🚀 Deployment

### Docker Deployment (Optional)

1. **Create Docker Compose file:**
   ```yaml
   version: '3.8'
   services:
     backend:
       build: ./backend
       ports:
         - "8080:8080"
       environment:
         - DB_HOST=postgres
         - REDIS_HOST=redis
       depends_on:
         - postgres
         - redis
     
     frontend:
       build: ./front
       ports:
         - "3000:3000"
       depends_on:
         - backend
     
     postgres:
       image: postgres:13
       environment:
         POSTGRES_DB: local_games
         POSTGRES_USER: postgres
         POSTGRES_PASSWORD: password
       volumes:
         - postgres_data:/var/lib/postgresql/data
     
     redis:
       image: redis:6-alpine
       ports:
         - "6379:6379"
   
   volumes:
     postgres_data:
   ```

2. **Run with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

## 🔮 Future Enhancements

### Planned Features
- 🔄 **Additional game types** (Word Association, Trivia, etc.)
- 🔄 **User accounts** with persistent profiles
- 🔄 **Friend system** and private games
- 🔄 **Achievement system** and leaderboards
- 🔄 **Custom word categories** creation
- 🔄 **Voice chat** integration
- 🔄 **Mobile app** (React Native/Flutter)

### Technical Improvements
- 🔄 **Microservices architecture** for scalability
- 🔄 **Kubernetes deployment** for production
- 🔄 **GraphQL API** for flexible data fetching
- 🔄 **WebRTC** for peer-to-peer communication
- 🔄 **Unit & integration tests** coverage
- 🔄 **CI/CD pipeline** automation

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/your-repo/local-games/issues) page
2. Create a new issue with detailed information
3. Join our [Discord](https://discord.gg/localgames) community

---

**Built with ❤️ using SvelteKit and Go**
