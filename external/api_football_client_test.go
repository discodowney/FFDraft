package external

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewAPIFootballClient(t *testing.T) {
	client := NewAPIFootballClient("http://test.com", "test-key", "123", "2023")

	if client.BaseURL != "http://test.com" {
		t.Errorf("Expected BaseURL to be 'http://test.com', got '%s'", client.BaseURL)
	}
	if client.APIKey != "test-key" {
		t.Errorf("Expected APIKey to be 'test-key', got '%s'", client.APIKey)
	}
	if client.LeagueID != "123" {
		t.Errorf("Expected LeagueID to be '123', got '%s'", client.LeagueID)
	}
	if client.Season != "2023" {
		t.Errorf("Expected Season to be '2023', got '%s'", client.Season)
	}
	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized")
	}
}

func TestFetchTeams(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("x-rapidapi-key") != "test-key" {
			t.Errorf("Expected x-rapidapi-key to be 'test-key', got '%s'", r.Header.Get("x-rapidapi-key"))
		}
		if r.Header.Get("x-rapidapi-host") != "api-football-v1.p.rapidapi.com" {
			t.Errorf("Expected x-rapidapi-host to be 'api-football-v1.p.rapidapi.com', got '%s'", r.Header.Get("x-rapidapi-host"))
		}

		// Verify request URL
		expectedURL := "/teams?league=123&season=2023"
		if r.URL.String() != expectedURL {
			t.Errorf("Expected URL to be '%s', got '%s'", expectedURL, r.URL.String())
		}

		// Return mock response
		response := map[string]interface{}{
			"response": []map[string]interface{}{
				{
					"team": map[string]interface{}{
						"id":   1,
						"name": "Test Team 1",
					},
				},
				{
					"team": map[string]interface{}{
						"id":   2,
						"name": "Test Team 2",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewAPIFootballClient(server.URL, "test-key", "123", "2023")

	// Call FetchTeams
	teams, err := client.FetchTeams()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify results
	if len(teams) != 2 {
		t.Errorf("Expected 2 teams, got %d", len(teams))
	}

	expectedTeams := []struct {
		name       string
		externalID int
	}{
		{"Test Team 1", 1},
		{"Test Team 2", 2},
	}

	for i, expected := range expectedTeams {
		if teams[i].Name != expected.name {
			t.Errorf("Expected team name to be '%s', got '%s'", expected.name, teams[i].Name)
		}
		if teams[i].ExternalId != expected.externalID {
			t.Errorf("Expected external ID to be %d, got %d", expected.externalID, teams[i].ExternalId)
		}
	}
}

func TestFetchTeamByExternalID(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("x-rapidapi-key") != "test-key" {
			t.Errorf("Expected x-rapidapi-key to be 'test-key', got '%s'", r.Header.Get("x-rapidapi-key"))
		}
		if r.Header.Get("x-rapidapi-host") != "api-football-v1.p.rapidapi.com" {
			t.Errorf("Expected x-rapidapi-host to be 'api-football-v1.p.rapidapi.com', got '%s'", r.Header.Get("x-rapidapi-host"))
		}

		// Verify request URL
		expectedURL := "/teams?id=123"
		if r.URL.String() != expectedURL {
			t.Errorf("Expected URL to be '%s', got '%s'", expectedURL, r.URL.String())
		}

		// Return mock response
		response := map[string]interface{}{
			"response": []map[string]interface{}{
				{
					"team": map[string]interface{}{
						"id":   123,
						"name": "Test Team",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewAPIFootballClient(server.URL, "test-key", "123", "2023")

	// Call FetchTeamByExternalID
	team, err := client.FetchTeamByExternalID(123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify results
	if team.Name != "Test Team" {
		t.Errorf("Expected team name to be 'Test Team', got '%s'", team.Name)
	}
	if team.ExternalId != 123 {
		t.Errorf("Expected external ID to be 123, got %d", team.ExternalId)
	}
}

func TestFetchTeamsErrorHandling(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewAPIFootballClient(server.URL, "test-key", "123", "2023")

	_, err := client.FetchTeams()
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestFetchTeamByExternalIDErrorHandling(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewAPIFootballClient(server.URL, "test-key", "123", "2023")

	_, err := client.FetchTeamByExternalID(123)
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestFetchTeamByExternalIDNotFound(t *testing.T) {
	// Create a test server that returns no teams
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"response": []map[string]interface{}{},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewAPIFootballClient(server.URL, "test-key", "123", "2023")

	_, err := client.FetchTeamByExternalID(123)
	if err == nil {
		t.Error("Expected an error for not found team, got nil")
	}
}
