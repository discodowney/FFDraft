package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go-app/config"
	"go-app/database"
	v1 "go-app/server/v1"
)

func main() {
	// Parse command line flags
	env := flag.String("env", "development", "Environment (development/production)")
	flag.Parse()

	// Load environment variables based on the environment
	if err := loadEnv(*env); err != nil {
		log.Fatalf("Error loading environment: %v", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	// Start the server using v1 handler
	v1.StartServer(db)
}

func loadEnv(env string) error {
	// Try to load .env first
	if _, err := os.Stat(".env"); err == nil {
		fmt.Printf("Loading environment from .env\n")
		return nil
	}

	// If .env doesn't exist, try environment-specific file
	envFile := fmt.Sprintf(".env.%s", env)
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return fmt.Errorf("neither .env nor %s found", envFile)
	}

	fmt.Printf("Loading environment from %s\n", envFile)
	return nil
}
