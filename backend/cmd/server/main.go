package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jimsyyap/tennis-tracker/backend/internal/database"
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
