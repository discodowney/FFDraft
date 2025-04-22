package mocks

import (
	"go-app/models"
	"go-app/services/league"

	"github.com/stretchr/testify/mock"
)

type MockLeagueService struct {
	mock.Mock
}

func (m *MockLeagueService) CreateLeague(league *models.League) (*models.League, error) {
	args := m.Called(league)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.League), args.Error(1)
}

func (m *MockLeagueService) GetLeague(id int) (*models.League, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.League), args.Error(1)
}

func (m *MockLeagueService) UpdateLeague(league *models.League) (*models.League, error) {
	args := m.Called(league)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.League), args.Error(1)
}

func (m *MockLeagueService) DeleteLeague(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLeagueService) ListLeagues() ([]*models.League, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	if args.Get(0) == nil {
		return nil, nil
	}
	return args.Get(0).([]*models.League), nil
}

func (m *MockLeagueService) ValidateLeague(league *models.League) error {
	args := m.Called(league)
	return args.Error(0)
}

func (m *MockLeagueService) GetLeagueByCode(code string) (*models.League, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.League), args.Error(1)
}

var _ league.LeagueService = (*MockLeagueService)(nil)
