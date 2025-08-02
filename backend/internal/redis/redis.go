package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"local-games/internal/config"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

// NewClient creates a new Redis client
func NewClient() (*redis.Client, error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.AppConfig.RedisHost, config.AppConfig.RedisPort),
		DB:       0,
		Password: "", // Add password if needed
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connection established successfully")
	return client, nil
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	return client
}

// SetGameState stores the current game state in Redis
func SetGameState(gameID string, state interface{}) error {
	ctx := context.Background()
	key := fmt.Sprintf("game:%s:state", gameID)
	
	// Store with 1 hour expiration
	return client.Set(ctx, key, state, time.Hour).Err()
}

// GetGameState retrieves the current game state from Redis
func GetGameState(gameID string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("game:%s:state", gameID)
	
	return client.Get(ctx, key).Result()
}

// SetPlayerSession stores player session information
func SetPlayerSession(playerID string, sessionData interface{}) error {
	ctx := context.Background()
	key := fmt.Sprintf("player:%s:session", playerID)
	
	// Store with 24 hour expiration
	return client.Set(ctx, key, sessionData, 24*time.Hour).Err()
}

// GetPlayerSession retrieves player session information
func GetPlayerSession(playerID string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("player:%s:session", playerID)
	
	return client.Get(ctx, key).Result()
}

// AddPlayerToGame adds a player to a game's active players list
func AddPlayerToGame(gameID, playerID string) error {
	ctx := context.Background()
	key := fmt.Sprintf("game:%s:players", gameID)
	
	return client.SAdd(ctx, key, playerID).Err()
}

// RemovePlayerFromGame removes a player from a game's active players list
func RemovePlayerFromGame(gameID, playerID string) error {
	ctx := context.Background()
	key := fmt.Sprintf("game:%s:players", gameID)
	
	return client.SRem(ctx, key, playerID).Err()
}

// GetGamePlayers retrieves all active players in a game
func GetGamePlayers(gameID string) ([]string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("game:%s:players", gameID)
	
	return client.SMembers(ctx, key).Result()
}

// PublishMessage publishes a message to a specific channel
func PublishMessage(channel string, message interface{}) error {
	ctx := context.Background()
	return client.Publish(ctx, channel, message).Err()
}

// SubscribeToChannel subscribes to a Redis channel
func SubscribeToChannel(channel string) *redis.PubSub {
	return client.Subscribe(context.Background(), channel)
} 