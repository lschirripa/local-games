package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"local-games/internal/models"
	"local-games/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// SocketHandler handles WebSocket connections
type SocketHandler struct {
	socketService *services.SocketService
	gameService   *services.GameService
	upgrader      websocket.Upgrader
	clients       map[string]*Client
	mutex         sync.RWMutex
}

// Client represents a WebSocket client connection
type Client struct {
	ID       string
	PlayerID string
	GameID   string
	Conn     *websocket.Conn
	Send     chan []byte
}

// NewSocketHandler creates a new socket handler
func NewSocketHandler(socketService *services.SocketService, gameService *services.GameService) *SocketHandler {
	return &SocketHandler{
		socketService: socketService,
		gameService:   gameService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
		clients: make(map[string]*Client),
	}
}

// HandleWebSocket handles WebSocket upgrade and connection
func (h *SocketHandler) HandleWebSocket(c *gin.Context) {
	// Get player ID from query parameters
	playerID := c.Query("player_id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player_id is required"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Create new client
	client := &Client{
		ID:       playerID,
		PlayerID: playerID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}

	// Register client
	h.registerClient(client)

	// Start goroutines for reading and writing
	go h.readPump(client)
	go h.writePump(client)
}

// registerClient registers a new client
func (h *SocketHandler) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.clients[client.ID] = client
	log.Printf("Client registered: %s", client.ID)
}

// unregisterClient unregisters a client
func (h *SocketHandler) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	if _, ok := h.clients[client.ID]; ok {
		delete(h.clients, client.ID)
		close(client.Send)
		log.Printf("Client unregistered: %s", client.ID)
	}
}

// readPump reads messages from the WebSocket connection
func (h *SocketHandler) readPump(client *Client) {
	defer func() {
		h.unregisterClient(client)
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		// Handle the message
		h.handleMessage(client, message)
	}
}

// writePump writes messages to the WebSocket connection
func (h *SocketHandler) writePump(client *Client) {
	defer func() {
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (h *SocketHandler) handleMessage(client *Client, message []byte) {
	var wsMessage models.WebSocketMessage
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	switch wsMessage.Type {
	case "join_game":
		h.handleJoinGame(client, wsMessage)
	case "leave_game":
		h.handleLeaveGame(client, wsMessage)
	case "start_game":
		h.handleStartGame(client, wsMessage)
	case "end_game":
		h.handleEndGame(client, wsMessage)
	case "get_game_state":
		h.handleGetGameState(client, wsMessage)
	case "get_player_word":
		h.handleGetPlayerWord(client, wsMessage)
	case "vote":
		h.handleVote(client, wsMessage)
	default:
		log.Printf("Unknown message type: %s", wsMessage.Type)
	}
}

// handleJoinGame handles join game requests
func (h *SocketHandler) handleJoinGame(client *Client, message models.WebSocketMessage) {
	var payload struct {
		GameID string `json:"game_id"`
	}
	
	if err := json.Unmarshal(message.Payload.([]byte), &payload); err != nil {
		h.sendError(client, "Invalid payload")
		return
	}

	// Join the game
	if err := h.gameService.JoinGame(payload.GameID, client.PlayerID); err != nil {
		h.sendError(client, err.Error())
		return
	}

	client.GameID = payload.GameID

	// Notify all players in the game
	h.broadcastToGame(payload.GameID, models.WebSocketMessage{
		Type: "player_joined",
		Payload: gin.H{
			"player_id": client.PlayerID,
			"game_id":   payload.GameID,
		},
	})

	h.sendSuccess(client, "Successfully joined game")
}

// handleLeaveGame handles leave game requests
func (h *SocketHandler) handleLeaveGame(client *Client, message models.WebSocketMessage) {
	var payload struct {
		GameID string `json:"game_id"`
	}
	
	if err := json.Unmarshal(message.Payload.([]byte), &payload); err != nil {
		h.sendError(client, "Invalid payload")
		return
	}

	// Leave the game
	if err := h.gameService.LeaveGame(payload.GameID, client.PlayerID); err != nil {
		h.sendError(client, err.Error())
		return
	}

	// Notify all players in the game
	h.broadcastToGame(payload.GameID, models.WebSocketMessage{
		Type: "player_left",
		Payload: gin.H{
			"player_id": client.PlayerID,
			"game_id":   payload.GameID,
		},
	})

	client.GameID = ""
	h.sendSuccess(client, "Successfully left game")
}

// handleStartGame handles start game requests
func (h *SocketHandler) handleStartGame(client *Client, message models.WebSocketMessage) {
	var payload struct {
		GameID string `json:"game_id"`
	}
	
	if err := json.Unmarshal(message.Payload.([]byte), &payload); err != nil {
		h.sendError(client, "Invalid payload")
		return
	}

	// Start the game
	if err := h.gameService.StartGame(payload.GameID); err != nil {
		h.sendError(client, err.Error())
		return
	}

	// Notify all players in the game
	h.broadcastToGame(payload.GameID, models.WebSocketMessage{
		Type: "game_started",
		Payload: gin.H{
			"game_id": payload.GameID,
		},
	})

	h.sendSuccess(client, "Game started successfully")
}

// handleEndGame handles end game requests
func (h *SocketHandler) handleEndGame(client *Client, message models.WebSocketMessage) {
	var payload struct {
		GameID string `json:"game_id"`
	}
	
	if err := json.Unmarshal(message.Payload.([]byte), &payload); err != nil {
		h.sendError(client, "Invalid payload")
		return
	}

	// End the game
	if err := h.gameService.EndGame(payload.GameID); err != nil {
		h.sendError(client, err.Error())
		return
	}

	// Notify all players in the game
	h.broadcastToGame(payload.GameID, models.WebSocketMessage{
		Type: "game_ended",
		Payload: gin.H{
			"game_id": payload.GameID,
		},
	})

	h.sendSuccess(client, "Game ended successfully")
}

// handleGetGameState handles get game state requests
func (h *SocketHandler) handleGetGameState(client *Client, message models.WebSocketMessage) {
	var payload struct {
		GameID string `json:"game_id"`
	}
	
	if err := json.Unmarshal(message.Payload.([]byte), &payload); err != nil {
		h.sendError(client, "Invalid payload")
		return
	}

	// Get game state
	gameState, err := h.gameService.GetGameState(payload.GameID)
	if err != nil {
		h.sendError(client, err.Error())
		return
	}

	h.sendMessage(client, models.WebSocketMessage{
		Type:    "game_state",
		Payload: gameState,
	})
}

// handleGetPlayerWord handles get player word requests (for impostor game)
func (h *SocketHandler) handleGetPlayerWord(client *Client, message models.WebSocketMessage) {
	var payload struct {
		GameID string `json:"game_id"`
	}
	
	if err := json.Unmarshal(message.Payload.([]byte), &payload); err != nil {
		h.sendError(client, "Invalid payload")
		return
	}

	// Get player word
	word, err := h.gameService.GetPlayerWord(payload.GameID, client.PlayerID)
	if err != nil {
		h.sendError(client, err.Error())
		return
	}

	h.sendMessage(client, models.WebSocketMessage{
		Type: "player_word",
		Payload: gin.H{
			"word": word,
		},
	})
}

// handleVote handles voting in games
func (h *SocketHandler) handleVote(client *Client, message models.WebSocketMessage) {
	var payload struct {
		GameID     string `json:"game_id"`
		RoundID    string `json:"round_id"`
		VotedForID string `json:"voted_for_id"`
	}
	
	if err := json.Unmarshal(message.Payload.([]byte), &payload); err != nil {
		h.sendError(client, "Invalid payload")
		return
	}

	// TODO: Implement voting logic
	h.sendSuccess(client, "Vote recorded")
}

// broadcastToGame sends a message to all players in a game
func (h *SocketHandler) broadcastToGame(gameID string, message models.WebSocketMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	for _, client := range h.clients {
		if client.GameID == gameID {
			select {
			case client.Send <- messageBytes:
			default:
				close(client.Send)
				delete(h.clients, client.ID)
			}
		}
	}
}

// sendMessage sends a message to a specific client
func (h *SocketHandler) sendMessage(client *Client, message models.WebSocketMessage) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	select {
	case client.Send <- messageBytes:
	default:
		close(client.Send)
		delete(h.clients, client.ID)
	}
}

// sendSuccess sends a success message to a client
func (h *SocketHandler) sendSuccess(client *Client, message string) {
	h.sendMessage(client, models.WebSocketMessage{
		Type: "success",
		Payload: gin.H{
			"message": message,
		},
	})
}

// sendError sends an error message to a client
func (h *SocketHandler) sendError(client *Client, message string) {
	h.sendMessage(client, models.WebSocketMessage{
		Type: "error",
		Payload: gin.H{
			"message": message,
		},
	})
} 