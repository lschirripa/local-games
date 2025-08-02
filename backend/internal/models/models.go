package models

import (
	"time"

	"github.com/google/uuid"
)

// Player represents a user in the system
type Player struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Game represents a game session
type Game struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Type        string    `json:"type" db:"type"`
	Status      string    `json:"status" db:"status"` // waiting, active, finished
	MaxPlayers  int       `json:"max_players" db:"max_players"`
	MinPlayers  int       `json:"min_players" db:"min_players"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	Settings    GameSettings `json:"settings" db:"settings"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// GameSettings contains game-specific configuration
type GameSettings struct {
	Categories     []string `json:"categories"`
	Rounds         int      `json:"rounds"`
	TimePerRound   int      `json:"time_per_round"` // in seconds
	VotingEnabled  bool     `json:"voting_enabled"`
	AutoStart      bool     `json:"auto_start"`
}

// GamePlayer represents a player in a specific game
type GamePlayer struct {
	ID        string    `json:"id" db:"id"`
	GameID    string    `json:"game_id" db:"game_id"`
	PlayerID  string    `json:"player_id" db:"player_id"`
	Role      string    `json:"role" db:"role"` // player, impostor, spectator
	Score     int       `json:"score" db:"score"`
	JoinedAt  time.Time `json:"joined_at" db:"joined_at"`
}

// GameRound represents a single round in a game
type GameRound struct {
	ID          string    `json:"id" db:"id"`
	GameID      string    `json:"game_id" db:"game_id"`
	RoundNumber int       `json:"round_number" db:"round_number"`
	Category    string    `json:"category" db:"category"`
	Word        string    `json:"word" db:"word"`
	ImpostorID  string    `json:"impostor_id" db:"impostor_id"`
	Status      string    `json:"status" db:"status"` // active, voting, finished
	StartedAt   time.Time `json:"started_at" db:"started_at"`
	EndedAt     *time.Time `json:"ended_at" db:"ended_at"`
}

// Vote represents a player's vote in a round
type Vote struct {
	ID         string    `json:"id" db:"id"`
	RoundID    string    `json:"round_id" db:"round_id"`
	VoterID    string    `json:"voter_id" db:"voter_id"`
	VotedForID string    `json:"voted_for_id" db:"voted_for_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// CreatePlayerRequest represents the request to create a new player
type CreatePlayerRequest struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
}

// CreateGameRequest represents the request to create a new game
type CreateGameRequest struct {
	Name       string        `json:"name" binding:"required,min=2,max=100"`
	Type       string        `json:"type" binding:"required"`
	MaxPlayers int           `json:"max_players" binding:"required,min=2,max=20"`
	MinPlayers int           `json:"min_players" binding:"required,min=2"`
	Settings   GameSettings  `json:"settings"`
}

// JoinGameRequest represents the request to join a game
type JoinGameRequest struct {
	PlayerID string `json:"player_id" binding:"required"`
}

// GameState represents the current state of a game for real-time updates
type GameState struct {
	Game       *Game        `json:"game"`
	Players    []GamePlayer `json:"players"`
	CurrentRound *GameRound `json:"current_round,omitempty"`
	RoundHistory []GameRound `json:"round_history,omitempty"`
}

// WebSocketMessage represents a message sent through WebSocket
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	PlayerID string     `json:"player_id,omitempty"`
}

// NewPlayer creates a new player with a generated ID
func NewPlayer(name string) *Player {
	return &Player{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewGame creates a new game with a generated ID
func NewGame(name, gameType, createdBy string, maxPlayers, minPlayers int, settings GameSettings) *Game {
	return &Game{
		ID:         uuid.New().String(),
		Name:       name,
		Type:       gameType,
		Status:     "waiting",
		MaxPlayers: maxPlayers,
		MinPlayers: minPlayers,
		CreatedBy:  createdBy,
		Settings:   settings,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// NewGamePlayer creates a new game player relationship
func NewGamePlayer(gameID, playerID, role string) *GamePlayer {
	return &GamePlayer{
		ID:       uuid.New().String(),
		GameID:   gameID,
		PlayerID: playerID,
		Role:     role,
		Score:    0,
		JoinedAt: time.Now(),
	}
}

// NewGameRound creates a new game round
func NewGameRound(gameID string, roundNumber int, category, word, impostorID string) *GameRound {
	return &GameRound{
		ID:          uuid.New().String(),
		GameID:      gameID,
		RoundNumber: roundNumber,
		Category:    category,
		Word:        word,
		ImpostorID:  impostorID,
		Status:      "active",
		StartedAt:   time.Now(),
	}
} 