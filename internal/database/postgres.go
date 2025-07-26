package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgresDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping Postgres: %w", err)
	}

	// Create a users table if it doesn't exist
	if err := createUsersTable(db); err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	return db, nil
}

func createUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
	    id SERIAL PRIMARY KEY,
	    login VARCHAR(255) NOT NULL UNIQUE,
	    password_hash VARCHAR(255) NOT NULL,
	    balance INTEGER DEFAULT 10000, -- 100.00 in cents as a starting balance
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	);

	CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);
`
	_, err := db.Exec(query)
	return err
}
