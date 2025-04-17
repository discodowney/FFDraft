package user

import (
	"fmt"
	"regexp"
	"time"

	"go-app/models"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user-related operations
type UserService interface {
	GetUser(id int) (*models.User, error)
	ListUsers() ([]*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id int) error
	ValidateUser(user *models.User) error
}

// UserService handles user-related operations
type userServiceImpl struct {
	db *sqlx.DB
}

// NewUserService creates a new UserService instance
func NewUserService(db *sqlx.DB) UserService {
	return &userServiceImpl{db: db}
}

// CreateUser creates a new user in the database
func (s *userServiceImpl) CreateUser(user *models.User) (*models.User, error) {
	// Validate user data
	if err := s.ValidateUser(user); err != nil {
		return nil, err
	}

	// Check if email already exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", user.Email).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Insert user into database
	var id int
	err = s.db.QueryRow(`
		INSERT INTO users (first_name, last_name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, user.FirstName, user.LastName, user.Email, string(hashedPassword), user.CreatedAt, user.UpdatedAt).Scan(&id)
	if err != nil {
		return nil, err
	}

	// Set the ID and return the user
	user.ID = id
	return user, nil
}

// GetUser retrieves a user by ID
func (s *userServiceImpl) GetUser(id int) (*models.User, error) {
	user := &models.User{}
	err := s.db.Get(user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates an existing user in the database
func (s *userServiceImpl) UpdateUser(user *models.User) (*models.User, error) {
	// Validate user data
	if err := s.ValidateUser(user); err != nil {
		return nil, err
	}

	// Check if user exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", user.ID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, fmt.Errorf("user with ID %d not found", user.ID)
	}

	// Check if email is being changed and if it already exists
	if user.Email != "" {
		var emailCount int
		err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1 AND id != $2", user.Email, user.ID).Scan(&emailCount)
		if err != nil {
			return nil, err
		}
		if emailCount > 0 {
			return nil, fmt.Errorf("email already exists")
		}
	}

	// Set updated timestamp
	user.UpdatedAt = time.Now()

	// Update user in database
	_, err = s.db.Exec(`
		UPDATE users 
		SET first_name = $1, last_name = $2, email = $3, updated_at = $4
		WHERE id = $5
	`, user.FirstName, user.LastName, user.Email, user.UpdatedAt, user.ID)
	if err != nil {
		return nil, err
	}

	// Return the updated user
	return user, nil
}

// DeleteUser deletes a user by ID
func (s *userServiceImpl) DeleteUser(id int) error {
	result, err := s.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

// ListUsers retrieves all users
func (s *userServiceImpl) ListUsers() ([]*models.User, error) {
	var users []*models.User
	err := s.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ValidateUser validates user data
func (s *userServiceImpl) ValidateUser(user *models.User) error {
	if user.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if user.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !isValidEmail(user.Email) {
		return fmt.Errorf("invalid email format")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(user.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}

// isValidEmail checks if the email is valid
func isValidEmail(email string) bool {
	// This is a simple regex for email validation
	// It matches most common email formats but may not catch all edge cases
	emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	return emailRegex.MatchString(email)
}
