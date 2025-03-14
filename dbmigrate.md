Let me help you create database migration scripts for your tennis tracking application. Database migrations allow you to version your database schema and make incremental changes in a controlled way.

For a Go application with PostgreSQL, there are several popular migration libraries. I'll show you how to set up migrations using `golang-migrate`, one of the most widely used tools.

## Setting Up Database Migrations

### Step 1: Install the migrate CLI tool

This is useful for creating and running migrations from the command line:

```bash
# For macOS
brew install golang-migrate

# For Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate

# For Windows (using scoop)
scoop install migrate
```

### Step 2: Create a migrations directory in your backend folder

```bash
mkdir -p backend/internal/database/migrations
```

### Step 3: Create your first migration files

```bash
cd backend
migrate create -ext sql -dir internal/database/migrations -seq initial_schema
```

This will create two files:
- `internal/database/migrations/000001_initial_schema.up.sql`
- `internal/database/migrations/000001_initial_schema.down.sql`

```sql
-- Tennis Tracker Initial Schema Migration

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create sessions table
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    opponent_name VARCHAR(255),
    session_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create errors table
CREATE TABLE errors (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES sessions(id) ON DELETE CASCADE,
    count INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create shared links table
CREATE TABLE shared_links (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    session_id INTEGER REFERENCES sessions(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_errors_session_id ON errors(session_id);
CREATE INDEX idx_shared_links_token ON shared_links(token);
CREATE INDEX idx_shared_links_user_id ON shared_links(user_id);
CREATE INDEX idx_shared_links_session_id ON shared_links(session_id);

```

```sql
-- Tennis Tracker Initial Schema Rollback

-- Drop indexes
DROP INDEX IF EXISTS idx_shared_links_session_id;
DROP INDEX IF EXISTS idx_shared_links_user_id;
DROP INDEX IF EXISTS idx_shared_links_token;
DROP INDEX IF EXISTS idx_errors_session_id;
DROP INDEX IF EXISTS idx_sessions_user_id;

-- Drop tables (in reverse order of creation to handle dependencies)
DROP TABLE IF EXISTS shared_links;
DROP TABLE IF EXISTS errors;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

```

### Step 4: Set up a database connection utility in Go

```go
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

```

### Step 5: Create a migration utility to run migrations from your Go code

```go
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

```

### Step 6: Update your main.go to run migrations at startup

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/tennis-tracker/backend/internal/database"
)

func main() {
	// Run database migrations
	if err := database.MigrateUp(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize database connection
	db, err := database.New()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Setup API server and other components here
	// ...

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

```

### Step 7: Add a Makefile to simplify migration commands

```makefile
.PHONY: migrate-up migrate-down migrate-create db-reset run

# Database migration commands
migrate-up:
	migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/tennis_tracker?sslmode=disable" up

migrate-down:
	migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/tennis_tracker?sslmode=disable" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/database/migrations -seq $$name

# Reset database (down all migrations then up)
db-reset:
	migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/tennis_tracker?sslmode=disable" down -all
	migrate -path internal/database/migrations -database "postgres://postgres:postgres@localhost:5432/tennis_tracker?sslmode=disable" up

# Run the application
run:
	go run cmd/server/main.go

```

## How to Use the Migration System

### Creating new migrations

When you need to make changes to your database schema (like adding a new table or modifying an existing one), create a new migration:

```bash
# Using the migrate CLI directly
migrate create -ext sql -dir internal/database/migrations -seq add_user_role

# Or using the Makefile
make migrate-create
# When prompted, enter: add_user_role
```

This will create two new files:
- `000002_add_user_role.up.sql` - Contains the changes to apply
- `000002_add_user_role.down.sql` - Contains how to reverse those changes

### Running migrations manually

```bash
# Apply all pending migrations
make migrate-up

# Rollback the latest migration
make migrate-down

# Reset the database (useful during development)
make db-reset
```

### Running migrations from your application

The application will automatically run migrations when it starts, thanks to the code in `main.go`.

## Benefits of This Migration Approach

1. **Version Control**: Your database schema changes are versioned alongside your code
2. **Reproducibility**: Anyone can recreate the database with the same schema
3. **Rollback Capability**: You can undo changes if something goes wrong
4. **CI/CD Friendly**: Migrations can be run automatically during deployment
5. **Development/Production Parity**: Ensures all environments have the same schema

Would you like me to explain any specific part of the migration process in more detail or help with creating additional migrations for future phases of your project?
