# Local Games - Architecture Documentation

## 🏗️ System Architecture

### Overview

Local Games is a real-time multiplayer gaming platform designed with scalability, performance, and maintainability in mind. The system follows a layered architecture pattern with clear separation of concerns.

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Client Layer                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │   Desktop   │  │   Tablet    │  │   Mobile    │          │
│  │   Browser   │  │   Browser   │  │   Browser   │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Presentation Layer                          │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              SvelteKit Frontend                        │    │
│  │  • Responsive UI Components                           │    │
│  │  • Real-time WebSocket Client                         │    │
│  │  • State Management                                   │    │
│  │  • Routing & Navigation                               │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway Layer                         │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              Gin HTTP Server                           │    │
│  │  • REST API Endpoints                                 │    │
│  │  • WebSocket Handlers                                 │    │
│  │  • Middleware (CORS, Auth, Logging)                  │    │
│  │  • Request/Response Validation                        │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Business Logic Layer                       │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                Services Layer                           │    │
│  │  • Game Service (Game Logic)                          │    │
│  │  • Player Service (Player Management)                 │    │
│  │  • Socket Service (Real-time Communication)           │    │
│  │  • Authentication Service                              │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Data Access Layer                         │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              Data Layer                                │    │
│  │  • PostgreSQL (Persistent Data)                       │    │
│  │  • Redis (Caching & Sessions)                         │    │
│  │  • File System (Static Assets)                        │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

## 🎯 Design Principles

### 1. Separation of Concerns
- **Frontend**: UI/UX, client-side state management, real-time communication
- **Backend**: Business logic, data persistence, API management
- **Database**: Data storage and retrieval
- **Cache**: Session management and performance optimization

### 2. Scalability
- **Horizontal Scaling**: Stateless backend services
- **Database Scaling**: Read replicas and connection pooling
- **Cache Scaling**: Redis cluster for high availability
- **Load Balancing**: Multiple backend instances

### 3. Performance
- **Caching Strategy**: Redis for sessions and game state
- **Database Optimization**: Indexed queries and connection pooling
- **Real-time Communication**: WebSocket for low-latency updates
- **CDN**: Static asset delivery

### 4. Security
- **Input Validation**: Server-side validation for all inputs
- **Authentication**: Session-based authentication
- **CORS**: Proper cross-origin resource sharing
- **Rate Limiting**: API rate limiting to prevent abuse

## 📊 Data Flow

### 1. Player Registration Flow
```
Client → Frontend → API Gateway → Player Service → Database
  ↑                                                    ↓
Client ← Frontend ← API Gateway ← Player Service ← Database
```

### 2. Game Creation Flow
```
Client → Frontend → API Gateway → Game Service → Database
  ↑                                                    ↓
Client ← Frontend ← API Gateway ← Game Service ← Database
```

### 3. Real-time Game Flow
```
Client ↔ WebSocket ↔ Socket Handler ↔ Game Service ↔ Redis
  ↑                                                    ↓
Client ↔ WebSocket ↔ Socket Handler ↔ Game Service ↔ Database
```

## 🗄️ Database Design

### Entity Relationship Diagram

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Players   │    │    Games    │    │ GamePlayers │
│             │    │             │    │             │
│ • id        │◄──►│ • id        │◄──►│ • id        │
│ • name      │    │ • name      │    │ • game_id   │
│ • created_at│    │ • type      │    │ • player_id │
│ • updated_at│    │ • status    │    │ • role      │
└─────────────┘    │ • settings  │    │ • score     │
                   │ • created_at│    │ • joined_at │
                   └─────────────┘    └─────────────┘
                            │
                            ▼
                   ┌─────────────┐    ┌─────────────┐
                   │ GameRounds  │    │    Votes    │
                   │             │    │             │
                   │ • id        │◄──►│ • id        │
                   │ • game_id   │    │ • round_id  │
                   │ • round_num │    │ • voter_id  │
                   │ • category  │    │ • voted_for │
                   │ • word      │    │ • created_at│
                   │ • impostor  │    └─────────────┘
                   │ • status    │
                   │ • started_at│
                   │ • ended_at  │
                   └─────────────┘
```

### Database Schema

#### Players Table
```sql
CREATE TABLE players (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

#### Games Table
```sql
CREATE TABLE games (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'waiting',
    max_players INTEGER NOT NULL,
    min_players INTEGER NOT NULL,
    created_by VARCHAR(36) NOT NULL,
    settings JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

#### GamePlayers Table
```sql
CREATE TABLE game_players (
    id VARCHAR(36) PRIMARY KEY,
    game_id VARCHAR(36) NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    player_id VARCHAR(36) NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'player',
    score INTEGER NOT NULL DEFAULT 0,
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(game_id, player_id)
);
```

#### GameRounds Table
```sql
CREATE TABLE game_rounds (
    id VARCHAR(36) PRIMARY KEY,
    game_id VARCHAR(36) NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    round_number INTEGER NOT NULL,
    category VARCHAR(100) NOT NULL,
    word VARCHAR(100) NOT NULL,
    impostor_id VARCHAR(36) REFERENCES players(id),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMP
);
```

#### Votes Table
```sql
CREATE TABLE votes (
    id VARCHAR(36) PRIMARY KEY,
    round_id VARCHAR(36) NOT NULL REFERENCES game_rounds(id) ON DELETE CASCADE,
    voter_id VARCHAR(36) NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    voted_for_id VARCHAR(36) NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(round_id, voter_id)
);
```

## 🔄 Real-time Communication

### WebSocket Message Types

#### Client to Server
- `join_game` - Join a specific game
- `leave_game` - Leave the current game
- `start_game` - Start the game (host only)
- `end_game` - End the game (host only)
- `get_game_state` - Request current game state
- `get_player_word` - Get player's word for current round
- `vote` - Submit a vote during voting phase

#### Server to Client
- `player_joined` - Notification when player joins
- `player_left` - Notification when player leaves
- `game_started` - Notification when game starts
- `game_ended` - Notification when game ends
- `round_started` - Notification when new round starts
- `round_ended` - Notification when round ends
- `vote_received` - Notification when vote is received
- `vote_results` - Final voting results
- `game_state_update` - Real-time game state updates
- `player_word` - Player's word for current round
- `success` - Success message
- `error` - Error message

### WebSocket Connection Flow

```
1. Client connects to /ws?player_id=:id
2. Server validates player_id
3. Server registers client connection
4. Client can send/receive messages
5. Server broadcasts updates to all players in game
6. Client disconnects when leaving or game ends
```

## 🎮 Game Logic

### Impostor Game Flow

#### 1. Game Setup
```go
// Create game with settings
game := NewGame(name, "impostor", createdBy, maxPlayers, minPlayers, settings)

// Add creator to game
gamePlayer := NewGamePlayer(game.ID, createdBy, "player")
```

#### 2. Game Start
```go
// Validate minimum players
if len(players) < game.MinPlayers {
    return error
}

// Initialize game state
game.Status = "active"

// Start first round
impostorGame := NewImpostorGame(game, players)
impostorGame.StartGame()
```

#### 3. Round Management
```go
// Select random category
category := selectRandomCategory(game.Settings.Categories)

// Select random word
word := selectRandomWord(category)

// Select random impostor
impostor := selectRandomImpostor(players)

// Create round
round := NewGameRound(game.ID, roundNumber, category.Name, word, impostor.ID)
```

#### 4. Word Distribution
```go
// Regular players get the word
if playerID != round.ImpostorID {
    return round.Word
}

// Impostor gets no word
return ""
```

#### 5. Voting System
```go
// Collect votes
votes := collectVotes(roundID)

// Count votes
voteCounts := countVotes(votes)

// Determine result
eliminatedPlayer := getMostVotedPlayer(voteCounts)

// Check if impostor was caught
if eliminatedPlayer.ID == round.ImpostorID {
    // Impostor caught - players win
    playersWin = true
} else {
    // Wrong player eliminated - impostor wins
    impostorWins = true
}
```

## 🔧 API Design

### RESTful Endpoints

#### Players
```
POST   /api/v1/players           # Create player
GET    /api/v1/players/:id       # Get player
GET    /api/v1/players/:id/session # Get session
GET    /api/v1/players/:id/games # Get player's games
GET    /api/v1/players/:id/stats # Get player stats
```

#### Games
```
POST   /api/v1/games             # Create game
GET    /api/v1/games             # List games
GET    /api/v1/games/:id         # Get game
POST   /api/v1/games/:id/join    # Join game
POST   /api/v1/games/:id/leave   # Leave game
POST   /api/v1/games/:id/start   # Start game
POST   /api/v1/games/:id/end     # End game
GET    /api/v1/games/:id/state   # Get game state
```

### Request/Response Examples

#### Create Player
```http
POST /api/v1/players
Content-Type: application/json

{
  "name": "John Doe"
}
```

Response:
```json
{
  "message": "Player created successfully",
  "player": {
    "id": "uuid-here",
    "name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Create Game
```http
POST /api/v1/games
Content-Type: application/json
X-Player-ID: player-uuid

{
  "name": "Friday Night Game",
  "type": "impostor",
  "max_players": 8,
  "min_players": 3,
  "settings": {
    "categories": ["football_players", "movies"],
    "rounds": 3,
    "time_per_round": 60,
    "voting_enabled": true,
    "auto_start": false
  }
}
```

## 🚀 Performance Considerations

### 1. Database Optimization
- **Indexes**: Proper indexing on frequently queried columns
- **Connection Pooling**: Efficient database connection management
- **Query Optimization**: Optimized SQL queries with proper joins

### 2. Caching Strategy
- **Redis**: Session storage and game state caching
- **TTL**: Automatic expiration for temporary data
- **Pub/Sub**: Real-time message broadcasting

### 3. WebSocket Optimization
- **Connection Pooling**: Efficient WebSocket connection management
- **Message Batching**: Batch updates when possible
- **Heartbeat**: Keep connections alive with ping/pong

### 4. Frontend Optimization
- **Code Splitting**: Lazy loading of components
- **Image Optimization**: Compressed images and lazy loading
- **Bundle Optimization**: Tree shaking and minification

## 🔒 Security Considerations

### 1. Input Validation
- **Server-side Validation**: All inputs validated on server
- **Type Safety**: Strong typing with TypeScript
- **SQL Injection Prevention**: Parameterized queries

### 2. Authentication & Authorization
- **Session Management**: Secure session handling
- **Player Identification**: Unique player IDs
- **Game Access Control**: Validate player permissions

### 3. CORS Configuration
```go
config := cors.DefaultConfig()
config.AllowAllOrigins = true
config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Player-ID"}
```

### 4. Rate Limiting
- **API Rate Limiting**: Prevent abuse of API endpoints
- **WebSocket Rate Limiting**: Limit message frequency
- **IP-based Limiting**: Prevent spam from single IPs

## 📈 Scalability Patterns

### 1. Horizontal Scaling
- **Stateless Backend**: Multiple backend instances
- **Load Balancing**: Distribute traffic across instances
- **Database Scaling**: Read replicas for read-heavy operations

### 2. Microservices Architecture (Future)
```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  Game API   │  │ Player API  │  │ Auth API    │
└─────────────┘  └─────────────┘  └─────────────┘
       │                │                │
       └────────────────┼────────────────┘
                        ▼
              ┌─────────────────┐
              │  API Gateway    │
              └─────────────────┘
```

### 3. Event-Driven Architecture
- **Redis Pub/Sub**: Real-time event broadcasting
- **Event Sourcing**: Track all game events
- **CQRS**: Separate read and write operations

## 🧪 Testing Strategy

### 1. Unit Testing
- **Backend**: Go testing package for business logic
- **Frontend**: Jest for component testing
- **Database**: Integration tests for data layer

### 2. Integration Testing
- **API Testing**: Test all REST endpoints
- **WebSocket Testing**: Test real-time communication
- **End-to-End Testing**: Full user journey testing

### 3. Performance Testing
- **Load Testing**: Simulate multiple concurrent users
- **Stress Testing**: Test system limits
- **WebSocket Testing**: Test real-time performance

## 📊 Monitoring & Observability

### 1. Logging
- **Structured Logging**: JSON format for easy parsing
- **Log Levels**: Debug, Info, Warn, Error
- **Context**: Include request IDs and user context

### 2. Metrics
- **Application Metrics**: Response times, error rates
- **Business Metrics**: Active games, player counts
- **Infrastructure Metrics**: CPU, memory, disk usage

### 3. Tracing
- **Distributed Tracing**: Track requests across services
- **Performance Profiling**: Identify bottlenecks
- **Error Tracking**: Monitor and alert on errors

## 🔮 Future Enhancements

### 1. Technical Improvements
- **GraphQL**: More flexible API queries
- **WebRTC**: Peer-to-peer communication
- **Service Mesh**: Advanced service communication
- **Kubernetes**: Container orchestration

### 2. Feature Enhancements
- **Voice Chat**: Real-time voice communication
- **Video Chat**: Face-to-face gaming
- **AI Integration**: Smart game moderation
- **Mobile Apps**: Native mobile applications

### 3. Game Types
- **Word Association**: Chain word associations
- **Trivia**: Knowledge-based games
- **Drawing Games**: Collaborative drawing
- **Puzzle Games**: Logic and strategy games

---

This architecture provides a solid foundation for a scalable, maintainable, and performant multiplayer gaming platform. The modular design allows for easy extension and modification as requirements evolve. 