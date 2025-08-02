package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Game struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Players    int       `json:"players"`
	MaxPlayers int       `json:"maxPlayers"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var games []Game
var players []Player

func main() {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/games", getGames).Methods("GET")
	api.HandleFunc("/games", createGame).Methods("POST")
	api.HandleFunc("/games/{id}", getGame).Methods("GET")
	api.HandleFunc("/players", getPlayers).Methods("GET")
	api.HandleFunc("/players", createPlayer).Methods("POST")

	// Enable CORS for frontend communication
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // SvelteKit dev server
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func getGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func createGame(w http.ResponseWriter, r *http.Request) {
	var game Game
	json.NewDecoder(r.Body).Decode(&game)
	game.ID = fmt.Sprintf("game_%d", time.Now().Unix())
	game.CreatedAt = time.Now()
	games = append(games, game)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func getGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["id"]
	
	for _, game := range games {
		if game.ID == gameID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(game)
			return
		}
	}
	
	http.Error(w, "Game not found", http.StatusNotFound)
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

func createPlayer(w http.ResponseWriter, r *http.Request) {
	var player Player
	json.NewDecoder(r.Body).Decode(&player)
	player.ID = fmt.Sprintf("player_%d", time.Now().Unix())
	players = append(players, player)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}