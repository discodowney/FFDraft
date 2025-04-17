package user

import (
	"bytes"
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

func setupTest(t *testing.T) (*gin.Engine, *mocks.MockUserService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(mocks.MockUserService)
	handler := &UserHandler{
		userService: mockService,
	}

	// Setup routes
	router.GET("/api/users/:id", handler.GetUser)
	router.GET("/api/users", handler.ListUsers)
	router.POST("/api/users", handler.CreateUser)
	router.PUT("/api/users/:id", handler.UpdateUser)
	router.DELETE("/api/users/:id", handler.DeleteUser)

	return router, mockService
}

func TestCreateUser(t *testing.T) {
	router, mockService := setupTest(t)

	// Test data
	user := &models.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Password:  "password123",
	}

	// Setup mock expectations
	mockService.On("ValidateUser", mock.Anything).Return(nil)
	mockService.On("CreateUser", mock.Anything).Return(user, nil)

	// Set request body with explicit fields
	requestBody := map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"password":   user.Password,
	}
	jsonData, _ := json.Marshal(requestBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusCreated, w.Code)

	var responseUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, user.FirstName, responseUser.FirstName)
	assert.Equal(t, user.LastName, responseUser.LastName)
	assert.Equal(t, user.Email, responseUser.Email)

	// Verify mock was called as expected
	mockService.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	router, mockService := setupTest(t)

	// Test data
	user := &models.User{
		ID:        1,
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
	}

	// Setup mock expectations
	mockService.On("GetUser", 1).Return(user, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users/1", nil)

	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var responseUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, responseUser.ID)
	assert.Equal(t, user.FirstName, responseUser.FirstName)
	assert.Equal(t, user.LastName, responseUser.LastName)
	assert.Equal(t, user.Email, responseUser.Email)

	// Verify mock was called as expected
	mockService.AssertExpectations(t)
}

func TestListUsers(t *testing.T) {
	router, mockService := setupTest(t)

	// Test data
	users := []*models.User{
		{
			ID:        1,
			FirstName: "Test1",
			LastName:  "User",
			Email:     "test1@example.com",
		},
		{
			ID:        2,
			FirstName: "Test2",
			LastName:  "User",
			Email:     "test2@example.com",
		},
	}

	// Setup mock expectations
	mockService.On("ListUsers").Return(users, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users", nil)

	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var responseUsers []models.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	assert.NoError(t, err)
	assert.Len(t, responseUsers, 2)
	assert.Equal(t, users[0].ID, responseUsers[0].ID)
	assert.Equal(t, users[1].ID, responseUsers[1].ID)

	// Verify mock was called as expected
	mockService.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	router, mockService := setupTest(t)

	// Test data
	user := &models.User{
		ID:        1,
		FirstName: "Updated",
		LastName:  "User",
		Email:     "updated@example.com",
	}

	// Setup mock expectations
	mockService.On("ValidateUser", user).Return(nil)
	mockService.On("UpdateUser", user).Return(user, nil)

	// Set request body
	jsonData, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var responseUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, responseUser.ID)
	assert.Equal(t, user.FirstName, responseUser.FirstName)
	assert.Equal(t, user.LastName, responseUser.LastName)
	assert.Equal(t, user.Email, responseUser.Email)

	// Verify mock was called as expected
	mockService.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	router, mockService := setupTest(t)

	// Setup mock expectations
	mockService.On("DeleteUser", 1).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/users/1", nil)

	router.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify mock was called as expected
	mockService.AssertExpectations(t)
}
