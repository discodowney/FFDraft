package mocks

import (
	"go-app/external"
	"go-app/models"
)

// MockAPIFootballClient is a mock implementation of the APIFootballClient
type MockAPIFootballClient struct {
	teams []*models.Team
	err   error
}

// Ensure MockAPIFootballClient implements APIFootballClientInterface
var _ external.APIFootballClientInterface = (*MockAPIFootballClient)(nil)

// NewMockAPIFootballClient creates a new mock API Football client
func NewMockAPIFootballClient(teams []*models.Team, err error) *MockAPIFootballClient {
	return &MockAPIFootballClient{
		teams: teams,
		err:   err,
	}
}

// FetchTeams returns the mock teams and error
func (m *MockAPIFootballClient) FetchTeams() ([]*models.Team, error) {
	return m.teams, m.err
}

// FetchTeamByExternalID returns a team by its external ID
func (m *MockAPIFootballClient) FetchTeamByExternalID(externalID int) (*models.Team, error) {
	for _, team := range m.teams {
		if team.ExternalId == externalID {
			return team, nil
		}
	}
	return nil, nil
}
