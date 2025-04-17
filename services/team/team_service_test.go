package team

import (
	"fmt"
	"testing"

	"go-app/database"
	"go-app/models"

	"github.com/stretchr/testify/assert"
)

var (
	testDB      *database.TestDB
	teamService TeamService
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

	teamService = NewTeamService(testDB.GetDB())
	m.Run()
}

func TestTeamService(t *testing.T) {
	// Test CreateTeam
	t.Run("CreateTeam", func(t *testing.T) {
		defer testDB.Clear()

		team := &models.Team{
			Name:       "Test Team",
			ExternalId: 123,
		}

		// Test successful creation
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)
		assert.NotZero(t, createdTeam.ID)
		assert.Equal(t, team.Name, createdTeam.Name)
	})

	// Test GetTeam
	t.Run("GetTeam", func(t *testing.T) {
		defer testDB.Clear()

		// Create a test team
		team := &models.Team{
			Name: "Get Team",
		}
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)

		// Test successful retrieval
		retrievedTeam, err := teamService.GetTeam(int64(createdTeam.ID))
		assert.NoError(t, err)
		assert.Equal(t, createdTeam.ID, retrievedTeam.ID)
		assert.Equal(t, team.Name, retrievedTeam.Name)

		// Test non-existent team
		_, err = teamService.GetTeam(int64(999))
		assert.Error(t, err)
	})

	// Test UpdateTeam
	t.Run("UpdateTeam", func(t *testing.T) {
		defer testDB.Clear()

		// Create a test team
		team := &models.Team{
			Name: "Original Team",
		}
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)

		// Update team
		createdTeam.Name = "Updated Team"
		updatedTeam, err := teamService.UpdateTeam(createdTeam)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Team", updatedTeam.Name)

		// Verify update in database
		retrievedTeam, err := teamService.GetTeam(int64(createdTeam.ID))
		assert.NoError(t, err)
		assert.Equal(t, "Updated Team", retrievedTeam.Name)
	})

	// Test DeleteTeam
	t.Run("DeleteTeam", func(t *testing.T) {
		defer testDB.Clear()

		// Create a test team
		team := &models.Team{
			Name: "Delete Team",
		}
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)

		// Delete team
		err = teamService.DeleteTeam(int64(createdTeam.ID))
		assert.NoError(t, err)

		// Verify deletion
		_, err = teamService.GetTeam(int64(createdTeam.ID))
		assert.Error(t, err)
	})

	// Test ListTeams
	t.Run("ListTeams", func(t *testing.T) {
		defer testDB.Clear()

		// Create multiple teams
		teams := []*models.Team{
			{
				Name:       "Arsenal",
				ExternalId: 1,
			},
			{
				Name:       "Manchester United",
				ExternalId: 2,
			},
			{
				Name:       "Manchester City",
				ExternalId: 3,
			},
		}

		for _, team := range teams {
			_, err := teamService.CreateTeam(team)
			assert.NoError(t, err)
		}

		// Test listing all teams
		allTeams, err := teamService.ListTeams()
		assert.NoError(t, err)
		assert.Len(t, allTeams, 3)
	})

	// Test ValidateTeam
	t.Run("ValidateTeam", func(t *testing.T) {
		defer testDB.Clear()

		// Test valid team
		validTeam := &models.Team{
			Name: "Valid Team",
		}
		err := teamService.ValidateTeam(validTeam)
		assert.NoError(t, err)

		// Test missing name
		invalidTeam := &models.Team{}
		err = teamService.ValidateTeam(invalidTeam)
		assert.Error(t, err)
	})

	// Test GetTeamByExternalID
	t.Run("GetTeamByExternalID", func(t *testing.T) {
		defer testDB.Clear()

		// Create a test team with external ID
		team := &models.Team{
			Name:       "External Team",
			ExternalId: 123,
		}
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)

		// Test successful retrieval
		retrievedTeam, err := teamService.GetTeamByExternalID(123)
		assert.NoError(t, err)
		assert.Equal(t, createdTeam.ID, retrievedTeam.ID)
		assert.Equal(t, team.Name, retrievedTeam.Name)
		assert.Equal(t, team.ExternalId, retrievedTeam.ExternalId)

		// Test non-existent team
		_, err = teamService.GetTeamByExternalID(999)
		assert.Error(t, err)
	})
}
