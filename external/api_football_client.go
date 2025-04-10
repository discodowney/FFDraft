package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go-app/models"
)

// APIFootballClient handles communication with the API-Football service
type APIFootballClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewAPIFootballClient creates a new APIFootballClient instance
func NewAPIFootballClient(apiKey string) *APIFootballClient {
	return &APIFootballClient{
		BaseURL: "https://api-football-v1.p.rapidapi.com/v3",
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// FetchTeams retrieves teams from the API-Football service
func (c *APIFootballClient) FetchTeams() ([]*models.Team, error) {
	// Construct the API URL
	url := fmt.Sprintf("%s/teams", c.BaseURL)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers
	req.Header.Add("x-rapidapi-host", "api-football-v1.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", c.APIKey)

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var response struct {
		Response []struct {
			Team struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"team"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Convert to our model
	teams := make([]*models.Team, 0, len(response.Response))
	for _, item := range response.Response {
		team := &models.Team{
			ID:         0, // Will be set by database
			Name:       item.Team.Name,
			ExternalId: item.Team.ID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		teams = append(teams, team)
	}

	return teams, nil
}

// FetchTeamByExternalID retrieves a team from the API-Football service by its external ID
func (c *APIFootballClient) FetchTeamByExternalID(externalID int) (*models.Team, error) {
	// Construct the API URL
	url := fmt.Sprintf("%s/teams?id=%d", c.BaseURL, externalID)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers
	req.Header.Add("x-rapidapi-host", "api-football-v1.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", c.APIKey)

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var response struct {
		Response []struct {
			Team struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"team"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Check if we got any results
	if len(response.Response) == 0 {
		return nil, fmt.Errorf("team with external ID %d not found", externalID)
	}

	// Convert to our model
	team := &models.Team{
		ID:         0, // Will be set by database
		Name:       response.Response[0].Team.Name,
		ExternalId: response.Response[0].Team.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return team, nil
}
