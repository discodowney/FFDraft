package team_sync

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"go-app/external"
	"go-app/models"
	"go-app/services/team"
)

// TeamSyncService handles synchronizing teams from external sources
type TeamSyncService struct {
	teamService       team.TeamService
	apiFootballClient external.APIFootballClientInterface
}

// NewTeamSyncService creates a new TeamSyncService instance
func NewTeamSyncService(teamService team.TeamService, apiFootballClient external.APIFootballClientInterface) *TeamSyncService {
	return &TeamSyncService{
		teamService:       teamService,
		apiFootballClient: apiFootballClient,
	}
}

// SyncTeamsFromExternalAPI fetches teams from the API-Football service and saves them to the database
func (s *TeamSyncService) SyncTeamsFromExternalAPI() error {
	log.Println("Starting team sync from API-Football")

	// Fetch teams from the API-Football service
	teams, err := s.apiFootballClient.FetchTeams()
	if err != nil {
		return fmt.Errorf("failed to fetch teams from API-Football: %w", err)
	}

	log.Printf("Fetched %d teams from API-Football", len(teams))

	// Process each team
	for _, team := range teams {
		// Check if team already exists
		existingTeam, err := s.teamService.GetTeamByExternalID(int64(team.ExternalId))
		if err != nil {
			log.Printf("Error checking for existing team %d: %v", team.ExternalId, err)
			continue
		}

		if existingTeam != nil {
			// Update existing team
			existingTeam.Name = team.Name
			_, err := s.teamService.UpdateTeam(existingTeam)
			if err != nil {
				log.Printf("Error updating team %d: %v", team.ExternalId, err)
				continue
			}
			log.Printf("Updated team: %s (ID: %d)", team.Name, team.ExternalId)
		} else {
			// Create new team
			newTeam := &models.Team{
				Name:       team.Name,
				ExternalId: team.ExternalId,
			}
			_, err := s.teamService.CreateTeam(newTeam)
			if err != nil {
				log.Printf("Error creating team %d: %v", team.ExternalId, err)
				continue
			}
			log.Printf("Created team: %s (ID: %d)", team.Name, team.ExternalId)
		}
	}

	log.Println("Team sync completed successfully")
	return nil
}

// SyncTeamByExternalID fetches a specific team from the API-Football service and saves it to the database
func (s *TeamSyncService) SyncTeamByExternalID(externalID int) error {
	log.Printf("Starting sync for team with external ID: %d", externalID)

	// Fetch team from the API-Football service
	team, err := s.apiFootballClient.FetchTeamByExternalID(externalID)
	if err != nil {
		return fmt.Errorf("failed to fetch team from API-Football: %w", err)
	}

	// Check if team already exists
	existingTeam, err := s.teamService.GetTeamByExternalID(int64(team.ExternalId))
	if err != nil {
		return fmt.Errorf("error checking for existing team: %w", err)
	}

	if existingTeam != nil {
		// Update existing team
		existingTeam.Name = team.Name
		_, err := s.teamService.UpdateTeam(existingTeam)
		if err != nil {
			return fmt.Errorf("error updating team: %w", err)
		}
		log.Printf("Updated team: %s (ID: %d)", team.Name, team.ExternalId)
	} else {
		// Create new team
		newTeam := &models.Team{
			Name:       team.Name,
			ExternalId: team.ExternalId,
		}
		_, err := s.teamService.CreateTeam(newTeam)
		if err != nil {
			return fmt.Errorf("error creating team: %w", err)
		}
		log.Printf("Created team: %s (ID: %d)", team.Name, team.ExternalId)
	}

	log.Printf("Team sync completed successfully for external ID: %d", externalID)
	return nil
}

// SyncTeams synchronizes teams from the external API
func (s *TeamSyncService) SyncTeams() error {
	teams, err := s.apiFootballClient.FetchTeams()
	if err != nil {
		return err
	}

	now := time.Now()

	for _, team := range teams {
		// Check if team already exists
		existingTeam, err := s.teamService.GetTeamByExternalID(int64(team.ExternalId))
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if existingTeam != nil {
			// Update existing team
			existingTeam.Name = team.Name
			_, err := s.teamService.UpdateTeam(existingTeam)
			if err != nil {
				log.Printf("Error updating team %d: %v", team.ExternalId, err)
				continue
			}
			log.Printf("Updated team: %s (ID: %d)", team.Name, team.ExternalId)
		} else {
			// Create new team
			newTeam := &models.Team{
				Name:       team.Name,
				ExternalId: team.ExternalId,
				CreatedAt:  now,
				UpdatedAt:  now,
			}
			_, err := s.teamService.CreateTeam(newTeam)
			if err != nil {
				log.Printf("Error creating team %d: %v", team.ExternalId, err)
				continue
			}
			log.Printf("Created team: %s (ID: %d)", team.Name, team.ExternalId)
		}
	}

	return nil
}
