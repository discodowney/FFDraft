package player

import (
	"fmt"
	"testing"

	"go-app/database"
	"go-app/models"
	"go-app/services/team"

	"github.com/stretchr/testify/assert"
)

var testDB *database.TestDB
var playerService PlayerService
var teamService team.TeamService

func TestMain(m *testing.M) {
	// Setup
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

	playerService = NewPlayerService(testDB.GetDB())
	teamService = team.NewTeamService(testDB.GetDB())

	// Run tests
	m.Run()
}

func TestPlayerService(t *testing.T) {
	// Test CreatePlayer
	t.Run("CreatePlayer", func(t *testing.T) {
		defer testDB.Clear()

		// First create a team
		team := &models.Team{
			Name: "Test Team",
		}
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)

		player := &models.Player{
			TeamID:    createdTeam.ID,
			FirstName: "Test",
			LastName:  "Player",
			Position:  models.PositionFWD,
		}

		// Test successful creation
		createdPlayer, err := playerService.CreatePlayer(player)
		assert.NoError(t, err)
		assert.NotZero(t, createdPlayer.ID)
		assert.Equal(t, player.TeamID, createdPlayer.TeamID)
		assert.Equal(t, player.FirstName, createdPlayer.FirstName)
		assert.Equal(t, player.LastName, createdPlayer.LastName)
		assert.Equal(t, player.Position, createdPlayer.Position)
	})

	// Test GetPlayer
	t.Run("GetPlayer", func(t *testing.T) {
		defer testDB.Clear()

		// First create a team
		team := &models.Team{
			Name: "Get Team",
		}
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)

		// Create a test player
		player := &models.Player{
			TeamID:    createdTeam.ID,
			FirstName: "Get",
			LastName:  "Player",
			Position:  models.PositionMID,
		}
		createdPlayer, err := playerService.CreatePlayer(player)
		assert.NoError(t, err)

		// Test successful retrieval
		retrievedPlayer, err := playerService.GetPlayer(createdPlayer.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdPlayer.ID, retrievedPlayer.ID)
		assert.Equal(t, player.TeamID, retrievedPlayer.TeamID)
		assert.Equal(t, player.FirstName, retrievedPlayer.FirstName)
		assert.Equal(t, player.LastName, retrievedPlayer.LastName)
		assert.Equal(t, player.Position, retrievedPlayer.Position)

		// Test non-existent player
		_, err = playerService.GetPlayer(999)
		assert.Error(t, err)
	})

	// Test UpdatePlayer
	t.Run("UpdatePlayer", func(t *testing.T) {
		defer testDB.Clear()

		// First create a team
		team := &models.Team{
			Name: "Update Team",
		}
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)

		// Create a test player
		player := &models.Player{
			TeamID:    createdTeam.ID,
			FirstName: "Original",
			LastName:  "Player",
			Position:  models.PositionDEF,
		}
		createdPlayer, err := playerService.CreatePlayer(player)
		assert.NoError(t, err)

		// Update player
		createdPlayer.FirstName = "Updated"
		createdPlayer.LastName = "Name"
		createdPlayer.Position = models.PositionGK
		updatedPlayer, err := playerService.UpdatePlayer(createdPlayer)
		assert.NoError(t, err)
		assert.Equal(t, "Updated", updatedPlayer.FirstName)
		assert.Equal(t, "Name", updatedPlayer.LastName)
		assert.Equal(t, models.PositionGK, updatedPlayer.Position)

		// Verify update in database
		retrievedPlayer, err := playerService.GetPlayer(createdPlayer.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated", retrievedPlayer.FirstName)
		assert.Equal(t, "Name", retrievedPlayer.LastName)
		assert.Equal(t, models.PositionGK, retrievedPlayer.Position)
	})

	// Test DeletePlayer
	t.Run("DeletePlayer", func(t *testing.T) {
		defer testDB.Clear()

		// First create a team
		team := &models.Team{
			Name: "Delete Team",
		}
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)

		// Create a test player
		player := &models.Player{
			TeamID:    createdTeam.ID,
			FirstName: "ToDelete",
			LastName:  "Player",
			Position:  models.PositionMID,
		}
		createdPlayer, err := playerService.CreatePlayer(player)
		assert.NoError(t, err)

		// Delete player
		err = playerService.DeletePlayer(createdPlayer.ID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = playerService.GetPlayer(createdPlayer.ID)
		assert.Error(t, err)
	})

	// Test ListPlayers
	t.Run("ListPlayers", func(t *testing.T) {
		defer testDB.Clear()

		// First create a team
		team := &models.Team{
			Name: "List Team",
		}
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)

		// Create test players
		players := []*models.Player{
			{
				TeamID:    createdTeam.ID,
				FirstName: "John",
				LastName:  "Doe",
				Position:  models.PositionFWD,
			},
			{
				TeamID:    createdTeam.ID,
				FirstName: "Jane",
				LastName:  "Smith",
				Position:  models.PositionMID,
			},
			{
				TeamID:    createdTeam.ID,
				FirstName: "Bob",
				LastName:  "Johnson",
				Position:  models.PositionDEF,
			},
		}

		for _, player := range players {
			_, err := playerService.CreatePlayer(player)
			assert.NoError(t, err)
		}

		// Test listing all players
		allPlayers, err := playerService.ListPlayers(nil)
		assert.NoError(t, err)
		assert.Len(t, allPlayers, 3)

		// Test filtering by position
		filteredPlayers, err := playerService.ListPlayers(&PlayerFilter{
			Position: models.PositionFWD,
		})
		assert.NoError(t, err)
		assert.Len(t, filteredPlayers, 1)
		assert.Equal(t, "John", filteredPlayers[0].FirstName)

		// Test filtering by first name
		filteredPlayers, err = playerService.ListPlayers(&PlayerFilter{
			FirstName: "Jane",
		})
		assert.NoError(t, err)
		assert.Len(t, filteredPlayers, 1)
		assert.Equal(t, "Smith", filteredPlayers[0].LastName)

		// Test filtering by last name
		filteredPlayers, err = playerService.ListPlayers(&PlayerFilter{
			LastName: "Johnson",
		})
		assert.NoError(t, err)
		assert.Len(t, filteredPlayers, 1)
		assert.Equal(t, "Bob", filteredPlayers[0].FirstName)

		// Test filtering by team ID
		filteredPlayers, err = playerService.ListPlayers(&PlayerFilter{
			TeamID: createdTeam.ID,
		})
		assert.NoError(t, err)
		assert.Len(t, filteredPlayers, 3)

		// Test multiple filters
		filteredPlayers, err = playerService.ListPlayers(&PlayerFilter{
			Position:  models.PositionMID,
			FirstName: "Jane",
		})
		assert.NoError(t, err)
		assert.Len(t, filteredPlayers, 1)
		assert.Equal(t, "Smith", filteredPlayers[0].LastName)
	})

	// Test ListPlayers with filters
	t.Run("ListPlayers with filters", func(t *testing.T) {
		defer testDB.Clear()

		// First create a team
		team := &models.Team{
			Name: "Filter Team",
		}
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)

		// Create test players
		players := []*models.Player{
			{
				TeamID:    createdTeam.ID,
				FirstName: "Filter",
				LastName:  "One",
				Position:  models.PositionFWD,
			},
			{
				TeamID:    createdTeam.ID,
				FirstName: "Filter",
				LastName:  "Two",
				Position:  models.PositionMID,
			},
			{
				TeamID:    createdTeam.ID,
				FirstName: "Other",
				LastName:  "Three",
				Position:  models.PositionDEF,
			},
		}

		for _, player := range players {
			_, err := playerService.CreatePlayer(player)
			assert.NoError(t, err)
		}

		// Test position filter
		filter := &PlayerFilter{
			Position: models.PositionFWD,
		}
		filteredPlayers, err := playerService.ListPlayers(filter)
		assert.NoError(t, err)
		assert.Len(t, filteredPlayers, 1)
		assert.Equal(t, models.PositionFWD, filteredPlayers[0].Position)

		// Test first name filter
		filter = &PlayerFilter{
			FirstName: "Filter",
		}
		filteredPlayers, err = playerService.ListPlayers(filter)
		assert.NoError(t, err)
		assert.Len(t, filteredPlayers, 2)
		for _, player := range filteredPlayers {
			assert.Equal(t, "Filter", player.FirstName)
		}

		// Test last name filter
		filter = &PlayerFilter{
			LastName: "Three",
		}
		filteredPlayers, err = playerService.ListPlayers(filter)
		assert.NoError(t, err)
		assert.Len(t, filteredPlayers, 1)
		assert.Equal(t, "Three", filteredPlayers[0].LastName)

		// Test team ID filter
		filter = &PlayerFilter{
			TeamID: createdTeam.ID,
		}
		filteredPlayers, err = playerService.ListPlayers(filter)
		assert.NoError(t, err)
		assert.Len(t, filteredPlayers, 3)
		for _, player := range filteredPlayers {
			assert.Equal(t, createdTeam.ID, player.TeamID)
		}

		// Test multiple filters
		filter = &PlayerFilter{
			TeamID:    createdTeam.ID,
			Position:  models.PositionFWD,
			FirstName: "Filter",
		}
		filteredPlayers, err = playerService.ListPlayers(filter)
		assert.NoError(t, err)
		assert.Len(t, filteredPlayers, 1)
		assert.Equal(t, createdTeam.ID, filteredPlayers[0].TeamID)
		assert.Equal(t, models.PositionFWD, filteredPlayers[0].Position)
		assert.Equal(t, "Filter", filteredPlayers[0].FirstName)
	})

	// Test GetPlayersByTeam
	t.Run("GetPlayersByTeam", func(t *testing.T) {
		defer testDB.Clear()

		// Create two teams
		team1 := &models.Team{
			Name: "Team 1",
		}
		createdTeam1, err := teamService.CreateTeam(team1)
		assert.NoError(t, err)

		team2 := &models.Team{
			Name: "Team 2",
		}
		createdTeam2, err := teamService.CreateTeam(team2)
		assert.NoError(t, err)

		// Create players for each team
		players := []*models.Player{
			{
				TeamID:    createdTeam1.ID,
				FirstName: "Team1",
				LastName:  "Player1",
				Position:  models.PositionFWD,
			},
			{
				TeamID:    createdTeam1.ID,
				FirstName: "Team1",
				LastName:  "Player2",
				Position:  models.PositionMID,
			},
			{
				TeamID:    createdTeam2.ID,
				FirstName: "Team2",
				LastName:  "Player1",
				Position:  models.PositionDEF,
			},
		}

		for _, player := range players {
			_, err := playerService.CreatePlayer(player)
			assert.NoError(t, err)
		}

		// Get players by team
		team1Players, err := playerService.GetPlayersByTeam(createdTeam1.ID)
		assert.NoError(t, err)
		assert.Len(t, team1Players, 2)

		team2Players, err := playerService.GetPlayersByTeam(createdTeam2.ID)
		assert.NoError(t, err)
		assert.Len(t, team2Players, 1)
	})

	// Test ValidatePlayer
	t.Run("ValidatePlayer", func(t *testing.T) {
		defer testDB.Clear()

		// First create a team
		team := &models.Team{
			Name: "Validate Team",
		}
		createdTeam, err := teamService.CreateTeam(team)
		assert.NoError(t, err)

		// Test valid player
		validPlayer := &models.Player{
			TeamID:    createdTeam.ID,
			FirstName: "Valid",
			LastName:  "Player",
			Position:  models.PositionFWD,
		}
		err = playerService.ValidatePlayer(validPlayer)
		assert.NoError(t, err)

		// Test missing team ID
		invalidPlayer := &models.Player{
			FirstName: "Player",
			LastName:  "Test",
			Position:  models.PositionFWD,
		}
		err = playerService.ValidatePlayer(invalidPlayer)
		assert.Error(t, err)

		// Test missing first name
		invalidPlayer = &models.Player{
			TeamID:   createdTeam.ID,
			LastName: "Test",
			Position: models.PositionFWD,
		}
		err = playerService.ValidatePlayer(invalidPlayer)
		assert.Error(t, err)

		// Test missing last name
		invalidPlayer = &models.Player{
			TeamID:    createdTeam.ID,
			FirstName: "Player",
			Position:  models.PositionFWD,
		}
		err = playerService.ValidatePlayer(invalidPlayer)
		assert.Error(t, err)

		// Test missing position
		invalidPlayer = &models.Player{
			TeamID:    createdTeam.ID,
			FirstName: "Player",
			LastName:  "Test",
		}
		err = playerService.ValidatePlayer(invalidPlayer)
		assert.Error(t, err)
	})

	// Test ValidatePosition
	t.Run("ValidatePosition", func(t *testing.T) {
		// Test valid positions
		err := playerService.ValidatePosition(models.PositionGK)
		assert.NoError(t, err)

		err = playerService.ValidatePosition(models.PositionDEF)
		assert.NoError(t, err)

		err = playerService.ValidatePosition(models.PositionMID)
		assert.NoError(t, err)

		err = playerService.ValidatePosition(models.PositionFWD)
		assert.NoError(t, err)

		// Test invalid position
		err = playerService.ValidatePosition("INVALID")
		assert.Error(t, err)
	})

	// Test GetPlayerStats - leave commented out until we have a player stats service to add player stats data
	// t.Run("GetPlayerStats", func(t *testing.T) {
	// 	defer testDB.Clear()

	// 	// First create a team
	// 	team := &models.Team{
	// 		Name: "Stats Team",
	// 	}
	// 	createdTeam, err := teamService.CreateTeam(team)
	// 	assert.NoError(t, err)

	// 	// Create a test player
	// 	player := &models.Player{
	// 		TeamID:    createdTeam.ID,
	// 		FirstName: "Stats",
	// 		LastName:  "Player",
	// 		Position:  models.PositionFWD,
	// 	}
	// 	createdPlayer, err := playerService.CreatePlayer(player)
	// 	assert.NoError(t, err)

	// 	// Test successful retrieval
	// 	stats, err := playerService.GetPlayerStats(createdPlayer.ID)
	// 	assert.NoError(t, err)
	// 	assert.NotNil(t, stats)

	// 	// Test non-existent player
	// 	stats, err = playerService.GetPlayerStats(999)
	// 	assert.NoError(t, err)
	// 	assert.Nil(t, stats)
	// })
}
