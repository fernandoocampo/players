package storages_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/fernandoocampo/players/internal/adapters/storages"
	"github.com/fernandoocampo/players/internal/appkit/e2etests"
	"github.com/fernandoocampo/players/internal/appkit/unittests"
	"github.com/fernandoocampo/players/internal/players"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresqlConnection(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test, to execute this test send e2e-test flag to true")
	}

	t.Log("executing integration test to create a postgresql client")

	// Given
	givenDBParameters := e2etests.DatabaseParametersFixture(t)

	// When
	client, err := storages.NewPostgresClient(givenDBParameters)

	// Then
	assert.NoError(t, err)
	if err != nil {
		t.Errorf("unexpected error trying to connect to default database: %s", err)
		t.FailNow()
	}
	err = storages.CloseConnection(client)

	assert.NoError(t, err)
}

func TestSavePlayer(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	newPlayer := e2etests.RandomPlayerFixture()
	ctx := context.TODO()

	storage, client := newStorage(t)
	defer closeConnection(t, client)

	// When
	err := storage.Save(ctx, newPlayer)

	// Then
	assert.NoError(t, err)
}

func TestUpdatePlayer(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	storage, client := newStorage(t)
	defer closeConnection(t, client)

	newPlayer := savePlayer(ctx, t, storage)

	newPlayer.LastName = "updated"
	newPlayer.DateUpdated = time.Now().UTC()

	// When
	err := storage.Update(ctx, *newPlayer)

	// Then
	assert.NoError(t, err)
}

func TestDeletePlayer(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	storage, client := newStorage(t)
	defer closeConnection(t, client)

	newPlayer := savePlayer(ctx, t, storage)

	// When
	err := storage.Delete(ctx, *newPlayer.ID)

	// Then
	assert.NoError(t, err)
}

func TestGetPlayerByID(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	storage, client := newStorage(t)
	defer closeConnection(t, client)

	newPlayer := savePlayer(ctx, t, storage)

	playerIDToFind := *newPlayer.ID

	// When
	got, err := storage.GetByID(ctx, playerIDToFind)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, newPlayer, got)
}

func TestGetPlayersByNicknameOrEmail(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	storage, client := newStorage(t)
	defer closeConnection(t, client)

	filter := players.PlayerFilter{
		Email:    "unknownplayer@anyemail.com",
		Nickname: "unknownplayer",
	}

	want := players.PlayerExistResult{
		EmailExist:    false,
		NicknameExist: false,
	}

	// When
	got, err := storage.GetPlayersWithEmailOrNickName(ctx, filter)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, &want, got)
}

func TestGetPlayersByNicknameOrEmailBothExist(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	storage, client := newStorage(t)
	defer closeConnection(t, client)

	newPlayer1 := savePlayer(ctx, t, storage)

	filter := players.PlayerFilter{
		Email:    newPlayer1.Email.Address,
		Nickname: newPlayer1.Nickname,
	}

	want := players.PlayerExistResult{
		EmailExist:    true,
		NicknameExist: true,
	}

	// When
	got, err := storage.GetPlayersWithEmailOrNickName(ctx, filter)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, &want, got)
}

func TestGetPlayersByNicknameOrEmailIgnore(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	storage, client := newStorage(t)
	defer closeConnection(t, client)

	newPlayer1 := savePlayer(ctx, t, storage)

	filter := players.PlayerFilter{
		Email:    newPlayer1.Email.Address,
		Nickname: newPlayer1.Nickname,
		IgnoreID: newPlayer1.ID,
	}

	want := players.PlayerExistResult{
		EmailExist:    false,
		NicknameExist: false,
	}

	// When
	got, err := storage.GetPlayersWithEmailOrNickName(ctx, filter)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, &want, got)
}

func TestSearchPlayers(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	storage, client := newStorage(t)
	defer closeConnection(t, client)

	countryFilter := "Colombia"
	searchCriteria := players.SearchCriteria{
		Country: &countryFilter,
		Limit:   5,
		Offset:  0,
	}

	for range 5 { // Save 5 players with the same country
		newPlayer := e2etests.RandomPlayerFixture()
		newPlayer.Country = "Colombia"
		err := storage.Save(ctx, newPlayer)
		require.NoError(t, err)
	}

	// When
	got, err := storage.Search(ctx, searchCriteria)

	// Then
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, got.Total, 5)
	assert.Equal(t, 5, len(got.Items))
}

func newStorage(t *testing.T) (*storages.Storage, *sql.DB) {
	t.Helper()

	givenDBParameters := e2etests.DatabaseParametersFixture(t)

	client, err := storages.NewPostgresClient(givenDBParameters)
	require.NoError(t, err)

	storageSetup := storages.StorageSetup{
		DB:     client,
		Logger: unittests.NewLogger(),
	}

	storage := storages.NewPlayerRepository(storageSetup)

	return storage, client
}

func savePlayer(ctx context.Context, t *testing.T, storage *storages.Storage) *players.Player {
	t.Helper()

	newPlayer := e2etests.RandomPlayerFixture()

	t.Logf("saving player: %+v", newPlayer)

	err := storage.Save(ctx, newPlayer)
	require.NoError(t, err)

	return &newPlayer
}

func closeConnection(t *testing.T, conn *sql.DB) {
	t.Helper()

	err := storages.CloseConnection(conn)
	if err != nil {
		t.Logf("unexpected error closing connection: %s", err)
	}
}
