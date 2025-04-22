package league

import (
	"database/sql"
	"fmt"
	"time"

	"go-app/models"

	"github.com/jmoiron/sqlx"
)

// LeagueService defines the interface for league-related operations
type LeagueService interface {
	CreateLeague(league *models.League) (*models.League, error)
	GetLeague(id int) (*models.League, error)
	UpdateLeague(league *models.League) (*models.League, error)
	DeleteLeague(id int) error
	ListLeagues() ([]*models.League, error)
	ValidateLeague(league *models.League) error
	GetLeagueByCode(code string) (*models.League, error)
}

// Implementation of the LeagueService interface
type leagueServiceImpl struct {
	db *sqlx.DB
}

// NewLeagueService creates a new LeagueService instance
func NewLeagueService(db *sqlx.DB) LeagueService {
	return &leagueServiceImpl{db: db}
}

// CreateLeague creates a new league
func (s *leagueServiceImpl) CreateLeague(league *models.League) (*models.League, error) {
	if err := s.ValidateLeague(league); err != nil {
		return nil, err
	}

	now := time.Now()
	league.CreatedAt = now
	league.UpdatedAt = now

	var id int
	err := s.db.QueryRow(`
		INSERT INTO leagues (code, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, league.Code, league.Name, league.CreatedAt, league.UpdatedAt).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error creating league: %w", err)
	}

	league.ID = id
	return league, nil
}

// GetLeague retrieves a league by ID
func (s *leagueServiceImpl) GetLeague(id int) (*models.League, error) {
	league := &models.League{}
	err := s.db.Get(league, "SELECT * FROM leagues WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return league, nil
}

// UpdateLeague updates an existing league in the database
func (s *leagueServiceImpl) UpdateLeague(league *models.League) (*models.League, error) {
	// Validate league data
	if err := s.ValidateLeague(league); err != nil {
		return nil, err
	}

	// Check if league exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM leagues WHERE id = $1", league.ID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("league with ID %d not found", league.ID)
	}

	// Set updated timestamp
	league.UpdatedAt = time.Now()

	// Update league in database
	_, err = s.db.Exec(`
		UPDATE leagues 
		SET name = $1, code = $2, updated_at = $3
		WHERE id = $4
	`, league.Name, league.Code, league.UpdatedAt, league.ID)
	if err != nil {
		return nil, err
	}

	// Return the updated league
	return league, nil
}

// DeleteLeague deletes a league by ID
func (s *leagueServiceImpl) DeleteLeague(id int) error {
	result, err := s.db.Exec("DELETE FROM leagues WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("league with ID %d not found", id)
	}

	return nil
}

// ListLeagues retrieves all leagues
func (s *leagueServiceImpl) ListLeagues() ([]*models.League, error) {
	var leagues []*models.League
	err := s.db.Select(&leagues, "SELECT * FROM leagues")
	if err != nil {
		return nil, err
	}
	return leagues, nil
}

// ValidateLeague validates league data
func (s *leagueServiceImpl) ValidateLeague(league *models.League) error {
	if league.Name == "" {
		return fmt.Errorf("name is required")
	}
	if league.Code == "" {
		return fmt.Errorf("code is required")
	}
	return nil
}

// GetLeagueByCode retrieves a league by its code
func (s *leagueServiceImpl) GetLeagueByCode(code string) (*models.League, error) {
	league := &models.League{}
	err := s.db.Get(league, "SELECT * FROM leagues WHERE code = $1", code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return league, nil
}
