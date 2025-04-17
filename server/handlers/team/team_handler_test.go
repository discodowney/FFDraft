package team

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

func setupTeamHandlerTest(t *testing.T) (*gin.Engine, *mocks.MockTeamService, *mocks.MockPlayerService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockTeamService := new(mocks.MockTeamService)
	mockPlayerService := new(mocks.MockPlayerService)
	handler := &TeamHandler{
		teamService:   mockTeamService,
		playerService: mockPlayerService,
	}

	// Setup routes
	router.GET("/teams/:id", handler.GetTeam)
	router.GET("/teams", handler.ListTeams)
	router.GET("/teams/:id/players", handler.GetTeamPlayers)

	return router, mockTeamService, mockPlayerService
}

func TestGetTeam(t *testing.T) {
	router, mockTeamService, _ := setupTeamHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		expectedTeam := &models.Team{
			ID:   1,
			Name: "Test Team",
		}

		mockTeamService.On("GetTeam", int64(1)).Return(expectedTeam, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/teams/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Team
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedTeam.ID, response.ID)
		assert.Equal(t, expectedTeam.Name, response.Name)
	})

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/teams/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockTeamService.On("GetTeam", int64(999)).Return(nil, sql.ErrNoRows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/teams/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestListTeams(t *testing.T) {
	router, mockTeamService, _ := setupTeamHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		expectedTeams := []models.Team{
			{
				ID:   1,
				Name: "Team A",
			},
			{
				ID:   2,
				Name: "Team B",
			},
		}

		mockTeamService.On("ListTeams", mock.Anything).Return(expectedTeams, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/teams", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []models.Team
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response, 2)
		assert.Equal(t, expectedTeams[0].ID, response[0].ID)
		assert.Equal(t, expectedTeams[1].ID, response[1].ID)
	})

	t.Run("with filters", func(t *testing.T) {
		expectedTeams := []models.Team{
			{
				ID:   1,
				Name: "Team A",
			},
		}

		mockTeamService.On("ListTeams", mock.Anything).Return(expectedTeams, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/teams?name=Team A&external_id=123", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []models.Team
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response, 1)
		assert.Equal(t, expectedTeams[0].ID, response[0].ID)
	})

	t.Run("invalid external_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/teams?external_id=invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetTeamPlayers(t *testing.T) {
	router, _, mockPlayerService := setupTeamHandlerTest(t)

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

		mockPlayerService.On("GetPlayersByTeam", 1).Return(expectedPlayers, nil)

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
		mockPlayerService.On("GetPlayersByTeam", 999).Return([]*models.Player{}, sql.ErrNoRows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/teams/999/players", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
