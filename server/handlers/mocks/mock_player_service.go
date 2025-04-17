package mocks

import (
	"go-app/models"
	"go-app/services/player"

	"github.com/stretchr/testify/mock"
)

type MockPlayerService struct {
	mock.Mock
}

func (m *MockPlayerService) GetPlayer(id int) (*models.Player, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Player), args.Error(1)
}

func (m *MockPlayerService) ListPlayers(filter *player.PlayerFilter) ([]*models.Player, error) {
	args := m.Called(filter)
	return args.Get(0).([]*models.Player), args.Error(1)
}

func (m *MockPlayerService) GetPlayersByTeam(teamID int) ([]*models.Player, error) {
	args := m.Called(teamID)
	return args.Get(0).([]*models.Player), args.Error(1)
}

func (m *MockPlayerService) GetPlayerStats(playerID int) (*models.PlayerStats, error) {
	args := m.Called(playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PlayerStats), args.Error(1)
}

func (m *MockPlayerService) ValidatePosition(position models.Position) error {
	args := m.Called(position)
	return args.Error(0)
}

func (m *MockPlayerService) DeletePlayer(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPlayerService) UpdatePlayer(player *models.Player) (*models.Player, error) {
	args := m.Called(player)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Player), args.Error(1)
}

func (m *MockPlayerService) CreatePlayer(player *models.Player) (*models.Player, error) {
	args := m.Called(player)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Player), args.Error(1)
}

func (m *MockPlayerService) ValidatePlayer(player *models.Player) error {
	args := m.Called(player)
	return args.Error(0)
}
