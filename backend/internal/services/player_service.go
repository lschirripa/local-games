package services

import (
	"fmt"
	"time"

	"local-games/internal/models"
	"local-games/internal/redis"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// PlayerService handles player-related business logic
type PlayerService struct {
	db    *sqlx.DB
	redis *redis.Client
}

// NewPlayerService creates a new player service
func NewPlayerService(db *sqlx.DB, redis *redis.Client) *PlayerService {
	return &PlayerService{
		db:    db,
		redis: redis,
	}
}

// CreatePlayer creates a new player
func (s *PlayerService) CreatePlayer(req models.CreatePlayerRequest) (*models.Player, error) {
	player := models.NewPlayer(req.Name)

	query := `
		INSERT INTO players (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`
	
	_, err := s.db.Exec(query, player.ID, player.Name, player.CreatedAt, player.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	// Store player session in Redis
	sessionData := map[string]interface{}{
		"player_id": player.ID,
		"name":      player.Name,
		"created_at": player.CreatedAt,
	}
	
	if err := redis.SetPlayerSession(player.ID, sessionData); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to store player session: %v\n", err)
	}

	return player, nil
}

// GetPlayer retrieves a player by ID
func (s *PlayerService) GetPlayer(playerID string) (*models.Player, error) {
	var player models.Player
	query := `SELECT * FROM players WHERE id = ?`
	
	err := s.db.Get(&player, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("player not found: %w", err)
	}

	return &player, nil
}

// GetPlayerSession retrieves player session data from Redis
func (s *PlayerService) GetPlayerSession(playerID string) (map[string]interface{}, error) {
	sessionStr, err := redis.GetPlayerSession(playerID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Parse session data (simplified for now)
	sessionData := map[string]interface{}{
		"player_id": playerID,
		"session":   sessionStr,
	}

	return sessionData, nil
}

// UpdatePlayer updates a player's information
func (s *PlayerService) UpdatePlayer(playerID string, name string) (*models.Player, error) {
	query := `UPDATE players SET name = ?, updated_at = ? WHERE id = ?`
	
	_, err := s.db.Exec(query, name, time.Now(), playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to update player: %w", err)
	}

	// Get updated player
	return s.GetPlayer(playerID)
}

// DeletePlayer deletes a player
func (s *PlayerService) DeletePlayer(playerID string) error {
	query := `DELETE FROM players WHERE id = ?`
	
	_, err := s.db.Exec(query, playerID)
	if err != nil {
		return fmt.Errorf("failed to delete player: %w", err)
	}

	// Remove from Redis
	// Note: In a real application, you'd want to handle this more carefully
	// and ensure the player is not in any active games

	return nil
}

// GetPlayerGames retrieves all games a player is participating in
func (s *PlayerService) GetPlayerGames(playerID string) ([]models.Game, error) {
	var games []models.Game
	query := `
		SELECT g.* FROM games g
		JOIN game_players gp ON g.id = gp.game_id
		WHERE gp.player_id = ?
		ORDER BY g.created_at DESC
	`
	
	err := s.db.Select(&games, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player games: %w", err)
	}

	return games, nil
}

// GetPlayerStats retrieves player statistics
func (s *PlayerService) GetPlayerStats(playerID string) (map[string]interface{}, error) {
	// Get total games played
	var totalGames int
	query := `SELECT COUNT(DISTINCT game_id) FROM game_players WHERE player_id = ?`
	err := s.db.Get(&totalGames, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total games: %w", err)
	}

	// Get games won (as impostor or as regular player)
	var gamesWon int
	query = `SELECT COUNT(DISTINCT gp.game_id) FROM game_players gp
		JOIN games g ON gp.game_id = g.id
		WHERE gp.player_id = ? AND g.status = 'finished'`
	err = s.db.Get(&gamesWon, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get games won: %w", err)
	}

	// Get total score
	var totalScore int
	query = `SELECT COALESCE(SUM(score), 0) FROM game_players WHERE player_id = ?`
	err = s.db.Get(&totalScore, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total score: %w", err)
	}

	stats := map[string]interface{}{
		"total_games": totalGames,
		"games_won":   gamesWon,
		"total_score": totalScore,
		"win_rate":    float64(gamesWon) / float64(totalGames) * 100,
	}

	return stats, nil
} 