package main

import (
	"log"
	"os"

	"local-games/internal/config"
	"local-games/internal/database"
	"local-games/internal/handlers"
	"local-games/internal/middleware"
	"local-games/internal/redis"
	"local-games/internal/server"
	"local-games/internal/services"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Initialize Redis
	redisClient, err := redis.NewClient()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	defer redisClient.Close()

	// Initialize Database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Initialize services
	gameService := services.NewGameService(db, redisClient)
	playerService := services.NewPlayerService(db, redisClient)
	socketService := services.NewSocketService(redisClient)

	// Initialize handlers
	gameHandler := handlers.NewGameHandler(gameService)
	playerHandler := handlers.NewPlayerHandler(playerService)
	socketHandler := handlers.NewSocketHandler(socketService, gameService)

	// Initialize server
	srv := server.NewServer()

	// Setup middleware
	srv.Use(middleware.CORS())
	srv.Use(middleware.Logger())

	// Setup routes
	api := srv.Group("/api/v1")
	{
		// Player routes
		api.POST("/players", playerHandler.CreatePlayer)
		api.GET("/players/:id", playerHandler.GetPlayer)
		
		// Game routes
		api.POST("/games", gameHandler.CreateGame)
		api.GET("/games", gameHandler.ListGames)
		api.GET("/games/:id", gameHandler.GetGame)
		api.POST("/games/:id/join", gameHandler.JoinGame)
		api.POST("/games/:id/leave", gameHandler.LeaveGame)
		api.POST("/games/:id/start", gameHandler.StartGame)
		api.POST("/games/:id/end", gameHandler.EndGame)
	}

	// WebSocket route
	srv.GET("/ws", socketHandler.HandleWebSocket)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := srv.Run(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
} 