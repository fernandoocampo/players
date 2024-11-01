package players_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fernandoocampo/players/internal/appkit/unittests"
	"github.com/fernandoocampo/players/internal/players"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePlayer(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	newPlayer := unittests.NewPlayerFixture(t)

	givenPlayerExistResult := players.PlayerExistResult{
		EmailExist:    false,
		NicknameExist: false,
	}

	dummyHash := []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")

	want := players.Player{
		FirstName: "Fernando",
		LastName:  "Ocampo",
		Nickname:  "focampo",
		Email:     *unittests.NewEmailAddress(t, "focampo@anyemail.com"),
		Password:  []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6"),
		Country:   "Spain",
	}

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetPlayersWithEmailOrNickName", ctx, mock.AnythingOfType("players.PlayerFilter")).Return(&givenPlayerExistResult, nil)
	storageMock.On("Save", ctx, mock.AnythingOfType("players.Player")).Return(nil)

	hasherMock := unittests.NewHasherMock()
	hasherMock.On("Hash", newPlayer.Password).Return(dummyHash, nil)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent"))

	service, _ := unittests.NewPlayerService(storageMock, hasherMock, notifierMock)

	// When
	player, err := service.Create(ctx, newPlayer)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, player)
	assert.NotNil(t, player.ID)
	assert.NotNil(t, player.DateCreated)
	want.ID = player.ID                   // random values
	want.DateCreated = player.DateCreated // random values
	want.DateUpdated = player.DateUpdated // random values
	assert.Equal(t, &want, player)
}

func TestCreateButPlayerWithEmailAlreadyExists(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	newPlayer := unittests.NewPlayerFixture(t)

	givenPlayerExistResult := players.PlayerExistResult{
		EmailExist:    true,
		NicknameExist: false,
	}

	dummyHash := []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")

	want := "player with the given email already exists"

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetPlayersWithEmailOrNickName", ctx, mock.AnythingOfType("players.PlayerFilter")).Return(&givenPlayerExistResult, nil)

	hasherMock := unittests.NewHasherMock()
	hasherMock.On("Hash", newPlayer.Password).Return(dummyHash, nil)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent"))

	service, _ := unittests.NewPlayerService(storageMock, hasherMock, notifierMock)

	// When
	player, err := service.Create(ctx, newPlayer)

	// Then
	assert.Error(t, err)
	assert.Nil(t, player)
	assert.Equal(t, want, err.Error())
}

func TestCreateButPlayerWithNicknameAlreadyExists(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	newPlayer := unittests.NewPlayerFixture(t)

	givenPlayerExistResult := players.PlayerExistResult{
		EmailExist:    false,
		NicknameExist: true,
	}

	dummyHash := []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")

	want := "player with the given nickname already exists"

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetPlayersWithEmailOrNickName", ctx, mock.AnythingOfType("players.PlayerFilter")).Return(&givenPlayerExistResult, nil)

	hasherMock := unittests.NewHasherMock()
	hasherMock.On("Hash", newPlayer.Password).Return(dummyHash, nil)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent"))

	service, _ := unittests.NewPlayerService(storageMock, hasherMock, notifierMock)

	// When
	player, err := service.Create(ctx, newPlayer)

	// Then
	assert.Error(t, err)
	assert.Nil(t, player)
	assert.Equal(t, want, err.Error())
}

func TestCreateButPlayerWithNicknameAndEmailAlreadyExists(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	newPlayer := unittests.NewPlayerFixture(t)

	givenPlayerExistResult := players.PlayerExistResult{
		EmailExist:    true,
		NicknameExist: true,
	}

	want := "player with the given email or nickname already exists"

	dummyHash := []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetPlayersWithEmailOrNickName", ctx, mock.AnythingOfType("players.PlayerFilter")).Return(&givenPlayerExistResult, nil)

	hasherMock := unittests.NewHasherMock()
	hasherMock.On("Hash", newPlayer.Password).Return(dummyHash, nil)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent"))

	service, _ := unittests.NewPlayerService(storageMock, hasherMock, notifierMock)

	// When
	player, err := service.Create(ctx, newPlayer)

	// Then
	assert.Error(t, err)
	assert.Nil(t, player)
	assert.Equal(t, want, err.Error())
}

func TestCreatePlayerWithError(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	newPlayer := unittests.NewPlayerFixture(t)

	givenPlayerExistResult := players.PlayerExistResult{
		EmailExist:    false,
		NicknameExist: false,
	}

	saveError := errors.New("unexpected create error")

	dummyHash := []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")

	want := "unable to create player: unexpected create error"

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetPlayersWithEmailOrNickName", ctx, mock.AnythingOfType("players.PlayerFilter")).Return(&givenPlayerExistResult, nil)
	storageMock.On("Save", ctx, mock.AnythingOfType("players.Player")).Return(saveError)

	hasherMock := unittests.NewHasherMock()
	hasherMock.On("Hash", newPlayer.Password).Return(dummyHash, nil)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent"))

	service, _ := unittests.NewPlayerService(storageMock, hasherMock, notifierMock)

	// When
	player, err := service.Create(ctx, newPlayer)

	// Then
	assert.Error(t, err)
	assert.Nil(t, player)
	assert.Equal(t, want, err.Error())
}

func TestCreatePlayerWithEndpointSuccessfully(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	newPlayer := unittests.NewPlayerFixture(t)

	givenPlayerExistResult := players.PlayerExistResult{
		EmailExist:    false,
		NicknameExist: false,
	}

	dummyHash := []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")

	want := players.CreatePlayerResult{}

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetPlayersWithEmailOrNickName", ctx, mock.AnythingOfType("players.PlayerFilter")).Return(&givenPlayerExistResult, nil)
	storageMock.On("Save", ctx, mock.AnythingOfType("players.Player")).Return(nil)

	hasherMock := unittests.NewHasherMock()
	hasherMock.On("Hash", newPlayer.Password).Return(dummyHash, nil)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent"))

	service, logger := unittests.NewPlayerService(storageMock, hasherMock, notifierMock)
	createPlayerEndpoint := players.MakeCreatePlayerEndpoint(service, logger)

	// When
	got, err := createPlayerEndpoint.Do(ctx, &newPlayer)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, got)
	want.ID = got.(players.CreatePlayerResult).ID // random values
	assert.Equal(t, want, got)
}

func TestCreatePlayerWithEndpointButError(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	newPlayer := unittests.NewPlayerFixture(t)

	givenPlayerExistResult := players.PlayerExistResult{
		EmailExist:    false,
		NicknameExist: false,
	}

	dummyHash := []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")

	saveError := errors.New("unexpected create error")

	want := players.CreatePlayerResult{
		Err: "unable to create player: unexpected create error",
	}

	storageMock := unittests.NewStorageMock()
	storageMock.On("GetPlayersWithEmailOrNickName", ctx, mock.AnythingOfType("players.PlayerFilter")).Return(&givenPlayerExistResult, nil)
	storageMock.On("Save", ctx, mock.AnythingOfType("players.Player")).Return(saveError)

	hasherMock := unittests.NewHasherMock()
	hasherMock.On("Hash", newPlayer.Password).Return(dummyHash, nil)

	notifierMock := unittests.NewNotifierMock()
	notifierMock.On("Notify", mock.AnythingOfType("players.NewEvent"))

	service, logger := unittests.NewPlayerService(storageMock, hasherMock, notifierMock)
	createPlayerEndpoint := players.MakeCreatePlayerEndpoint(service, logger)

	// When
	got, err := createPlayerEndpoint.Do(ctx, &newPlayer)

	// Then
	assert.NoError(t, err)
	assert.Nil(t, got.(players.CreatePlayerResult).ID)
	assert.Equal(t, want, got)
}
