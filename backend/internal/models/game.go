package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameStatus string

const (
	GameStatusWaiting GameStatus = "waiting"
	GameStatusActive  GameStatus = "active"
	GameStatusEnded   GameStatus = "ended"
)

type GameType string

const (
	GameTypeImpostor GameType = "impostor"
	// Add more game types here
)

type Game struct {
	ID             uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GameType       GameType        `json:"game_type" gorm:"size:50;not null"`
	RoomCode       string          `json:"room_code" gorm:"uniqueIndex;size:10;not null"`
	Status         GameStatus      `json:"status" gorm:"size:20;default:'waiting'"`
	MaxPlayers     int             `json:"max_players" gorm:"default:8"`
	CurrentPlayers int             `json:"current_players" gorm:"default:0"`
	GameConfig     json.RawMessage `json:"game_config" gorm:"type:jsonb"`
	CreatedAt      time.Time       `json:"created_at" gorm:"autoCreateTime"`
	StartedAt      *time.Time      `json:"started_at"`
	EndedAt        *time.Time      `json:"ended_at"`
	
	// Relationships
	Players []GamePlayer `json:"players" gorm:"foreignKey:GameID"`
}

func (g *Game) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

type GameResponse struct {
	ID             string          `json:"id"`
	GameType       string          `json:"game_type"`
	RoomCode       string          `json:"room_code"`
	Status         string          `json:"status"`
	MaxPlayers     int             `json:"max_players"`
	CurrentPlayers int             `json:"current_players"`
	GameConfig     json.RawMessage `json:"game_config"`
	CreatedAt      time.Time       `json:"created_at"`
	StartedAt      *time.Time      `json:"started_at"`
	EndedAt        *time.Time      `json:"ended_at"`
	Players        []GamePlayerResponse `json:"players"`
}

func (g *Game) ToResponse() GameResponse {
	players := make([]GamePlayerResponse, len(g.Players))
	for i, player := range g.Players {
		players[i] = player.ToResponse()
	}

	return GameResponse{
		ID:             g.ID.String(),
		GameType:       string(g.GameType),
		RoomCode:       g.RoomCode,
		Status:         string(g.Status),
		MaxPlayers:     g.MaxPlayers,
		CurrentPlayers: g.CurrentPlayers,
		GameConfig:     g.GameConfig,
		CreatedAt:      g.CreatedAt,
		StartedAt:      g.StartedAt,
		EndedAt:        g.EndedAt,
		Players:        players,
	}
}

type GamePlayer struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GameID    uuid.UUID       `json:"game_id" gorm:"type:uuid;not null"`
	PlayerID  uuid.UUID       `json:"player_id" gorm:"type:uuid;not null"`
	Role      string          `json:"role" gorm:"size:50"`
	GameData  json.RawMessage `json:"game_data" gorm:"type:jsonb"`
	JoinedAt  time.Time       `json:"joined_at" gorm:"autoCreateTime"`
	
	// Relationships
	Player Player `json:"player" gorm:"foreignKey:PlayerID"`
}

func (gp *GamePlayer) BeforeCreate(tx *gorm.DB) error {
	if gp.ID == uuid.Nil {
		gp.ID = uuid.New()
	}
	return nil
}

type GamePlayerResponse struct {
	ID       string          `json:"id"`
	GameID   string          `json:"game_id"`
	PlayerID string          `json:"player_id"`
	Role     string          `json:"role"`
	GameData json.RawMessage `json:"game_data"`
	JoinedAt time.Time       `json:"joined_at"`
	Player   PlayerResponse  `json:"player"`
}

func (gp *GamePlayer) ToResponse() GamePlayerResponse {
	return GamePlayerResponse{
		ID:       gp.ID.String(),
		GameID:   gp.GameID.String(),
		PlayerID: gp.PlayerID.String(),
		Role:     gp.Role,
		GameData: gp.GameData,
		JoinedAt: gp.JoinedAt,
		Player:   gp.Player.ToResponse(),
	}
} 