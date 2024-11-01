package unittests

import (
	"context"

	"github.com/fernandoocampo/players/internal/players"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func NewStorageMock() *MockStorage {
	return &MockStorage{}
}

func (m *MockStorage) Save(ctx context.Context, player players.Player) error {
	args := m.Called(ctx, player)

	return args.Error(0)
}

func (m *MockStorage) GetPlayersWithEmailOrNickName(ctx context.Context, filter players.PlayerFilter) (*players.PlayerExistResult, error) {
	args := m.Called(ctx, filter)

	return args.Get(0).(*players.PlayerExistResult), args.Error(1)
}

func (m *MockStorage) Delete(ctx context.Context, playerID players.PlayerID) error {
	args := m.Called(ctx, playerID)

	return args.Error(0)
}

func (m *MockStorage) Search(ctx context.Context, searchCriteria players.SearchCriteria) (*players.SearchResult, error) {
	args := m.Called(ctx, searchCriteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*players.SearchResult), args.Error(1)
}

func (m *MockStorage) Update(ctx context.Context, player players.Player) error {
	args := m.Called(ctx, player)

	return args.Error(0)
}

func (m *MockStorage) GetByID(ctx context.Context, id players.PlayerID) (*players.Player, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(*players.Player), args.Error(1)
}
