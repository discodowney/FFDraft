package external

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go-app/models"
)

// TeamAPIClient handles communication with the external team API
type TeamAPIClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewTeamAPIClient creates a new TeamAPIClient instance
func NewTeamAPIClient(baseURL, apiKey string) *TeamAPIClient {
	return &TeamAPIClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// FetchTeams retrieves teams from the external API
func (c *TeamAPIClient) FetchTeams() ([]*models.Team, error) {
	// Construct the API URL
	url := fmt.Sprintf("%s/teams", c.BaseURL)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add headers
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Add("Content-Type", "application/json")

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
	var teams []*models.Team
	if err := json.NewDecoder(resp.Body).Decode(&teams); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return teams, nil
}
