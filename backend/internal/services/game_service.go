package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"local-games/internal/database"
	"local-games/internal/games"
	"local-games/internal/models"
	"local-games/internal/redis"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// GameService handles game-related business logic
type GameService struct {
	db    *sqlx.DB
	redis *redis.Client
}

// NewGameService creates a new game service
func NewGameService(db *sqlx.DB, redis *redis.Client) *GameService {
	return &GameService{
		db:    db,
		redis: redis,
	}
}

// CreateGame creates a new game
func (s *GameService) CreateGame(req models.CreateGameRequest, createdBy string) (*models.Game, error) {
	// Validate game settings if it's an impostor game
	if req.Type == "impostor" {
		impostorGame := games.NewImpostorGame(nil, nil)
		if err := impostorGame.ValidateGameSettings(req.Settings); err != nil {
			return nil, err
		}
	}

	game := models.NewGame(req.Name, req.Type, createdBy, req.MaxPlayers, req.MinPlayers, req.Settings)

	// Insert into database
	query := `
		INSERT INTO games (id, name, type, status, max_players, min_players, created_by, settings, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	settingsJSON, err := json.Marshal(req.Settings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal settings: %w", err)
	}

	_, err = s.db.Exec(query,
		game.ID, game.Name, game.Type, game.Status, game.MaxPlayers, game.MinPlayers,
		game.CreatedBy, settingsJSON, game.CreatedAt, game.UpdatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	// Add creator to game
	gamePlayer := models.NewGamePlayer(game.ID, createdBy, "player")
	if err := s.addPlayerToGame(gamePlayer); err != nil {
		return nil, fmt.Errorf("failed to add creator to game: %w", err)
	}

	return game, nil
}

// GetGame retrieves a game by ID
func (s *GameService) GetGame(gameID string) (*models.Game, error) {
	var game models.Game
	query := `SELECT * FROM games WHERE id = ?`
	
	err := s.db.Get(&game, query, gameID)
	if err != nil {
		return nil, fmt.Errorf("game not found: %w", err)
	}

	return &game, nil
}

// ListGames retrieves all available games
func (s *GameService) ListGames() ([]models.Game, error) {
	var games []models.Game
	query := `SELECT * FROM games WHERE status = 'waiting' ORDER BY created_at DESC`
	
	err := s.db.Select(&games, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}

	return games, nil
}

// JoinGame adds a player to a game
func (s *GameService) JoinGame(gameID, playerID string) error {
	// Check if game exists and is joinable
	game, err := s.GetGame(gameID)
	if err != nil {
		return err
	}

	if game.Status != "waiting" {
		return fmt.Errorf("game is not joinable")
	}

	// Check if player is already in the game
	exists, err := s.isPlayerInGame(gameID, playerID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("player is already in the game")
	}

	// Check if game is full
	playerCount, err := s.getGamePlayerCount(gameID)
	if err != nil {
		return err
	}
	if playerCount >= game.MaxPlayers {
		return fmt.Errorf("game is full")
	}

	// Add player to game
	gamePlayer := models.NewGamePlayer(gameID, playerID, "player")
	return s.addPlayerToGame(gamePlayer)
}

// LeaveGame removes a player from a game
func (s *GameService) LeaveGame(gameID, playerID string) error {
	// Remove from database
	query := `DELETE FROM game_players WHERE game_id = ? AND player_id = ?`
	_, err := s.db.Exec(query, gameID, playerID)
	if err != nil {
		return fmt.Errorf("failed to leave game: %w", err)
	}

	// Remove from Redis
	return redis.RemovePlayerFromGame(gameID, playerID)
}

// StartGame starts a game
func (s *GameService) StartGame(gameID string) error {
	game, err := s.GetGame(gameID)
	if err != nil {
		return err
	}

	if game.Status != "waiting" {
		return fmt.Errorf("game cannot be started")
	}

	// Get game players
	players, err := s.getGamePlayers(gameID)
	if err != nil {
		return err
	}

	if len(players) < game.MinPlayers {
		return fmt.Errorf("not enough players to start game")
	}

	// Update game status
	query := `UPDATE games SET status = 'active', updated_at = ? WHERE id = ?`
	_, err = s.db.Exec(query, time.Now(), gameID)
	if err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}

	// Initialize game logic based on type
	if game.Type == "impostor" {
		impostorGame := games.NewImpostorGame(game, players)
		if err := impostorGame.StartGame(); err != nil {
			return fmt.Errorf("failed to start impostor game: %w", err)
		}
		
		// Store game state in Redis
		gameState := impostorGame.GetGameState()
		if err := redis.SetGameState(gameID, gameState); err != nil {
			log.Printf("Warning: failed to store game state in Redis: %v", err)
		}
	}

	return nil
}

// EndGame ends a game
func (s *GameService) EndGame(gameID string) error {
	query := `UPDATE games SET status = 'finished', updated_at = ? WHERE id = ?`
	_, err := s.db.Exec(query, time.Now(), gameID)
	if err != nil {
		return fmt.Errorf("failed to end game: %w", err)
	}

	return nil
}

// GetGameState retrieves the current state of a game
func (s *GameService) GetGameState(gameID string) (*models.GameState, error) {
	game, err := s.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	players, err := s.getGamePlayers(gameID)
	if err != nil {
		return nil, err
	}

	// Try to get from Redis first
	if gameStateStr, err := redis.GetGameState(gameID); err == nil {
		var gameState models.GameState
		if err := json.Unmarshal([]byte(gameStateStr), &gameState); err == nil {
			return &gameState, nil
		}
	}

	// Fallback to database
	gameState := models.GameState{
		Game:    game,
		Players: players,
	}

	return &gameState, nil
}

// GetPlayerWord returns the word for a specific player in an impostor game
func (s *GameService) GetPlayerWord(gameID, playerID string) (string, error) {
	game, err := s.GetGame(gameID)
	if err != nil {
		return "", err
	}

	if game.Type != "impostor" {
		return "", fmt.Errorf("game is not an impostor game")
	}

	players, err := s.getGamePlayers(gameID)
	if err != nil {
		return "", err
	}

	impostorGame := games.NewImpostorGame(game, players)
	return impostorGame.GetPlayerWord(playerID), nil
}

// Helper methods
func (s *GameService) addPlayerToGame(gamePlayer models.GamePlayer) error {
	query := `
		INSERT INTO game_players (id, game_id, player_id, role, score, joined_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	_, err := s.db.Exec(query,
		gamePlayer.ID, gamePlayer.GameID, gamePlayer.PlayerID,
		gamePlayer.Role, gamePlayer.Score, gamePlayer.JoinedAt)
	
	if err != nil {
		return fmt.Errorf("failed to add player to game: %w", err)
	}

	// Add to Redis
	return redis.AddPlayerToGame(gamePlayer.GameID, gamePlayer.PlayerID)
}

func (s *GameService) isPlayerInGame(gameID, playerID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM game_players WHERE game_id = ? AND player_id = ?`
	
	err := s.db.Get(&count, query, gameID, playerID)
	if err != nil {
		return false, fmt.Errorf("failed to check player in game: %w", err)
	}

	return count > 0, nil
}

func (s *GameService) getGamePlayerCount(gameID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM game_players WHERE game_id = ?`
	
	err := s.db.Get(&count, query, gameID)
	if err != nil {
		return 0, fmt.Errorf("failed to get game player count: %w", err)
	}

	return count, nil
}

func (s *GameService) getGamePlayers(gameID string) ([]models.GamePlayer, error) {
	var players []models.GamePlayer
	query := `SELECT * FROM game_players WHERE game_id = ? ORDER BY joined_at`
	
	err := s.db.Select(&players, query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game players: %w", err)
	}

	return players, nil
} 