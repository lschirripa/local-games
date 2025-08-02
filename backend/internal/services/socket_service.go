package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"local-games/internal/models"
	"local-games/internal/redis"

	"github.com/redis/go-redis/v9"
)

// SocketService handles WebSocket-related operations and Redis pub/sub
type SocketService struct {
	redis *redis.Client
}

// NewSocketService creates a new socket service
func NewSocketService(redis *redis.Client) *SocketService {
	return &SocketService{
		redis: redis,
	}
}

// PublishGameEvent publishes a game event to Redis
func (s *SocketService) PublishGameEvent(gameID string, eventType string, payload interface{}) error {
	message := models.WebSocketMessage{
		Type:    eventType,
		Payload: payload,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	channel := fmt.Sprintf("game:%s", gameID)
	return redis.PublishMessage(channel, string(messageBytes))
}

// SubscribeToGame subscribes to game events
func (s *SocketService) SubscribeToGame(gameID string) *redis.PubSub {
	channel := fmt.Sprintf("game:%s", gameID)
	return redis.SubscribeToChannel(channel)
}

// BroadcastToGame broadcasts a message to all players in a game
func (s *SocketService) BroadcastToGame(gameID string, message models.WebSocketMessage) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	channel := fmt.Sprintf("game:%s", gameID)
	return redis.PublishMessage(channel, string(messageBytes))
}

// SendToPlayer sends a message to a specific player
func (s *SocketService) SendToPlayer(playerID string, message models.WebSocketMessage) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	channel := fmt.Sprintf("player:%s", playerID)
	return redis.PublishMessage(channel, string(messageBytes))
}

// SubscribeToPlayer subscribes to player-specific events
func (s *SocketService) SubscribeToPlayer(playerID string) *redis.PubSub {
	channel := fmt.Sprintf("player:%s", playerID)
	return redis.SubscribeToChannel(channel)
}

// NotifyGameStateChange notifies all players in a game about state changes
func (s *SocketService) NotifyGameStateChange(gameID string, gameState *models.GameState) error {
	message := models.WebSocketMessage{
		Type:    "game_state_update",
		Payload: gameState,
	}

	return s.BroadcastToGame(gameID, message)
}

// NotifyPlayerJoined notifies all players when someone joins
func (s *SocketService) NotifyPlayerJoined(gameID, playerID, playerName string) error {
	message := models.WebSocketMessage{
		Type: "player_joined",
		Payload: map[string]interface{}{
			"player_id":   playerID,
			"player_name": playerName,
			"game_id":     gameID,
		},
	}

	return s.BroadcastToGame(gameID, message)
}

// NotifyPlayerLeft notifies all players when someone leaves
func (s *SocketService) NotifyPlayerLeft(gameID, playerID, playerName string) error {
	message := models.WebSocketMessage{
		Type: "player_left",
		Payload: map[string]interface{}{
			"player_id":   playerID,
			"player_name": playerName,
			"game_id":     gameID,
		},
	}

	return s.BroadcastToGame(gameID, message)
}

// NotifyGameStarted notifies all players when a game starts
func (s *SocketService) NotifyGameStarted(gameID string) error {
	message := models.WebSocketMessage{
		Type: "game_started",
		Payload: map[string]interface{}{
			"game_id": gameID,
		},
	}

	return s.BroadcastToGame(gameID, message)
}

// NotifyGameEnded notifies all players when a game ends
func (s *SocketService) NotifyGameEnded(gameID string, results map[string]interface{}) error {
	message := models.WebSocketMessage{
		Type: "game_ended",
		Payload: map[string]interface{}{
			"game_id": gameID,
			"results": results,
		},
	}

	return s.BroadcastToGame(gameID, message)
}

// NotifyRoundStart notifies all players when a new round starts
func (s *SocketService) NotifyRoundStart(gameID string, round *models.GameRound) error {
	message := models.WebSocketMessage{
		Type: "round_started",
		Payload: map[string]interface{}{
			"game_id": gameID,
			"round":   round,
		},
	}

	return s.BroadcastToGame(gameID, message)
}

// NotifyRoundEnd notifies all players when a round ends
func (s *SocketService) NotifyRoundEnd(gameID string, round *models.GameRound, results map[string]interface{}) error {
	message := models.WebSocketMessage{
		Type: "round_ended",
		Payload: map[string]interface{}{
			"game_id": gameID,
			"round":   round,
			"results": results,
		},
	}

	return s.BroadcastToGame(gameID, message)
}

// NotifyVoteReceived notifies all players when a vote is received
func (s *SocketService) NotifyVoteReceived(gameID, roundID, voterID, votedForID string) error {
	message := models.WebSocketMessage{
		Type: "vote_received",
		Payload: map[string]interface{}{
			"game_id":     gameID,
			"round_id":    roundID,
			"voter_id":    voterID,
			"voted_for_id": votedForID,
		},
	}

	return s.BroadcastToGame(gameID, message)
}

// NotifyVoteResults notifies all players of vote results
func (s *SocketService) NotifyVoteResults(gameID, roundID string, results map[string]interface{}) error {
	message := models.WebSocketMessage{
		Type: "vote_results",
		Payload: map[string]interface{}{
			"game_id":  gameID,
			"round_id": roundID,
			"results":  results,
		},
	}

	return s.BroadcastToGame(gameID, message)
}

// StartMessageListener starts a background goroutine to listen for Redis messages
func (s *SocketService) StartMessageListener(ctx context.Context, gameID string, messageHandler func(models.WebSocketMessage)) {
	sub := s.SubscribeToGame(gameID)
	defer sub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-sub.Channel():
			var wsMessage models.WebSocketMessage
			if err := json.Unmarshal([]byte(msg.Payload), &wsMessage); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			messageHandler(wsMessage)
		}
	}
}

// GetActivePlayers gets the list of active players in a game
func (s *SocketService) GetActivePlayers(gameID string) ([]string, error) {
	return redis.GetGamePlayers(gameID)
}

// AddPlayerToGame adds a player to the active players list
func (s *SocketService) AddPlayerToGame(gameID, playerID string) error {
	return redis.AddPlayerToGame(gameID, playerID)
}

// RemovePlayerFromGame removes a player from the active players list
func (s *SocketService) RemovePlayerFromGame(gameID, playerID string) error {
	return redis.RemovePlayerFromGame(gameID, playerID)
} 