package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// API Keys
	APIFootballBaseURL  string
	APIFootballAPIKey   string
	APIFootballLeagueID string
	APIFootballSeason   string

	// Database
	DatabaseURL string

	// Server
	ServerPort string
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Try to load .env file if it exists
	envPath := ".env"
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			return nil, err
		}
	}

	// Create config with default values
	config := &Config{
		// API Keys
		APIFootballBaseURL:  getEnvOrDefault("API_FOOTBALL_BASE_URL", "https://v3.football.api-sports.io"),
		APIFootballAPIKey:   getEnvOrDefault("API_FOOTBALL_API_KEY", ""),
		APIFootballLeagueID: getEnvOrDefault("API_FOOTBALL_LEAGUE_ID", "39"), // Default to Premier League
		APIFootballSeason:   getEnvOrDefault("API_FOOTBALL_SEASON", "2023"),  // Default to current season

		// Database
		DatabaseURL: getEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/go_app_db?sslmode=disable"),

		// Server
		ServerPort: getEnvOrDefault("SERVER_PORT", "8080"),
	}

	return config, nil
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Validate checks if required configuration is present
func (c *Config) Validate() error {
	if c.APIFootballAPIKey == "" {
		return ErrMissingRequiredConfig("API_FOOTBALL_API_KEY")
	}

	// Add validation for other required fields as needed

	return nil
}

// ErrMissingRequiredConfig is returned when a required configuration is missing
type ErrMissingRequiredConfig string

func (e ErrMissingRequiredConfig) Error() string {
	return "missing required configuration: " + string(e)
}
