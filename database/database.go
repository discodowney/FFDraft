package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

// InitDB initializes the database connection
func InitDB(dbURL string) (*sqlx.DB, error) {
	var err error
	db, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// GetDB returns the database connectiona
func GetDB() *sqlx.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// GetItems retrieves all items from the database
func GetItems() ([]struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}, error) {
	var items []struct {
		ID   int    `db:"id" json:"id"`
		Name string `db:"name" json:"name"`
	}

	err := db.Select(&items, "SELECT id, name FROM items")
	if err != nil {
		return nil, err
	}

	return items, nil
}

// CreateItem creates a new item in the database
func CreateItem(name string) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO items (name) VALUES ($1) RETURNING id", name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
