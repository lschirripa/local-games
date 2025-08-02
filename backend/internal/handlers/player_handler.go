package handlers

import (
	"net/http"

	"local-games/internal/models"
	"local-games/internal/services"

	"github.com/gin-gonic/gin"
)

// PlayerHandler handles player-related HTTP requests
type PlayerHandler struct {
	playerService *services.PlayerService
}

// NewPlayerHandler creates a new player handler
func NewPlayerHandler(playerService *services.PlayerService) *PlayerHandler {
	return &PlayerHandler{
		playerService: playerService,
	}
}

// CreatePlayer handles POST /api/v1/players
func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var req models.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := h.playerService.CreatePlayer(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Player created successfully",
		"player":  player,
	})
}

// GetPlayer handles GET /api/v1/players/:id
func (h *PlayerHandler) GetPlayer(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player ID is required"})
		return
	}

	player, err := h.playerService.GetPlayer(playerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"player": player,
	})
}

// GetPlayerSession handles GET /api/v1/players/:id/session
func (h *PlayerHandler) GetPlayerSession(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player ID is required"})
		return
	}

	session, err := h.playerService.GetPlayerSession(playerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session": session,
	})
}

// GetPlayerGames handles GET /api/v1/players/:id/games
func (h *PlayerHandler) GetPlayerGames(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player ID is required"})
		return
	}

	games, err := h.playerService.GetPlayerGames(playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"games": games,
	})
}

// GetPlayerStats handles GET /api/v1/players/:id/stats
func (h *PlayerHandler) GetPlayerStats(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player ID is required"})
		return
	}

	stats, err := h.playerService.GetPlayerStats(playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
} 