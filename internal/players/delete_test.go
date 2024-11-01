package players_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fernandoocampo/players/internal/appkit/unittests"
	"github.com/fernandoocampo/players/internal/players"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeletePlayer(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	givenPlayerID := unittests.PlayerIDFixture(t, "b10d6af5-22f3-4db2-ade6-94cfcc819f91")

	existingPlayer := players.Player{
		ID:          givenPlayerID,
		FirstName:   "Fernando",
		LastName:    "Ocampo",
		Nickname:    "focampo",
		Email:       *unittests.NewEmailAddress(t, "focampo@anyemail.com"),
		Password:    []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6"),
		Country:     "Spain",
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetByID", ctx, *givenPlayerID).Return(&existingPlayer, nil)
	storageMock.On("Delete", ctx, *givenPlayerID).Return(nil)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent")).Return(nil)

	service, _ := unittests.NewPlayerServiceWithStorageAndNotifier(storageMock, notifierMock)

	// When
	err := service.Delete(ctx, *givenPlayerID)

	// Then
	assert.NoError(t, err)
}

func TestDeletePlayerButErrorWhileDeleting(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	givenPlayerID := unittests.PlayerIDFixture(t, "b10d6af5-22f3-4db2-ade6-94cfcc819f91")

	existingPlayer := players.Player{
		ID:          givenPlayerID,
		FirstName:   "Fernando",
		LastName:    "Ocampo",
		Nickname:    "focampo",
		Email:       *unittests.NewEmailAddress(t, "focampo@anyemail.com"),
		Password:    []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6"),
		Country:     "Spain",
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}

	deleteError := errors.New("unexpected delete error")

	want := "unable to delete player: unexpected delete error"

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetByID", ctx, *givenPlayerID).Return(&existingPlayer, nil)
	storageMock.On("Delete", ctx, *givenPlayerID).Return(deleteError)

	service, _ := unittests.NewPlayerServiceWithStorage(storageMock)

	// When
	err := service.Delete(ctx, *givenPlayerID)

	// Then
	assert.Error(t, err)
	assert.Equal(t, want, err.Error())
}

func TestDeletePlayerWithEndpointSuccessfully(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	givenPlayerID := unittests.PlayerIDFixture(t, "b10d6af5-22f3-4db2-ade6-94cfcc819f91")

	existingPlayer := players.Player{
		ID:          givenPlayerID,
		FirstName:   "Fernando",
		LastName:    "Ocampo",
		Nickname:    "focampo",
		Email:       *unittests.NewEmailAddress(t, "focampo@anyemail.com"),
		Password:    []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6"),
		Country:     "Spain",
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}

	want := players.DeletePlayerResult{}

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetByID", ctx, *givenPlayerID).Return(&existingPlayer, nil)
	storageMock.On("Delete", ctx, *givenPlayerID).Return(nil)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent")).Return(nil)

	service, logger := unittests.NewPlayerServiceWithStorageAndNotifier(storageMock, notifierMock)
	deletePlayerEndpoint := players.MakeDeletePlayerEndpoint(service, logger)

	// When
	got, err := deletePlayerEndpoint.Do(ctx, *givenPlayerID)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, want, got)
}

func TestDeletePlayerWithEndpointButError(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	givenPlayerID := unittests.PlayerIDFixture(t, "b10d6af5-22f3-4db2-ade6-94cfcc819f91")

	existingPlayer := players.Player{
		ID:          givenPlayerID,
		FirstName:   "Fernando",
		LastName:    "Ocampo",
		Nickname:    "focampo",
		Email:       *unittests.NewEmailAddress(t, "focampo@anyemail.com"),
		Password:    []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6"),
		Country:     "Spain",
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}

	deleteError := errors.New("unexpected delete error")

	want := players.DeletePlayerResult{
		Err: "unable to delete player: unexpected delete error",
	}

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetByID", ctx, *givenPlayerID).Return(&existingPlayer, nil)
	storageMock.On("Delete", ctx, *givenPlayerID).Return(deleteError)

	service, logger := unittests.NewPlayerServiceWithStorage(storageMock)
	deletePlayerEndpoint := players.MakeDeletePlayerEndpoint(service, logger)

	// When
	got, err := deletePlayerEndpoint.Do(ctx, *givenPlayerID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
