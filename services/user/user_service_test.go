package user

import (
	"fmt"
	"testing"

	"go-app/database"
	"go-app/models"

	"github.com/stretchr/testify/assert"
)

var testDB *database.TestDB
var userService UserService

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

	userService = NewUserService(testDB.GetDB())

	// Run tests
	m.Run()
}

func TestUserService(t *testing.T) {
	// Test CreateUser
	t.Run("CreateUser", func(t *testing.T) {
		defer testDB.Clear()

		user := &models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@example.com",
			Password:  "password123",
		}

		// Test successful creation
		createdUser, err := userService.CreateUser(user)
		assert.NoError(t, err)
		assert.NotZero(t, createdUser.ID)
		assert.Equal(t, user.FirstName, createdUser.FirstName)
		assert.Equal(t, user.LastName, createdUser.LastName)
		assert.Equal(t, user.Email, createdUser.Email)
		assert.NotEmpty(t, createdUser.Password) // Password should be hashed

		// Test duplicate email
		_, err = userService.CreateUser(user)
		assert.Error(t, err)
	})

	// Test GetUser
	t.Run("GetUser", func(t *testing.T) {
		defer testDB.Clear()

		// Create a test user
		user := &models.User{
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane@example.com",
			Password:  "password123",
		}
		createdUser, err := userService.CreateUser(user)
		assert.NoError(t, err)

		// Test successful retrieval
		retrievedUser, err := userService.GetUser(createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdUser.ID, retrievedUser.ID)
		assert.Equal(t, user.FirstName, retrievedUser.FirstName)
		assert.Equal(t, user.LastName, retrievedUser.LastName)
		assert.Equal(t, user.Email, retrievedUser.Email)

		// Test non-existent user
		_, err = userService.GetUser(999)
		assert.Error(t, err)
	})

	// Test UpdateUser
	t.Run("UpdateUser", func(t *testing.T) {
		defer testDB.Clear()

		// Create a test user
		user := &models.User{
			FirstName: "Original",
			LastName:  "User",
			Email:     "original@example.com",
			Password:  "password123",
		}
		createdUser, err := userService.CreateUser(user)
		assert.NoError(t, err)

		// Update user
		createdUser.FirstName = "Updated"
		createdUser.LastName = "Name"
		createdUser.Email = "updated@example.com"
		updatedUser, err := userService.UpdateUser(createdUser)
		assert.NoError(t, err)
		assert.Equal(t, "Updated", updatedUser.FirstName)
		assert.Equal(t, "Name", updatedUser.LastName)
		assert.Equal(t, "updated@example.com", updatedUser.Email)

		// Verify update in database
		retrievedUser, err := userService.GetUser(createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated", retrievedUser.FirstName)
		assert.Equal(t, "Name", retrievedUser.LastName)
		assert.Equal(t, "updated@example.com", retrievedUser.Email)
	})

	// Test DeleteUser
	t.Run("DeleteUser", func(t *testing.T) {
		defer testDB.Clear()

		// Create a test user
		user := &models.User{
			FirstName: "ToDelete",
			LastName:  "User",
			Email:     "delete@example.com",
			Password:  "password123",
		}
		createdUser, err := userService.CreateUser(user)
		assert.NoError(t, err)

		// Delete user
		err = userService.DeleteUser(createdUser.ID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = userService.GetUser(createdUser.ID)
		assert.Error(t, err)
	})

	// Test ListUsers
	t.Run("ListUsers", func(t *testing.T) {
		defer testDB.Clear()

		// Create multiple users
		users := []*models.User{
			{
				FirstName: "User1",
				LastName:  "Test",
				Email:     "user1@example.com",
				Password:  "password123",
			},
			{
				FirstName: "User2",
				LastName:  "Test",
				Email:     "user2@example.com",
				Password:  "password123",
			},
		}

		for _, user := range users {
			_, err := userService.CreateUser(user)
			assert.NoError(t, err)
		}

		// List users
		retrievedUsers, err := userService.ListUsers()
		assert.NoError(t, err)
		assert.Len(t, retrievedUsers, 2)
	})

	// Test ValidateUser
	t.Run("ValidateUser", func(t *testing.T) {
		defer testDB.Clear()

		// Test valid user
		validUser := &models.User{
			FirstName: "Valid",
			LastName:  "User",
			Email:     "valid@example.com",
			Password:  "password123",
		}
		err := userService.ValidateUser(validUser)
		assert.NoError(t, err)

		// Test missing first name
		invalidUser := &models.User{
			LastName: "User",
			Email:    "test@example.com",
			Password: "password123",
		}
		err = userService.ValidateUser(invalidUser)
		assert.Error(t, err)

		// Test missing last name
		invalidUser = &models.User{
			FirstName: "User",
			Email:     "test@example.com",
			Password:  "password123",
		}
		err = userService.ValidateUser(invalidUser)
		assert.Error(t, err)

		// Test missing email
		invalidUser = &models.User{
			FirstName: "User",
			LastName:  "Test",
			Password:  "password123",
		}
		err = userService.ValidateUser(invalidUser)
		assert.Error(t, err)

		// Test invalid email format
		invalidUser = &models.User{
			FirstName: "User",
			LastName:  "Test",
			Email:     "invalid-email",
			Password:  "password123",
		}
		err = userService.ValidateUser(invalidUser)
		assert.Error(t, err)

		// Test missing password
		invalidUser = &models.User{
			FirstName: "User",
			LastName:  "Test",
			Email:     "test@example.com",
		}
		err = userService.ValidateUser(invalidUser)
		assert.Error(t, err)
	})
}
