package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// TestDB represents a test database instance
type TestDB struct {
	db *sqlx.DB
}

// NewTestDB creates a new test database instance
func NewTestDB() (*TestDB, error) {
	// Get test database URL from environment variable or use default
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/fantasy_football_test?sslmode=disable"
	}

	// First try to connect to the database
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		// If connection fails, try to create the database
		// Extract database name from URL
		dbName := "fantasy_football_test"
		// Connect to default postgres database
		defaultDB, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
		if err != nil {
			return nil, fmt.Errorf("failed to connect to default database: %v", err)
		}
		defer defaultDB.Close()

		// Create the test database
		_, err = defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return nil, fmt.Errorf("failed to create test database: %v", err)
		}

		// Now try to connect to the newly created database
		db, err = sqlx.Connect("postgres", dbURL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to test database after creation: %v", err)
		}
	}

	// Drop existing tables to ensure clean state
	if err := dropTestTables(db); err != nil {
		return nil, fmt.Errorf("failed to drop existing tables: %v", err)
	}

	// Create test tables
	if err := createTestTables(db); err != nil {
		return nil, fmt.Errorf("failed to create test tables: %v", err)
	}

	return &TestDB{db: db}, nil
}

// GetDB returns the database connection
func (t *TestDB) GetDB() *sqlx.DB {
	return t.db
}

// Close closes the database connection and cleans up test data
func (t *TestDB) Close() error {
	if err := dropTestTables(t.db); err != nil {
		return fmt.Errorf("failed to drop test tables: %v", err)
	}
	if err := t.db.Close(); err != nil {
		return fmt.Errorf("failed to close test database: %v", err)
	}
	return nil
}

// createTestTables creates all necessary tables for testing
func createTestTables(db *sqlx.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	// Create teams table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS teams (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			external_id INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create teams table: %v", err)
	}

	// Create players table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS players (
			id SERIAL PRIMARY KEY,
			team_id INTEGER REFERENCES teams(id),
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			position VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create players table: %v", err)
	}

	// Create user_teams table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_teams (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create user_teams table: %v", err)
	}

	// Create user_team_players table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_team_players (
			id SERIAL PRIMARY KEY,
			user_team_id INTEGER REFERENCES user_teams(id),
			player_id INTEGER REFERENCES players(id),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create user_team_players table: %v", err)
	}

	// Create player_stats table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS player_stats (
			id SERIAL PRIMARY KEY,
			player_id INTEGER REFERENCES players(id),
			goals INTEGER DEFAULT 0,
			assists INTEGER DEFAULT 0,
			yellow_cards INTEGER DEFAULT 0,
			red_cards INTEGER DEFAULT 0,
			clean_sheets INTEGER DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create player_stats table: %v", err)
	}

	// Create matches table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS matches (
			id SERIAL PRIMARY KEY,
			league_id INTEGER,
			home_team_id INTEGER REFERENCES teams(id),
			away_team_id INTEGER REFERENCES teams(id),
			match_date TIMESTAMP NOT NULL,
			home_score INTEGER DEFAULT 0,
			away_score INTEGER DEFAULT 0,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create matches table: %v", err)
	}

	// Create match_incidents table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS match_incidents (
			id SERIAL PRIMARY KEY,
			match_id INTEGER REFERENCES matches(id),
			player_id INTEGER REFERENCES players(id),
			type VARCHAR(50) NOT NULL,
			minute INTEGER NOT NULL,
			description TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create match_incidents table: %v", err)
	}

	return nil
}

// dropTestTables drops all test tables
func dropTestTables(db *sqlx.DB) error {
	tables := []string{
		"match_incidents",
		"matches",
		"player_stats",
		"user_team_players",
		"user_teams",
		"players",
		"teams",
		"users",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to drop table %s: %v", table, err)
		}
	}

	return nil
}

// Clear removes all data from the test database
func (t *TestDB) Clear() error {
	tables := []string{
		"match_incidents",
		"matches",
		"player_stats",
		"user_team_players",
		"user_teams",
		"players",
		"teams",
		"users",
	}

	for _, table := range tables {
		_, err := t.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to clear table %s: %v", table, err)
		}
	}

	return nil
}
