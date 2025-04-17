package player

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-app/models"
	"go-app/server/handlers/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupPlayerHandlerTest(t *testing.T) (*gin.Engine, *mocks.MockPlayerService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(mocks.MockPlayerService)
	handler := &PlayerHandler{
		playerService: mockService,
	}

	// Setup routes
	router.GET("/players/:id", handler.GetPlayer)
	router.GET("/players", handler.ListPlayers)
	router.GET("/teams/:teamId/players", handler.GetPlayersByTeam)
	router.GET("/players/:id/stats", handler.GetPlayerStats)

	return router, mockService
}

func TestGetPlayer(t *testing.T) {
	router, mockService := setupPlayerHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		expectedPlayer := &models.Player{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Position:  "MID",
			TeamID:    1,
		}

		mockService.On("GetPlayer", 1).Return(expectedPlayer, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/players/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Player
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedPlayer.ID, response.ID)
		assert.Equal(t, expectedPlayer.FirstName, response.FirstName)
		assert.Equal(t, expectedPlayer.LastName, response.LastName)
	})

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/players/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("GetPlayer", 999).Return(nil, sql.ErrNoRows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/players/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestListPlayers(t *testing.T) {
	router, mockService := setupPlayerHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		expectedPlayers := []*models.Player{
			{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Position:  "MID",
				TeamID:    1,
			},
			{
				ID:        2,
				FirstName: "Jane",
				LastName:  "Smith",
				Position:  "FWD",
				TeamID:    1,
			},
		}

		mockService.On("ListPlayers", mock.Anything).Return(expectedPlayers, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/players", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.Player
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response, 2)
		assert.Equal(t, expectedPlayers[0].ID, response[0].ID)
		assert.Equal(t, expectedPlayers[1].ID, response[1].ID)
	})

	t.Run("invalid position", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/players?position=INVALID", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid team_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/players?team_id=invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetPlayersByTeam(t *testing.T) {
	router, mockService := setupPlayerHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		expectedPlayers := []*models.Player{
			{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Position:  "MID",
				TeamID:    1,
			},
			{
				ID:        2,
				FirstName: "Jane",
				LastName:  "Smith",
				Position:  "FWD",
				TeamID:    1,
			},
		}

		mockService.On("GetPlayersByTeam", 1).Return(expectedPlayers, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/teams/1/players", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.Player
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response, 2)
		assert.Equal(t, expectedPlayers[0].ID, response[0].ID)
		assert.Equal(t, expectedPlayers[1].ID, response[1].ID)
	})

	t.Run("invalid team id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/teams/invalid/players", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("team not found", func(t *testing.T) {
		mockService.On("GetPlayersByTeam", 999).Return([]*models.Player{}, sql.ErrNoRows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/teams/999/players", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetPlayerStats(t *testing.T) {
	router, mockService := setupPlayerHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		expectedStats := &models.PlayerStats{
			PlayerID: 1,
			Goals:    10,
			Assists:  5,
		}

		mockService.On("GetPlayerStats", 1).Return(expectedStats, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/players/1/stats", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.PlayerStats
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedStats.PlayerID, response.PlayerID)
		assert.Equal(t, expectedStats.Goals, response.Goals)
		assert.Equal(t, expectedStats.Assists, response.Assists)
	})

	t.Run("invalid player id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/players/invalid/stats", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("stats not found", func(t *testing.T) {
		mockService.On("GetPlayerStats", 999).Return(nil, sql.ErrNoRows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/players/999/stats", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
