package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateUp runs all up migrations
func MigrateUp() error {
	m, err := getMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrated successfully")
	return nil
}

// MigrateDown rolls back all migrations
func MigrateDown() error {
	m, err := getMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Println("Database rolled back successfully")
	return nil
}

// MigrateTo migrates to a specific version
func MigrateTo(version uint) error {
	m, err := getMigrate()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Migrate(version); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	log.Printf("Database migrated to version %d successfully", version)
	return nil
}

// getMigrate creates a new migrate instance
func getMigrate() (*migrate.Migrate, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/tennis_tracker?sslmode=disable"
	}

	// Path to migration files - assumes you're running from the project root
	migrationsPath := "file://internal/database/migrations"

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}
