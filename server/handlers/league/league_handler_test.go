package league

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-app/models"
	"go-app/server/handlers/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupLeagueHandlerTest(t *testing.T) (*gin.Engine, *mocks.MockLeagueService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockLeagueService := new(mocks.MockLeagueService)
	handler := &LeagueHandler{
		leagueService: mockLeagueService,
	}

	// Setup routes
	router.GET("/leagues/:id", handler.GetLeague)
	router.GET("/leagues", handler.ListLeagues)
	router.POST("/leagues", handler.CreateLeague)
	router.PUT("/leagues/:id", handler.UpdateLeague)
	router.DELETE("/leagues/:id", handler.DeleteLeague)
	router.GET("/leagues/code/:code", handler.GetLeagueByCode)

	return router, mockLeagueService
}

func clearMockExpectations(mockService *mocks.MockLeagueService) {
	mockService.ExpectedCalls = nil
}

func TestGetLeague(t *testing.T) {
	router, mockLeagueService := setupLeagueHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		expectedLeague := &models.League{
			ID:   1,
			Name: "Test League",
			Code: "TEST123",
		}

		mockLeagueService.On("GetLeague", 1).Return(expectedLeague, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/leagues/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.League
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedLeague.ID, response.ID)
		assert.Equal(t, expectedLeague.Name, response.Name)
		assert.Equal(t, expectedLeague.Code, response.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/leagues/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		mockLeagueService.On("GetLeague", 999).Return(nil, sql.ErrNoRows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/leagues/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestListLeagues(t *testing.T) {
	router, mockLeagueService := setupLeagueHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		expectedLeagues := []*models.League{
			{
				ID:   1,
				Name: "League A",
				Code: "LGA",
			},
			{
				ID:   2,
				Name: "League B",
				Code: "LGB",
			},
		}

		mockLeagueService.On("ListLeagues").Return(expectedLeagues, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/leagues", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.League
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response, 2)
		assert.Equal(t, expectedLeagues[0].ID, response[0].ID)
		assert.Equal(t, expectedLeagues[1].ID, response[1].ID)
	})

	t.Run("error", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		mockLeagueService.On("ListLeagues").Return(nil, sql.ErrConnDone)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/leagues", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCreateLeague(t *testing.T) {
	router, mockLeagueService := setupLeagueHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		league := &models.League{
			Name: "New League",
			Code: "NEW123",
		}
		createdLeague := &models.League{
			ID:   1,
			Name: "New League",
			Code: "NEW123",
		}

		mockLeagueService.On("CreateLeague", league).Return(createdLeague, nil)

		body, _ := json.Marshal(league)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/leagues", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.League
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, createdLeague.ID, response.ID)
		assert.Equal(t, createdLeague.Name, response.Name)
		assert.Equal(t, createdLeague.Code, response.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/leagues", nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdateLeague(t *testing.T) {
	router, mockLeagueService := setupLeagueHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		league := &models.League{
			ID:   1,
			Name: "Updated League",
			Code: "UPD123",
		}

		mockLeagueService.On("UpdateLeague", league).Return(league, nil)

		body, _ := json.Marshal(league)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/leagues/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.League
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, league.ID, response.ID)
		assert.Equal(t, league.Name, response.Name)
		assert.Equal(t, league.Code, response.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/leagues/invalid", nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/leagues/1", nil)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteLeague(t *testing.T) {
	router, mockLeagueService := setupLeagueHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		mockLeagueService.On("DeleteLeague", 1).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/leagues/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/leagues/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		mockLeagueService.On("DeleteLeague", 999).Return(sql.ErrNoRows)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/leagues/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetLeagueByCode(t *testing.T) {
	router, mockLeagueService := setupLeagueHandlerTest(t)

	t.Run("success", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		expectedLeague := &models.League{
			ID:   1,
			Name: "Test League",
			Code: "TEST123",
		}

		mockLeagueService.On("GetLeagueByCode", "TEST123").Return(expectedLeague, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/leagues/code/TEST123", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.League
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expectedLeague.ID, response.ID)
		assert.Equal(t, expectedLeague.Name, response.Name)
		assert.Equal(t, expectedLeague.Code, response.Code)
	})

	t.Run("not found", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		mockLeagueService.On("GetLeagueByCode", "NONEXIST").Return(nil, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/leagues/code/NONEXIST", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		defer clearMockExpectations(mockLeagueService)
		mockLeagueService.On("GetLeagueByCode", "ERROR").Return(nil, sql.ErrConnDone)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/leagues/code/ERROR", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
