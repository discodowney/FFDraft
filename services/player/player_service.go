package player

import (
	"fmt"
	"time"

	"go-app/models"

	"github.com/jmoiron/sqlx"
)

// PlayerService defines the interface for player-related operations
type PlayerService interface {
	GetPlayer(id int) (*models.Player, error)
	CreatePlayer(player *models.Player) (*models.Player, error)
	UpdatePlayer(player *models.Player) (*models.Player, error)
	DeletePlayer(id int) error
	ListPlayers(filter *PlayerFilter) ([]*models.Player, error)
	GetPlayersByTeam(teamID int) ([]*models.Player, error)
	GetPlayerStats(playerID int) (*models.PlayerStats, error)
	ValidatePlayer(player *models.Player) error
	ValidatePosition(position models.Position) error
}

// PlayerFilter represents the filter criteria for listing players
type PlayerFilter struct {
	Position  models.Position
	TeamID    int
	FirstName string
	LastName  string
}

// Implementation of the PlayerService interface
type playerServiceImpl struct {
	db *sqlx.DB
}

// NewPlayerService creates a new PlayerService instance
func NewPlayerService(db *sqlx.DB) *playerServiceImpl {
	return &playerServiceImpl{db: db}
}

// GetPlayer retrieves a player by ID
func (s *playerServiceImpl) GetPlayer(id int) (*models.Player, error) {
	player := &models.Player{}
	err := s.db.Get(player, "SELECT * FROM players WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return player, nil
}

// CreatePlayer creates a new player in the database
func (s *playerServiceImpl) CreatePlayer(player *models.Player) (*models.Player, error) {
	// Validate player data
	if err := s.ValidatePlayer(player); err != nil {
		return nil, err
	}

	// Set timestamps
	now := time.Now()
	player.CreatedAt = now
	player.UpdatedAt = now

	// Insert player into database
	var id int
	err := s.db.QueryRow(`
		INSERT INTO players (team_id, first_name, last_name, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, player.TeamID, player.FirstName, player.LastName, player.Position, player.CreatedAt, player.UpdatedAt).Scan(&id)
	if err != nil {
		return nil, err
	}

	// Set the ID and return the player
	player.ID = id
	return player, nil
}

// UpdatePlayer updates an existing player in the database
func (s *playerServiceImpl) UpdatePlayer(player *models.Player) (*models.Player, error) {
	// Validate player data
	if err := s.ValidatePlayer(player); err != nil {
		return nil, err
	}

	// Check if player exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM players WHERE id = $1", player.ID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("player with ID %d not found", player.ID)
	}

	// Set updated timestamp
	player.UpdatedAt = time.Now()

	// Update player in database
	_, err = s.db.Exec(`
		UPDATE players 
		SET team_id = $1, first_name = $2, last_name = $3, position = $4, updated_at = $5
		WHERE id = $6
	`, player.TeamID, player.FirstName, player.LastName, player.Position, player.UpdatedAt, player.ID)
	if err != nil {
		return nil, err
	}

	// Return the updated player
	return player, nil
}

// DeletePlayer deletes a player by ID
func (s *playerServiceImpl) DeletePlayer(id int) error {
	result, err := s.db.Exec("DELETE FROM players WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("player with ID %d not found", id)
	}

	return nil
}

// ListPlayers retrieves all players with optional filtering
func (s *playerServiceImpl) ListPlayers(filter *PlayerFilter) ([]*models.Player, error) {
	query := "SELECT * FROM players WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.TeamID != 0 {
			query += fmt.Sprintf(" AND team_id = $%d", argCount)
			args = append(args, filter.TeamID)
			argCount++
		}
		if filter.FirstName != "" {
			query += fmt.Sprintf(" AND first_name ILIKE $%d", argCount)
			args = append(args, "%"+filter.FirstName+"%")
			argCount++
		}
		if filter.LastName != "" {
			query += fmt.Sprintf(" AND last_name ILIKE $%d", argCount)
			args = append(args, "%"+filter.LastName+"%")
			argCount++
		}
		if filter.Position != "" {
			query += fmt.Sprintf(" AND position = $%d", argCount)
			args = append(args, filter.Position)
			argCount++
		}
	}

	var players []*models.Player
	err := s.db.Select(&players, query, args...)
	if err != nil {
		return nil, err
	}
	return players, nil
}

// GetPlayersByTeam retrieves all players for a specific team
func (s *playerServiceImpl) GetPlayersByTeam(teamID int) ([]*models.Player, error) {
	var players []*models.Player
	err := s.db.Select(&players, "SELECT * FROM players WHERE team_id = $1", teamID)
	if err != nil {
		return nil, err
	}
	return players, nil
}

// ValidatePlayer validates player data
func (s *playerServiceImpl) ValidatePlayer(player *models.Player) error {
	if player.TeamID == 0 {
		return fmt.Errorf("team ID is required")
	}
	if player.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if player.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if player.Position == "" {
		return fmt.Errorf("position is required")
	}
	return nil
}

// ValidatePosition validates a player position
func (s *playerServiceImpl) ValidatePosition(position models.Position) error {
	switch position {
	case models.PositionGK, models.PositionDEF, models.PositionMID, models.PositionFWD:
		return nil
	default:
		return fmt.Errorf("invalid position: %s", position)
	}
}

// GetPlayerStats retrieves a player's statistics
func (s *playerServiceImpl) GetPlayerStats(playerID int) (*models.PlayerStats, error) {
	stats := &models.PlayerStats{}
	err := s.db.Get(stats, "SELECT * FROM player_stats WHERE player_id = $1", playerID)
	if err != nil {
		return nil, nil
	}
	return stats, nil
}
