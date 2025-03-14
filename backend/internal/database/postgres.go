package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DB represents the database connection pool
type DB struct {
	Pool *pgxpool.Pool
}

// New creates a new database connection
func New() (*DB, error) {
	// Get database connection string from environment variable
	// Format: "postgres://username:password@localhost:5432/database_name"
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		// Fallback to default connection string if not set
		connString = "postgres://postgres:postgres@localhost:5432/tennis_tracker?sslmode=disable"
	}

	// Create a connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
