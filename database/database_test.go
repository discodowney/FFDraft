package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testDB *TestDB

// TestMain runs before and after all tests in the package
func TestMain(m *testing.M) {
	// Setup
	var err error
	testDB, err = NewTestDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to create test database: %v", err))
	}

	defer func() {
		if err := testDB.Close(); err != nil {
			panic(fmt.Sprintf("Failed to close test database: %v", err))
		}
	}()

	// Run tests
	m.Run()
}

func TestDatabaseConnection(t *testing.T) {
	// Clear the database at the end of the test, even if it fails
	defer testDB.Clear()

	db := testDB.GetDB()
	assert.NotNil(t, db)

	// Test the connection
	err := db.Ping()
	assert.NoError(t, err)
}

func TestDatabaseOperations(t *testing.T) {
	db := testDB.GetDB()

	// Test user operations
	t.Run("User Operations", func(t *testing.T) {
		// Clear the database at the end of the test, even if it fails
		defer testDB.Clear()

		// Create a user
		var userID int
		err := db.QueryRow(`
			INSERT INTO users (first_name, last_name, email, password, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, "John", "Doe", "john@example.com", "password", time.Now(), time.Now()).Scan(&userID)
		assert.NoError(t, err)
		assert.NotZero(t, userID)

		// Get the user
		var firstName, lastName, email string
		err = db.QueryRow(`
			SELECT first_name, last_name, email
			FROM users
			WHERE id = $1
		`, userID).Scan(&firstName, &lastName, &email)
		assert.NoError(t, err)
		assert.Equal(t, "John", firstName)
		assert.Equal(t, "Doe", lastName)
		assert.Equal(t, "john@example.com", email)
	})

	// Test team operations
	t.Run("Team Operations", func(t *testing.T) {
		// Clear the database at the end of the test, even if it fails
		defer testDB.Clear()

		// Create a team
		var teamID int
		err := db.QueryRow(`
			INSERT INTO teams (name, created_at, updated_at)
			VALUES ($1, $2, $3)
			RETURNING id
		`, "Test Team", time.Now(), time.Now()).Scan(&teamID)
		assert.NoError(t, err)
		assert.NotZero(t, teamID)

		// Get the team
		var name string
		err = db.QueryRow(`
			SELECT name
			FROM teams
			WHERE id = $1
		`, teamID).Scan(&name)
		assert.NoError(t, err)
		assert.Equal(t, "Test Team", name)
	})

	// Test player operations
	t.Run("Player Operations", func(t *testing.T) {
		// Clear the database at the end of the test, even if it fails
		defer testDB.Clear()

		// First create a team
		var teamID int
		err := db.QueryRow(`
			INSERT INTO teams (name, country, created_at, updated_at)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`, "Player Team", "Player Country", time.Now(), time.Now()).Scan(&teamID)
		assert.NoError(t, err)

		// Create a player
		var playerID int
		err = db.QueryRow(`
			INSERT INTO players (team_id, first_name, last_name, position, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, teamID, "Test", "Player", "Forward", time.Now(), time.Now()).Scan(&playerID)
		assert.NoError(t, err)
		assert.NotZero(t, playerID)

		// Get the player
		var firstName, lastName, position string
		err = db.QueryRow(`
			SELECT first_name, last_name, position
			FROM players
			WHERE id = $1
		`, playerID).Scan(&firstName, &lastName, &position)
		assert.NoError(t, err)
		assert.Equal(t, "Test", firstName)
		assert.Equal(t, "Player", lastName)
		assert.Equal(t, "Forward", position)
	})
}
