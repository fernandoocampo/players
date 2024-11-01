package unittests

import (
	"github.com/fernandoocampo/players/internal/players"
	"github.com/stretchr/testify/mock"
)

type MockNotifier struct {
	mock.Mock
}

func NewNotifierMock() *MockNotifier {
	return &MockNotifier{}
}

func (m *MockNotifier) Notify(event players.NewEvent) {
	_ = m.Called(event)
}
