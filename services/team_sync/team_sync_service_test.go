package team_sync

import (
	"fmt"
	"testing"

	"go-app/database"
	"go-app/mocks"
	"go-app/models"
	"go-app/services/team"

	"github.com/stretchr/testify/assert"
)

var (
	testDB          *database.TestDB
	teamService     team.TeamService
	teamSyncService *TeamSyncService
	mockClient      *mocks.MockAPIFootballClient
)

func TestMain(m *testing.M) {
	// Initialize test database
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

	// Initialize services
	teamService = team.NewTeamService(testDB.GetDB())
	mockClient = mocks.NewMockAPIFootballClient(
		[]*models.Team{
			{
				Name:       "Team A",
				ExternalId: 1,
			},
			{
				Name:       "Team B",
				ExternalId: 2,
			},
		},
		nil,
	)
	teamSyncService = NewTeamSyncService(teamService, mockClient)

	// Run tests
	m.Run()
}

func TestTeamSyncService(t *testing.T) {
	t.Run("SyncTeams", func(t *testing.T) {
		// Sync teams
		err := teamSyncService.SyncTeams()
		assert.NoError(t, err)

		// Verify teams were synced
		syncedTeams, err := teamService.ListTeams()
		assert.NoError(t, err)
		assert.Len(t, syncedTeams, 2)
		assert.Equal(t, "Team A", syncedTeams[0].Name)
		assert.Equal(t, "Team B", syncedTeams[1].Name)
	})
}
