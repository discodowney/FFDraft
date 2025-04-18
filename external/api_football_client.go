package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go-app/models"
)

// APIFootballClientInterface defines the interface for API-Football client operations
type APIFootballClientInterface interface {
	FetchTeams() ([]*models.Team, error)
	FetchTeamByExternalID(externalID int) (*models.Team, error)
}

// APIFootballClient handles communication with the API-Football service
type APIFootballClient struct {
	BaseURL    string
	APIKey     string
	LeagueID   string
	Season     string
	HTTPClient *http.Client
}

// Ensure APIFootballClient implements APIFootballClientInterface
var _ APIFootballClientInterface = (*APIFootballClient)(nil)

// NewAPIFootballClient creates a new APIFootballClient instance
func NewAPIFootballClient(baseURL, apiKey, leagueID, season string) *APIFootballClient {
	return &APIFootballClient{
		BaseURL:  baseURL,
		APIKey:   apiKey,
		LeagueID: leagueID,
		Season:   season,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// FetchTeams retrieves teams from the API-Football service
func (c *APIFootballClient) FetchTeams() ([]*models.Team, error) {
	url := fmt.Sprintf("%s/teams?league=%s&season=%s", c.BaseURL, c.LeagueID, c.Season)

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
