package team

import (
	"fmt"
	"time"

	"go-app/models"

	"github.com/jmoiron/sqlx"
)

// TeamService defines the interface for team-related operations
type TeamService interface {
	CreateTeam(team *models.Team) (*models.Team, error)
	GetTeam(id int64) (*models.Team, error)
	UpdateTeam(team *models.Team) (*models.Team, error)
	DeleteTeam(id int64) error
	ListTeams() ([]*models.Team, error)
	ValidateTeam(team *models.Team) error
	GetTeamByExternalID(externalID int64) (*models.Team, error)
}

// TeamFilter represents the filter criteria for listing teams
type TeamFilter struct {
	Name       string
	ExternalID string
}

// Implementation of the TeamService interface
type teamServiceImpl struct {
	db *sqlx.DB
}

// NewTeamService creates a new TeamService instance
func NewTeamService(db *sqlx.DB) TeamService {
	return &teamServiceImpl{db: db}
}

// CreateTeam creates a new team
func (s *teamServiceImpl) CreateTeam(team *models.Team) (*models.Team, error) {
	if err := s.ValidateTeam(team); err != nil {
		return nil, err
	}

	now := time.Now()
	team.CreatedAt = now
	team.UpdatedAt = now

	var id int
	err := s.db.QueryRow(`
		INSERT INTO teams (name, external_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, team.Name, team.ExternalId, team.CreatedAt, team.UpdatedAt).Scan(&id)
	if err != nil {
		return nil, err
	}

	team.ID = id
	return team, nil
}

// GetTeam retrieves a team by ID
func (s *teamServiceImpl) GetTeam(id int64) (*models.Team, error) {
	team := &models.Team{}
	err := s.db.Get(team, "SELECT * FROM teams WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return team, nil
}

// UpdateTeam updates an existing team in the database
func (s *teamServiceImpl) UpdateTeam(team *models.Team) (*models.Team, error) {
	// Validate team data
	if err := s.ValidateTeam(team); err != nil {
		return nil, err
	}

	// Check if team exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM teams WHERE id = $1", team.ID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("team with ID %d not found", team.ID)
	}

	// Set updated timestamp
	team.UpdatedAt = time.Now()

	// Update team in database
	_, err = s.db.Exec(`
		UPDATE teams 
		SET name = $1, updated_at = $2
		WHERE id = $3
	`, team.Name, team.UpdatedAt, team.ID)
	if err != nil {
		return nil, err
	}

	// Return the updated team
	return team, nil
}

// DeleteTeam deletes a team by ID
func (s *teamServiceImpl) DeleteTeam(id int64) error {
	result, err := s.db.Exec("DELETE FROM teams WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("team with ID %d not found", id)
	}

	return nil
}

// ListTeams retrieves all teams
func (s *teamServiceImpl) ListTeams() ([]*models.Team, error) {
	var teams []*models.Team
	err := s.db.Select(&teams, "SELECT * FROM teams")
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// ValidateTeam validates team data
func (s *teamServiceImpl) ValidateTeam(team *models.Team) error {
	if team.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// GetTeamByExternalID retrieves a team by its external ID
func (s *teamServiceImpl) GetTeamByExternalID(externalID int64) (*models.Team, error) {
	team := &models.Team{}
	err := s.db.Get(team, "SELECT * FROM teams WHERE external_id = $1", externalID)
	if err != nil {
		return nil, err
	}
	return team, nil
}
