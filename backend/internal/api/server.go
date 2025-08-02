package api

import (
	"fmt"
	"net/http"
	"time"

	"local-games/backend/internal/config"
	"local-games/backend/internal/database"
	"local-games/backend/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config *config.Config
	db     *database.DB
	redis  *database.RedisClient
	hub    *websocket.Hub
	router *gin.Engine
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

func NewServer(cfg *config.Config, db *database.DB, redis *database.RedisClient, hub *websocket.Hub) *Server {
	server := &Server{
		config: cfg,
		db:     db,
		redis:  redis,
		hub:    hub,
		router: gin.Default(),
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

func (s *Server) setupMiddleware() {
	// CORS middleware
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:80"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Logging middleware
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)

	// API routes
	api := s.router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/session", s.createSession)
			auth.DELETE("/session", s.deleteSession)
			auth.GET("/me", s.getCurrentPlayer)
		}

		// Game routes
		games := api.Group("/games")
		{
			games.GET("", s.getGames)
			games.POST("", s.createGame)
			games.GET("/:id", s.getGame)
			games.PUT("/:id", s.updateGame)
			games.DELETE("/:id", s.deleteGame)
			games.POST("/:id/join", s.joinGame)
			games.POST("/:id/leave", s.leaveGame)
		}

		// Player routes
		players := api.Group("/players")
		{
			players.GET("/:id", s.getPlayer)
			players.PUT("/:id", s.updatePlayer)
		}
	}

	// WebSocket endpoint
	s.router.GET("/ws", s.handleWebSocket)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	logrus.Infof("Starting server on %s", addr)
	return s.router.Run(addr)
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "Server is running",
	})
}

func (s *Server) handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &websocket.Client{
		Hub:  s.hub,
		ID:   generateClientID(),
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	client.Hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

func generateClientID() string {
	// Simple client ID generation - in production, use UUID
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
} 