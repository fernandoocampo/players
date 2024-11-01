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

func TestUpdatePlayer(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	playerID := unittests.NewPlayerID()
	newPasswordValue := "newpwd"
	dateCreated := time.Now().UTC()
	updatePlayer := players.UpdatePlayer{
		ID:       playerID,
		Email:    unittests.NewEmailAddress(t, "focampo@anotheremail.com"),
		Password: &newPasswordValue,
	}

	existingPlayer := existingPlayerToUpdateFixture(t, &playerID, dateCreated)

	wantedPlayer := players.Player{
		ID:          &playerID,
		FirstName:   "Fernando",
		LastName:    "Ocampo",
		Nickname:    "focampo",
		Email:       *unittests.NewEmailAddress(t, "focampo@anotheremail.com"),
		Password:    []byte("$3a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6"),
		Country:     "Spain",
		DateCreated: dateCreated,
	}

	givenPlayerExistResult := players.PlayerExistResult{
		EmailExist:    false,
		NicknameExist: false,
	}

	dummyHash := []byte("$3a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")

	hasherMock := unittests.NewHasherMock()
	hasherMock.On("Hash", *updatePlayer.Password).Return(dummyHash, nil)

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetByID", ctx, playerID).Return(&existingPlayer, nil)
	storageMock.On("GetPlayersWithEmailOrNickName", ctx, mock.AnythingOfType("players.PlayerFilter")).Return(&givenPlayerExistResult, nil)
	storageMock.On("Update", ctx, mock.AnythingOfType("players.Player")).Return(nil)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent"))

	service, _ := unittests.NewPlayerService(storageMock, hasherMock, notifierMock)

	// When
	got, err := service.Update(ctx, updatePlayer)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, got)
	wantedPlayer.DateUpdated = got.DateUpdated // random values
	assert.Equal(t, &wantedPlayer, got)
}

func TestUpdatePlayerButErrorWhileUpdating(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	playerID := unittests.NewPlayerID()
	newPasswordValue := "newpwd"
	dateCreated := time.Now().UTC()
	updatePlayer := players.UpdatePlayer{
		ID:       playerID,
		Email:    unittests.NewEmailAddress(t, "focampo@anotheremail.com"),
		Password: &newPasswordValue,
	}

	existingPlayer := existingPlayerToUpdateFixture(t, &playerID, dateCreated)

	givenPlayerExistResult := players.PlayerExistResult{
		EmailExist:    false,
		NicknameExist: false,
	}

	saveError := errors.New("unexpected update error")

	dummyHash := []byte("$3a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")

	want := "unable to update player: unexpected update error"

	hasherMock := unittests.NewHasherMock()
	hasherMock.On("Hash", *updatePlayer.Password).Return(dummyHash, nil)

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetByID", ctx, playerID).Return(&existingPlayer, nil)
	storageMock.On("GetPlayersWithEmailOrNickName", ctx, mock.AnythingOfType("players.PlayerFilter")).Return(&givenPlayerExistResult, nil)
	storageMock.On("Update", ctx, mock.AnythingOfType("players.Player")).Return(saveError)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent"))

	service, _ := unittests.NewPlayerService(storageMock, hasherMock, notifierMock)

	// When
	player, err := service.Update(ctx, updatePlayer)

	// Then
	assert.Error(t, err)
	assert.Nil(t, player)
	assert.Equal(t, want, err.Error())
}

func TestUpdatePlayerWithEndpointSuccessfully(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	playerID := unittests.NewPlayerID()
	newPasswordValue := "newpwd"
	dateCreated := time.Now().UTC()
	updatePlayer := players.UpdatePlayer{
		ID:       playerID,
		Email:    unittests.NewEmailAddress(t, "focampo@anotheremail.com"),
		Password: &newPasswordValue,
	}

	existingPlayer := existingPlayerToUpdateFixture(t, &playerID, dateCreated)

	wantedResult := players.UpdatePlayerResult{}

	givenPlayerExistResult := players.PlayerExistResult{
		EmailExist:    false,
		NicknameExist: false,
	}

	dummyHash := []byte("$3a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")

	hasherMock := unittests.NewHasherMock()
	hasherMock.On("Hash", *updatePlayer.Password).Return(dummyHash, nil)

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetByID", ctx, playerID).Return(&existingPlayer, nil)
	storageMock.On("GetPlayersWithEmailOrNickName", ctx, mock.AnythingOfType("players.PlayerFilter")).Return(&givenPlayerExistResult, nil)
	storageMock.On("Update", ctx, mock.AnythingOfType("players.Player")).Return(nil)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent"))

	service, logger := unittests.NewPlayerService(storageMock, hasherMock, notifierMock)
	updatePlayerEndpoint := players.MakeUpdatePlayerEndpoint(service, logger)

	// When
	got, err := updatePlayerEndpoint.Do(ctx, &updatePlayer)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, wantedResult, got)
}

func TestUpdatePlayerWithEndpointButError(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	playerID := unittests.NewPlayerID()
	newPasswordValue := "newpwd"
	dateCreated := time.Now().UTC()
	updatePlayer := players.UpdatePlayer{
		ID:       playerID,
		Email:    unittests.NewEmailAddress(t, "focampo@anotheremail.com"),
		Password: &newPasswordValue,
	}

	existingPlayer := existingPlayerToUpdateFixture(t, &playerID, dateCreated)

	givenPlayerExistResult := players.PlayerExistResult{
		EmailExist:    false,
		NicknameExist: false,
	}

	saveError := errors.New("unexpected update error")

	dummyHash := []byte("$3a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")

	wantedResult := players.UpdatePlayerResult{
		Err: "unable to update player: unexpected update error",
	}

	hasherMock := unittests.NewHasherMock()
	hasherMock.On("Hash", *updatePlayer.Password).Return(dummyHash, nil)

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetByID", ctx, playerID).Return(&existingPlayer, nil)
	storageMock.On("GetPlayersWithEmailOrNickName", ctx, mock.AnythingOfType("players.PlayerFilter")).Return(&givenPlayerExistResult, nil)
	storageMock.On("Update", ctx, mock.AnythingOfType("players.Player")).Return(saveError)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent"))

	service, logger := unittests.NewPlayerService(storageMock, hasherMock, notifierMock)
	updatePlayerEndpoint := players.MakeUpdatePlayerEndpoint(service, logger)

	// When
	got, err := updatePlayerEndpoint.Do(ctx, &updatePlayer)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, wantedResult, got)
}

func existingPlayerToUpdateFixture(t *testing.T, playerID *players.PlayerID, dateTx time.Time) players.Player {
	t.Helper()

	return players.Player{
		ID:          playerID,
		FirstName:   "Fernando",
		LastName:    "Ocampo",
		Nickname:    "focampo",
		Email:       *unittests.NewEmailAddress(t, "focampo@anyemail.com"),
		Password:    []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6"),
		Country:     "Spain",
		DateCreated: dateTx,
		DateUpdated: dateTx,
	}
}
