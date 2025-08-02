package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Player struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SessionID string    `json:"session_id" gorm:"uniqueIndex;not null"`
	Username  string    `json:"username" gorm:"size:50"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	LastSeen  time.Time `json:"last_seen" gorm:"autoUpdateTime"`
}

func (p *Player) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type PlayerResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}

func (p *Player) ToResponse() PlayerResponse {
	return PlayerResponse{
		ID:        p.ID.String(),
		Username:  p.Username,
		CreatedAt: p.CreatedAt,
		LastSeen:  p.LastSeen,
	}
} 