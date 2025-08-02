package database

import (
	"fmt"
	"log"

	"local-games/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewConnection creates a new database connection
func NewConnection() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connection established successfully")
	return db, nil
}

// InitTables creates the necessary database tables
func InitTables(db *sqlx.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS players (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS games (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			type VARCHAR(50) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'waiting',
			max_players INTEGER NOT NULL,
			min_players INTEGER NOT NULL,
			created_by VARCHAR(36) NOT NULL,
			settings JSONB,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS game_players (
			id VARCHAR(36) PRIMARY KEY,
			game_id VARCHAR(36) NOT NULL REFERENCES games(id) ON DELETE CASCADE,
			player_id VARCHAR(36) NOT NULL REFERENCES players(id) ON DELETE CASCADE,
			role VARCHAR(20) NOT NULL DEFAULT 'player',
			score INTEGER NOT NULL DEFAULT 0,
			joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(game_id, player_id)
		)`,
		`CREATE TABLE IF NOT EXISTS game_rounds (
			id VARCHAR(36) PRIMARY KEY,
			game_id VARCHAR(36) NOT NULL REFERENCES games(id) ON DELETE CASCADE,
			round_number INTEGER NOT NULL,
			category VARCHAR(100) NOT NULL,
			word VARCHAR(100) NOT NULL,
			impostor_id VARCHAR(36) REFERENCES players(id),
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			started_at TIMESTAMP NOT NULL DEFAULT NOW(),
			ended_at TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS votes (
			id VARCHAR(36) PRIMARY KEY,
			round_id VARCHAR(36) NOT NULL REFERENCES game_rounds(id) ON DELETE CASCADE,
			voter_id VARCHAR(36) NOT NULL REFERENCES players(id) ON DELETE CASCADE,
			voted_for_id VARCHAR(36) NOT NULL REFERENCES players(id) ON DELETE CASCADE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(round_id, voter_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_games_status ON games(status)`,
		`CREATE INDEX IF NOT EXISTS idx_game_players_game_id ON game_players(game_id)`,
		`CREATE INDEX IF NOT EXISTS idx_game_rounds_game_id ON game_rounds(game_id)`,
		`CREATE INDEX IF NOT EXISTS idx_votes_round_id ON votes(round_id)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	log.Println("Database tables initialized successfully")
	return nil
} 