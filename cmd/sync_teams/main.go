package main

import (
	"flag"
	"fmt"
	"log"

	"go-app/config"
	"go-app/database"
	"go-app/external"
	"go-app/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	flag.Parse()

	// Initialize database
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create API-Football client
	apiFootballClient := external.NewAPIFootballClient(cfg.APIFootballAPIKey)

	// Create services
	teamService := services.NewTeamService(db)
	teamSyncService := services.NewTeamSyncService(teamService, apiFootballClient)

	// Sync teams

	// Sync all teams
	if err := teamSyncService.SyncTeamsFromExternalAPI(); err != nil {
		log.Fatalf("Failed to sync teams: %v", err)
	}
	fmt.Println("Successfully synced all teams")

}
