package db

import (
	"database/sql"
	"fmt"
	"log"

	"capstone1/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// ConnectDB connects to PostgreSQL using config
func ConnectDB(cfg *config.Config) *sql.DB {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	log.Println("Connecting to DB:", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	DB = db
	log.Println("Connected to DB successfully")
	return db
}

// InitTables creates necessary tables if they don't exist
func InitTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			totp_secret TEXT NOT NULL,
			role TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			action TEXT NOT NULL,
			timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS vault (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			data BYTEA NOT NULL,
			user_id INT NOT NULL
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatal("Failed to create table:", err)
		}
	}
}
