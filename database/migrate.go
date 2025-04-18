package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Migration represents a database migration
type Migration struct {
	Version string
	SQL     string
}

// RunMigrations executes all pending migrations
func RunMigrations(db *sqlx.DB) error {
	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %w", err)
	}

	// Get list of applied migrations
	var appliedMigrations []string
	err = db.Select(&appliedMigrations, "SELECT version FROM migrations ORDER BY version")
	if err != nil {
		return fmt.Errorf("error getting applied migrations: %w", err)
	}

	// Read migration files
	migrationsDir := "database/migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %w", err)
	}

	// Sort files by name (which should be numbered)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Apply pending migrations
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		version := strings.TrimSuffix(file.Name(), ".sql")
		if contains(appliedMigrations, version) {
			continue
		}

		// Read migration file
		content, err := ioutil.ReadFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %w", file.Name(), err)
		}

		// Start transaction
		tx, err := db.Beginx()
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}

		// Execute migration
		_, err = tx.Exec(string(content))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing migration %s: %w", version, err)
		}

		// Record migration
		_, err = tx.Exec("INSERT INTO migrations (version) VALUES ($1)", version)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error recording migration %s: %w", version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing migration %s: %w", version, err)
		}

		log.Printf("Applied migration: %s", version)
	}

	return nil
}

// Helper function to check if a string is in a slice
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
