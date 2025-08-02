package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisHost  string
	RedisPort  string
	RedisDB    string
	JWTSecret  string
	Port       string
}

var AppConfig Config

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		// Don't fail if .env doesn't exist, use system env vars
	}

	AppConfig = Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "local_games"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		RedisDB:    getEnv("REDIS_DB", "0"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		Port:       getEnv("PORT", "8080"),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 