package mocks

import (
	"go-app/models"

	"github.com/stretchr/testify/mock"
)

type MockTeamService struct {
	mock.Mock
}

func (m *MockTeamService) GetTeam(id int64) (*models.Team, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamService) GetTeamByExternalID(externalID int64) (*models.Team, error) {
	args := m.Called(externalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamService) ListTeams() ([]*models.Team, error) {
	args := m.Called()
	return args.Get(0).([]*models.Team), args.Error(1)
}

func (m *MockTeamService) CreateTeam(team *models.Team) (*models.Team, error) {
	args := m.Called(team)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamService) UpdateTeam(team *models.Team) (*models.Team, error) {
	args := m.Called(team)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamService) DeleteTeam(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTeamService) ValidateTeam(team *models.Team) error {
	args := m.Called(team)
	return args.Error(0)
}
