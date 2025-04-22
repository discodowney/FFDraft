package league

import (
	"fmt"
	"testing"

	"go-app/database"
	"go-app/models"

	"github.com/stretchr/testify/assert"
)

var (
	testDB        *database.TestDB
	leagueService LeagueService
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = database.NewTestDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to create test database: %v", err))
	}
	defer func() {
		if err := testDB.Close(); err != nil {
			panic(fmt.Sprintf("Failed to close test database: %v", err))
		}
	}()

	leagueService = NewLeagueService(testDB.GetDB())
	m.Run()
}

func TestLeagueService(t *testing.T) {
	// Test CreateLeague
	t.Run("CreateLeague", func(t *testing.T) {
		defer testDB.Clear()

		league := &models.League{
			Name: "Test League",
			Code: "TEST123",
		}

		// Test successful creation
		createdLeague, err := leagueService.CreateLeague(league)
		assert.NoError(t, err)
		assert.NotZero(t, createdLeague.ID)
		assert.Equal(t, league.Name, createdLeague.Name)
		assert.Equal(t, league.Code, createdLeague.Code)

		// Test duplicate code
		duplicateLeague := &models.League{
			Name: "Duplicate League",
			Code: "TEST123",
		}
		_, err = leagueService.CreateLeague(duplicateLeague)
		assert.Error(t, err)
	})

	// Test GetLeague
	t.Run("GetLeague", func(t *testing.T) {
		defer testDB.Clear()

		// Create a test league
		league := &models.League{
			Name: "Get League",
			Code: "GET123",
		}
		createdLeague, err := leagueService.CreateLeague(league)
		assert.NoError(t, err)

		// Test successful retrieval
		retrievedLeague, err := leagueService.GetLeague(createdLeague.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdLeague.ID, retrievedLeague.ID)
		assert.Equal(t, league.Name, retrievedLeague.Name)
		assert.Equal(t, league.Code, retrievedLeague.Code)

		// Test non-existent league
		_, err = leagueService.GetLeague(999)
		assert.Error(t, err)
	})

	// Test UpdateLeague
	t.Run("UpdateLeague", func(t *testing.T) {
		defer testDB.Clear()

		// Create a test league
		league := &models.League{
			Name: "Original League",
			Code: "ORG123",
		}
		createdLeague, err := leagueService.CreateLeague(league)
		assert.NoError(t, err)

		// Update league
		createdLeague.Name = "Updated League"
		createdLeague.Code = "UPD123"
		updatedLeague, err := leagueService.UpdateLeague(createdLeague)
		assert.NoError(t, err)
		assert.Equal(t, "Updated League", updatedLeague.Name)
		assert.Equal(t, "UPD123", updatedLeague.Code)

		// Verify update in database
		retrievedLeague, err := leagueService.GetLeague(createdLeague.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated League", retrievedLeague.Name)
		assert.Equal(t, "UPD123", retrievedLeague.Code)
	})

	// Test DeleteLeague
	t.Run("DeleteLeague", func(t *testing.T) {
		defer testDB.Clear()

		// Create a test league
		league := &models.League{
			Name: "Delete League",
			Code: "DEL123",
		}
		createdLeague, err := leagueService.CreateLeague(league)
		assert.NoError(t, err)

		// Delete league
		err = leagueService.DeleteLeague(createdLeague.ID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = leagueService.GetLeague(createdLeague.ID)
		assert.Error(t, err)
	})

	// Test ListLeagues
	t.Run("ListLeagues", func(t *testing.T) {
		defer testDB.Clear()

		// Create multiple leagues
		leagues := []*models.League{
			{
				Name: "Premier League",
				Code: "PL",
			},
			{
				Name: "La Liga",
				Code: "LL",
			},
			{
				Name: "Bundesliga",
				Code: "BL",
			},
		}

		for _, league := range leagues {
			_, err := leagueService.CreateLeague(league)
			assert.NoError(t, err)
		}

		// Test listing all leagues
		allLeagues, err := leagueService.ListLeagues()
		assert.NoError(t, err)
		assert.Len(t, allLeagues, 3)
	})

	// Test ValidateLeague
	t.Run("ValidateLeague", func(t *testing.T) {
		defer testDB.Clear()

		// Test valid league
		validLeague := &models.League{
			Name: "Valid League",
			Code: "VL123",
		}
		err := leagueService.ValidateLeague(validLeague)
		assert.NoError(t, err)

		// Test missing name
		invalidLeague := &models.League{
			Code: "INV123",
		}
		err = leagueService.ValidateLeague(invalidLeague)
		assert.Error(t, err)

		// Test missing code
		invalidLeague = &models.League{
			Name: "Invalid League",
		}
		err = leagueService.ValidateLeague(invalidLeague)
		assert.Error(t, err)
	})

	// Test GetLeagueByCode
	t.Run("GetLeagueByCode", func(t *testing.T) {
		defer testDB.Clear()

		// Create a test league with code
		league := &models.League{
			Name: "Code League",
			Code: "CODE123",
		}
		createdLeague, err := leagueService.CreateLeague(league)
		assert.NoError(t, err)

		// Test successful retrieval
		retrievedLeague, err := leagueService.GetLeagueByCode("CODE123")
		assert.NoError(t, err)
		assert.NotNil(t, retrievedLeague)
		assert.Equal(t, createdLeague.ID, retrievedLeague.ID)
		assert.Equal(t, league.Name, retrievedLeague.Name)
		assert.Equal(t, league.Code, retrievedLeague.Code)

		// Test non-existent league
		nonExistentLeague, err := leagueService.GetLeagueByCode("NONEXIST")
		assert.NoError(t, err)
		assert.Nil(t, nonExistentLeague)
	})
}
