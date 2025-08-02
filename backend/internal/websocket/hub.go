package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	ID       string
	PlayerID string
	GameID   string
	Conn     *websocket.Conn
	Send     chan []byte
	mu       sync.Mutex
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	rooms      map[string]map[*Client]bool
	mu         sync.RWMutex
}

type Message struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	GameID  string      `json:"game_id,omitempty"`
	PlayerID string     `json:"player_id,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if client.GameID != "" {
				if h.rooms[client.GameID] == nil {
					h.rooms[client.GameID] = make(map[*Client]bool)
				}
				h.rooms[client.GameID][client] = true
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			if client.GameID != "" {
				if room, exists := h.rooms[client.GameID]; exists {
					delete(room, client)
					if len(room) == 0 {
						delete(h.rooms, client.GameID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastToRoom(gameID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, exists := h.rooms[gameID]; exists {
		for client := range room {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
				delete(room, client)
			}
		}
	}
}

func (h *Hub) SendToPlayer(playerID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.PlayerID == playerID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Parse and handle message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error parsing message: %v", err)
			continue
		}

		// Handle different message types
		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case "join_game":
		c.GameID = msg.GameID
		c.PlayerID = msg.PlayerID
		// Notify other players in the room
		response := Message{
			Type:     "player_joined",
			PlayerID: c.PlayerID,
			GameID:   c.GameID,
		}
		if data, err := json.Marshal(response); err == nil {
			c.Hub.BroadcastToRoom(c.GameID, data)
		}

	case "game_action":
		// Handle game-specific actions
		response := Message{
			Type:     "game_update",
			Data:     msg.Data,
			GameID:   c.GameID,
			PlayerID: c.PlayerID,
		}
		if data, err := json.Marshal(response); err == nil {
			c.Hub.BroadcastToRoom(c.GameID, data)
		}

	case "leave_game":
		// Notify other players
		response := Message{
			Type:     "player_left",
			PlayerID: c.PlayerID,
			GameID:   c.GameID,
		}
		if data, err := json.Marshal(response); err == nil {
			c.Hub.BroadcastToRoom(c.GameID, data)
		}
		c.GameID = ""
	}
}

func (c *Client) SendMessage(msgType string, data interface{}) {
	msg := Message{
		Type:     msgType,
		Data:     data,
		GameID:   c.GameID,
		PlayerID: c.PlayerID,
	}

	if message, err := json.Marshal(msg); err == nil {
		c.Send <- message
	}
} 