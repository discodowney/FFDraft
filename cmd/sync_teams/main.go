package main

import (
	"flag"
	"log"

	"go-app/config"
	"go-app/database"
	"go-app/external"
	"go-app/services/team"
	"go-app/services/team_sync"
)

func main() {
	log.Println("Starting team sync process...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

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
	log.Println("Database connection established")

	// Initialize API Football client
	apiFootballClient := external.NewAPIFootballClient(
		cfg.APIFootballBaseURL,
		cfg.APIFootballAPIKey,
		cfg.APIFootballLeagueID,
		cfg.APIFootballSeason,
	)
	log.Println("API Football client initialized")

	// Create services
	teamService := team.NewTeamService(db)
	teamSyncService := team_sync.NewTeamSyncService(teamService, apiFootballClient)
	log.Println("Services initialized")

	// Sync teams
	log.Println("Starting team sync from external API...")
	if err := teamSyncService.SyncTeamsFromExternalAPI(); err != nil {
		log.Fatalf("Failed to sync teams: %v", err)
	}
	log.Println("Team sync completed successfully")
}
