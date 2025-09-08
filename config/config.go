package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all app configuration
type Config struct {
	DBUser     string
	DBPass     string
	DBHost     string
	DBPort     string
	DBName     string
	JWTSecret  string
	VaultKey   string
	ServerPort string
}

// LoadConfig loads environment variables from .env file or system
func LoadConfig() *Config {
	// Load .env first
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPass:     getEnv("DB_PASS", "postgres"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "idam"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecretkey"),
		VaultKey:   getEnv("VAULT_KEY", "32byteslongsecretkeyforvault123456"), // must be 32 bytes
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
	log.Println("Loaded DB user:", cfg.DBUser)
	return cfg
}

// helper to read env with fallback
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
