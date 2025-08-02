package main

import (
	"log"
	"os"

	"local-games/backend/internal/api"
	"local-games/backend/internal/config"
	"local-games/backend/internal/database"
	"local-games/backend/internal/websocket"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.Database)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Redis
	redis, err := database.InitializeRedis(cfg.Redis)
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize and start server
	server := api.NewServer(cfg, db, redis, hub)
	
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 