package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func NewPostgres(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Database ping failed: %w", err)
	}

	return db, nil
}