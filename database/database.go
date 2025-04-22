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
