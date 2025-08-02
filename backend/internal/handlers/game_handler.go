package handlers

import (
	"net/http"

	"local-games/internal/models"
	"local-games/internal/services"

	"github.com/gin-gonic/gin"
)

// GameHandler handles game-related HTTP requests
type GameHandler struct {
	gameService *services.GameService
}

// NewGameHandler creates a new game handler
func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// CreateGame handles POST /api/v1/games
func (h *GameHandler) CreateGame(c *gin.Context) {
	var req models.CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get player ID from context (in a real app, this would come from authentication)
	playerID := c.GetHeader("X-Player-ID")
	if playerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "player ID required"})
		return
	}

	game, err := h.gameService.CreateGame(req, playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Game created successfully",
		"game":    game,
	})
}

// ListGames handles GET /api/v1/games
func (h *GameHandler) ListGames(c *gin.Context) {
	games, err := h.gameService.ListGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"games": games,
	})
}

// GetGame handles GET /api/v1/games/:id
func (h *GameHandler) GetGame(c *gin.Context) {
	gameID := c.Param("id")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game ID is required"})
		return
	}

	game, err := h.gameService.GetGame(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game": game,
	})
}

// JoinGame handles POST /api/v1/games/:id/join
func (h *GameHandler) JoinGame(c *gin.Context) {
	gameID := c.Param("id")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game ID is required"})
		return
	}

	var req models.JoinGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.gameService.JoinGame(gameID, req.PlayerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully joined game",
	})
}

// LeaveGame handles POST /api/v1/games/:id/leave
func (h *GameHandler) LeaveGame(c *gin.Context) {
	gameID := c.Param("id")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game ID is required"})
		return
	}

	playerID := c.GetHeader("X-Player-ID")
	if playerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "player ID required"})
		return
	}

	if err := h.gameService.LeaveGame(gameID, playerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully left game",
	})
}

// StartGame handles POST /api/v1/games/:id/start
func (h *GameHandler) StartGame(c *gin.Context) {
	gameID := c.Param("id")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game ID is required"})
		return
	}

	if err := h.gameService.StartGame(gameID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Game started successfully",
	})
}

// EndGame handles POST /api/v1/games/:id/end
func (h *GameHandler) EndGame(c *gin.Context) {
	gameID := c.Param("id")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game ID is required"})
		return
	}

	if err := h.gameService.EndGame(gameID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Game ended successfully",
	})
}

// GetGameState handles GET /api/v1/games/:id/state
func (h *GameHandler) GetGameState(c *gin.Context) {
	gameID := c.Param("id")
	if gameID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game ID is required"})
		return
	}

	gameState, err := h.gameService.GetGameState(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"game_state": gameState,
	})
}

// GetPlayerWord handles GET /api/v1/games/:id/player/:player_id/word
func (h *GameHandler) GetPlayerWord(c *gin.Context) {
	gameID := c.Param("id")
	playerID := c.Param("player_id")
	
	if gameID == "" || playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game ID and player ID are required"})
		return
	}

	word, err := h.gameService.GetPlayerWord(gameID, playerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"word": word,
	})
} 