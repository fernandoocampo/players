package unittests

import "github.com/stretchr/testify/mock"

type MockHasher struct {
	mock.Mock
}

func NewHasherMock() *MockHasher {
	return &MockHasher{}
}

func (m *MockHasher) Hash(password string) ([]byte, error) {
	args := m.Called(password)

	return args.Get(0).([]byte), args.Error(1)
}
