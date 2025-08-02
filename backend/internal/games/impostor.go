package games

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"local-games/internal/models"
)

// ImpostorGame represents the Impostor game logic
type ImpostorGame struct {
	game     *models.Game
	players  []models.GamePlayer
	rounds   []models.GameRound
	currentRound *models.GameRound
}

// Category represents a word category with its words
type Category struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
}

// ImpostorCategories contains predefined categories and words
var ImpostorCategories = map[string]Category{
	"football_players": {
		Name: "Football Players",
		Words: []string{
			"Lionel Messi", "Cristiano Ronaldo", "Neymar", "Kylian Mbappé",
			"Erling Haaland", "Kevin De Bruyne", "Mohamed Salah", "Robert Lewandowski",
			"Karim Benzema", "Luka Modrić", "Virgil van Dijk", "Sadio Mané",
			"Harry Kane", "Raheem Sterling", "Jadon Sancho", "Phil Foden",
		},
	},
	"movies": {
		Name: "Movies",
		Words: []string{
			"The Godfather", "Titanic", "Avatar", "Star Wars",
			"Jurassic Park", "The Lion King", "Frozen", "Black Panther",
			"Avengers: Endgame", "Spider-Man", "Batman", "Superman",
			"Wonder Woman", "Iron Man", "Thor", "Captain America",
		},
	},
	"animals": {
		Name: "Animals",
		Words: []string{
			"Lion", "Tiger", "Elephant", "Giraffe",
			"Penguin", "Dolphin", "Eagle", "Shark",
			"Kangaroo", "Koala", "Panda", "Gorilla",
			"Zebra", "Hippo", "Rhino", "Cheetah",
		},
	},
	"countries": {
		Name: "Countries",
		Words: []string{
			"United States", "Canada", "Mexico", "Brazil",
			"Argentina", "Chile", "Peru", "Colombia",
			"France", "Germany", "Italy", "Spain",
			"United Kingdom", "Japan", "China", "India",
		},
	},
	"foods": {
		Name: "Foods",
		Words: []string{
			"Pizza", "Burger", "Sushi", "Pasta",
			"Taco", "Burrito", "Curry", "Ramen",
			"Steak", "Salmon", "Chicken", "Lobster",
			"Ice Cream", "Chocolate", "Cake", "Cookie",
		},
	},
	"colors": {
		Name: "Colors",
		Words: []string{
			"Red", "Blue", "Green", "Yellow",
			"Purple", "Orange", "Pink", "Brown",
			"Black", "White", "Gray", "Cyan",
			"Magenta", "Teal", "Navy", "Maroon",
		},
	},
}

// NewImpostorGame creates a new Impostor game instance
func NewImpostorGame(game *models.Game, players []models.GamePlayer) *ImpostorGame {
	return &ImpostorGame{
		game:    game,
		players: players,
		rounds:  []models.GameRound{},
	}
}

// StartGame initializes and starts the Impostor game
func (ig *ImpostorGame) StartGame() error {
	if len(ig.players) < ig.game.MinPlayers {
		return fmt.Errorf("not enough players to start game")
	}

	// Update game status
	ig.game.Status = "active"
	
	// Start first round
	return ig.startNewRound()
}

// startNewRound starts a new round in the game
func (ig *ImpostorGame) startNewRound() error {
	roundNumber := len(ig.rounds) + 1
	
	// Select a random category
	category := ig.selectRandomCategory()
	
	// Select a random word from the category
	word := ig.selectRandomWord(category)
	
	// Select a random impostor
	impostor := ig.selectRandomImpostor()
	
	// Create new round
	ig.currentRound = models.NewGameRound(ig.game.ID, roundNumber, category.Name, word, impostor.PlayerID)
	ig.rounds = append(ig.rounds, *ig.currentRound)
	
	return nil
}

// selectRandomCategory selects a random category from available categories
func (ig *ImpostorGame) selectRandomCategory() Category {
	// Get available categories from game settings
	var availableCategories []Category
	
	if ig.game.Settings.Categories != nil && len(ig.game.Settings.Categories) > 0 {
		// Use categories from game settings
		for _, catName := range ig.game.Settings.Categories {
			if cat, exists := ImpostorCategories[catName]; exists {
				availableCategories = append(availableCategories, cat)
			}
		}
	}
	
	// If no categories specified or none found, use all available
	if len(availableCategories) == 0 {
		for _, cat := range ImpostorCategories {
			availableCategories = append(availableCategories, cat)
		}
	}
	
	// Select random category
	rand.Seed(time.Now().UnixNano())
	return availableCategories[rand.Intn(len(availableCategories))]
}

// selectRandomWord selects a random word from the given category
func (ig *ImpostorGame) selectRandomWord(category Category) string {
	rand.Seed(time.Now().UnixNano())
	return category.Words[rand.Intn(len(category.Words))]
}

// selectRandomImpostor selects a random player to be the impostor
func (ig *ImpostorGame) selectRandomImpostor() models.GamePlayer {
	rand.Seed(time.Now().UnixNano())
	return ig.players[rand.Intn(len(ig.players))]
}

// GetPlayerWord returns the word for a specific player
func (ig *ImpostorGame) GetPlayerWord(playerID string) string {
	if ig.currentRound == nil {
		return ""
	}
	
	// If player is the impostor, return a different word or empty
	if playerID == ig.currentRound.ImpostorID {
		return "" // Impostor gets no word
	}
	
	return ig.currentRound.Word
}

// GetPlayerRole returns the role for a specific player
func (ig *ImpostorGame) GetPlayerRole(playerID string) string {
	if ig.currentRound == nil {
		return "player"
	}
	
	if playerID == ig.currentRound.ImpostorID {
		return "impostor"
	}
	
	return "player"
}

// GetGameState returns the current game state
func (ig *ImpostorGame) GetGameState() models.GameState {
	return models.GameState{
		Game:         ig.game,
		Players:      ig.players,
		CurrentRound: ig.currentRound,
		RoundHistory: ig.rounds,
	}
}

// EndRound ends the current round
func (ig *ImpostorGame) EndRound() error {
	if ig.currentRound == nil {
		return fmt.Errorf("no active round to end")
	}
	
	ig.currentRound.Status = "finished"
	now := time.Now()
	ig.currentRound.EndedAt = &now
	
	// Check if game should end
	if len(ig.rounds) >= ig.game.Settings.Rounds {
		ig.game.Status = "finished"
		return nil
	}
	
	// Start next round
	return ig.startNewRound()
}

// GetCategories returns all available categories
func (ig *ImpostorGame) GetCategories() map[string]Category {
	return ImpostorCategories
}

// GetCategoryByName returns a specific category by name
func (ig *ImpostorGame) GetCategoryByName(name string) (Category, bool) {
	cat, exists := ImpostorCategories[name]
	return cat, exists
}

// ValidateGameSettings validates the game settings
func (ig *ImpostorGame) ValidateGameSettings(settings models.GameSettings) error {
	if settings.Rounds < 1 {
		return fmt.Errorf("rounds must be at least 1")
	}
	
	if settings.TimePerRound < 30 {
		return fmt.Errorf("time per round must be at least 30 seconds")
	}
	
	if settings.Categories != nil {
		for _, catName := range settings.Categories {
			if _, exists := ImpostorCategories[catName]; !exists {
				return fmt.Errorf("invalid category: %s", catName)
			}
		}
	}
	
	return nil
} 