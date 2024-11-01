package unittests

import (
	"log/slog"
	"net/mail"
	"testing"

	"github.com/fernandoocampo/players/internal/players"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func NewPlayerID() players.PlayerID {
	return players.PlayerID(uuid.New())
}

func NewPlayerFixture(t *testing.T) players.NewPlayer {
	t.Helper()

	return players.NewPlayer{
		FirstName: "Fernando",
		LastName:  "Ocampo",
		Nickname:  "focampo",
		Email:     *NewEmailAddress(t, "focampo@anyemail.com"),
		Password:  "th1s1s@dummypwd",
		Country:   "Spain",
	}
}

func NewEmailAddress(t *testing.T, address string) *mail.Address {
	t.Helper()

	emailAddress, err := mail.ParseAddress(address)
	if err != nil {
		t.Errorf("email address is not valid: %s", address)
		t.FailNow()
	}

	return emailAddress
}

func NewPlayerService(storageMock players.Storage, hasherMock players.Hasher, notifierMock players.Notifier) (*players.Service, *slog.Logger) {
	serviceSetup := players.ServiceSetup{
		Storage:  storageMock,
		Hasher:   hasherMock,
		Notifier: notifierMock,
		Logger:   NewLogger(),
	}

	service := players.NewService(&serviceSetup)

	return service, serviceSetup.Logger
}

func NewPlayerServiceWithStorage(storageMock players.Storage) (*players.Service, *slog.Logger) {
	return NewPlayerService(storageMock, NewHasherMock(), NewNotifierMock())
}

func NewPlayerServiceWithStorageAndNotifier(storageMock players.Storage, notifierMock players.Notifier) (*players.Service, *slog.Logger) {
	return NewPlayerService(storageMock, nil, notifierMock)
}

func PlayerIDFixture(t *testing.T, givenPlayerID string) *players.PlayerID {
	t.Helper()

	playerID, err := uuid.Parse(givenPlayerID)
	require.NoError(t, err)

	newPlayerID := players.PlayerID(playerID)

	return &newPlayerID
}

func SearchPlayersResultFixture(t *testing.T) []players.PlayerItem {
	t.Helper()

	return []players.PlayerItem{
		{
			ID:        *PlayerIDFixture(t, "b10d6af5-22f3-4db2-ade6-94cfcc819f91"),
			FirstName: "Jim",
			LastName:  "Raynor",
			Nickname:  "jraynor",
			Country:   "UK",
		},
		{
			ID:        *PlayerIDFixture(t, "8b1cb38a-d21c-4e36-b7f0-814f96180b5e"),
			FirstName: "Sarah",
			LastName:  "Kerrigan",
			Nickname:  "skerrigan",
			Country:   "UK",
		},
		{
			ID:        *PlayerIDFixture(t, "f7ee12ee-191d-4f7e-819e-dbb7b115bed2"),
			FirstName: "Graven",
			LastName:  "Hill",
			Nickname:  "ghill",
			Country:   "UK",
		},
	}
}
