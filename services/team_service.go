package services

import (
	"errors"
	"fmt"
	"go-app/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// TeamService handles database operations for teams
type TeamService struct {
	db *sqlx.DB
}

// NewTeamService creates a new TeamService instance
func NewTeamService(db *sqlx.DB) *TeamService {
	return &TeamService{
		db: db,
	}
}

// GetTeam retrieves a team by ID
func (s *TeamService) GetTeam(id int) (*models.Team, error) {
	// TODO: Implement database query to get team by ID
	// For now, return a placeholder
	return &models.Team{
		ID:         id,
		Name:       "Team Name",
		ExternalId: 1000 + id, // Example external ID
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

// ListTeams retrieves teams based on filters
func (s *TeamService) ListTeams(filters map[string]interface{}) ([]*models.Team, error) {
	// TODO: Implement database query with filters
	// For now, return placeholder data
	teams := []*models.Team{
		{
			ID:         1,
			Name:       "Team 1",
			ExternalId: 1001,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         2,
			Name:       "Team 2",
			ExternalId: 1002,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}
	return teams, nil
}

// GetTeamsByLeague retrieves all teams for a specific league
func (s *TeamService) GetTeamsByLeague(leagueID int) ([]*models.Team, error) {
	// TODO: Implement database query to get teams by league ID
	// For now, return placeholder data
	teams := []*models.Team{
		{
			ID:         1,
			Name:       "League Team 1",
			ExternalId: 1001,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         2,
			Name:       "League Team 2",
			ExternalId: 1002,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}
	return teams, nil
}

// ValidateTeam validates a team's data
func (s *TeamService) ValidateTeam(team *models.Team) error {
	if team.Name == "" {
		return errors.New("team name is required")
	}
	return nil
}

// GetTeamByExternalID retrieves a team by its external ID
func (s *TeamService) GetTeamByExternalID(externalID int) (*models.Team, error) {
	query := `SELECT * FROM teams WHERE external_id = $1`

	team := &models.Team{}
	err := s.db.Get(team, query, externalID)
	if err != nil {
		return nil, fmt.Errorf("error getting team by external ID: %w", err)
	}

	return team, nil
}

// CreateTeam creates a new team in the database
func (s *TeamService) CreateTeam(team *models.Team) error {
	query := `INSERT INTO teams (name, external_id, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4) 
			  RETURNING id`

	now := time.Now()
	team.CreatedAt = now
	team.UpdatedAt = now

	err := s.db.QueryRow(
		query,
		team.Name,
		team.ExternalId,
		team.CreatedAt,
		team.UpdatedAt,
	).Scan(&team.ID)

	if err != nil {
		return fmt.Errorf("error creating team: %w", err)
	}

	return nil
}

// UpdateTeam updates an existing team in the database
func (s *TeamService) UpdateTeam(team *models.Team) error {
	query := `UPDATE teams 
			  SET name = $1, updated_at = $2 
			  WHERE id = $3`

	team.UpdatedAt = time.Now()

	result, err := s.db.Exec(query, team.Name, team.UpdatedAt, team.ID)
	if err != nil {
		return fmt.Errorf("error updating team: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("team with ID %d not found", team.ID)
	}

	return nil
}
