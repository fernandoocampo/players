package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/fernandoocampo/players/internal/adapters/grpc"
	"github.com/fernandoocampo/players/internal/appkit/unittests"
	"github.com/fernandoocampo/players/internal/players"
	pb "github.com/fernandoocampo/players/pkg/pb/players"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePlayer(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	playerID := players.ToPlayerID(uuid.New())
	givenPlayerResult := players.Player{
		ID:          playerID,
		FirstName:   "Fernando",
		LastName:    "Ocampo",
		Nickname:    "focampo",
		Email:       *unittests.NewEmailAddress(t, "focampo@anyemail.com"),
		Password:    []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6"),
		Country:     "Spain",
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}
	newPlayerRequest := newCreatePlayerFixture()
	service := newServiceMock()
	service.On("Create", ctx, mock.AnythingOfType("players.NewPlayer")).Return(&givenPlayerResult, nil)
	server := newGRPCHandler(service)
	want := pb.CreatePlayerReply{
		PlayerId: playerID.String(),
	}

	// When
	reply, err := server.CreatePlayer(ctx, newPlayerRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, &want, reply)
}

func TestUpdatePlayer(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	playerID := players.ToPlayerID(uuid.New())
	givenPlayerResult := players.Player{
		ID:          playerID,
		FirstName:   "Fernando",
		LastName:    "Ocampo",
		Nickname:    "focampo",
		Email:       *unittests.NewEmailAddress(t, "focampo@anyemail.com"),
		Password:    []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6"),
		Country:     "Spain",
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}
	updatePlayerRequest := newUpdatePlayerFixture(playerID)
	service := newServiceMock()
	service.On("Update", ctx, mock.AnythingOfType("players.UpdatePlayer")).Return(&givenPlayerResult, nil)
	server := newGRPCHandler(service)
	want := pb.UpdatePlayerReply{
		Ok: true,
	}

	// When
	reply, err := server.UpdatePlayer(ctx, updatePlayerRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, &want, reply)
}

func TestDeletePlayer(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	playerID := uuid.New().String()
	deletePlayerRequest := pb.DeletePlayerRequest{
		PlayerId: playerID,
	}
	service := newServiceMock()
	service.On("Delete", ctx, mock.AnythingOfType("players.PlayerID")).Return(nil)
	server := newGRPCHandler(service)
	want := pb.DeletePlayerReply{
		Ok: true,
	}

	// When
	reply, err := server.DeletePlayer(ctx, &deletePlayerRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, &want, reply)
}

func TestSearchPlayer(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.TODO()
	searchRequest := pb.SearchPlayersRequest{
		Country: "UK",
		Limit:   3,
		Offset:  0,
	}
	givenSearchResult := players.SearchResult{
		Items:  unittests.SearchPlayersResultFixture(t),
		Total:  9,
		Limit:  3,
		Offset: 0,
	}
	service := newServiceMock()
	service.On("List", ctx, mock.AnythingOfType("players.SearchCriteria")).Return(&givenSearchResult, nil)
	server := newGRPCHandler(service)
	want := pb.SearchPlayersReply{
		PlayerItems: searchResultFixture(),
		Total:       9,
		Limit:       3,
		Offset:      0,
	}

	// When
	reply, err := server.SearchPlayers(ctx, &searchRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, &want, reply)
}

type MockService struct {
	mock.Mock
}

func newServiceMock() *MockService {
	return &MockService{}
}

func (m *MockService) Create(ctx context.Context, newPlayer players.NewPlayer) (*players.Player, error) {
	args := m.Called(ctx, newPlayer)

	return args.Get(0).(*players.Player), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, updatePlayer players.UpdatePlayer) (*players.Player, error) {
	args := m.Called(ctx, updatePlayer)

	return args.Get(0).(*players.Player), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, playerID players.PlayerID) error {
	args := m.Called(ctx, playerID)

	return args.Error(0)
}

func (m *MockService) List(ctx context.Context, searchCriteria players.SearchCriteria) (*players.SearchResult, error) {
	args := m.Called(ctx, searchCriteria)

	return args.Get(0).(*players.SearchResult), args.Error(1)
}

func newGRPCHandler(service *MockService) *grpc.Handler {
	handlerSetup := grpc.HandlerSetup{
		Service: service,
		Logger:  unittests.NewLogger(),
	}

	return grpc.NewHandler(handlerSetup)
}

func newCreatePlayerFixture() *pb.CreatePlayerRequest {
	return &pb.CreatePlayerRequest{
		Firstname: "Fernando",
		Lastname:  "Ocampo",
		Nickname:  "focampo",
		Email:     "focampo@anyemail.com",
		Password:  "th1s1s@dummypwd",
		Country:   "Spain",
	}
}

func newUpdatePlayerFixture(playerID *players.PlayerID) *pb.UpdatePlayerRequest {
	return &pb.UpdatePlayerRequest{
		PlayerId:  playerID.String(),
		Firstname: "Fernando",
		Lastname:  "Ocampo",
	}
}

func searchResultFixture() []*pb.PlayerItem {
	return []*pb.PlayerItem{
		{
			Id:        "b10d6af5-22f3-4db2-ade6-94cfcc819f91",
			Firstname: "Jim",
			Lastname:  "Raynor",
			Nickname:  "jraynor",
			Country:   "UK",
		},
		{
			Id:        "8b1cb38a-d21c-4e36-b7f0-814f96180b5e",
			Firstname: "Sarah",
			Lastname:  "Kerrigan",
			Nickname:  "skerrigan",
			Country:   "UK",
		},
		{
			Id:        "f7ee12ee-191d-4f7e-819e-dbb7b115bed2",
			Firstname: "Graven",
			Lastname:  "Hill",
			Nickname:  "ghill",
			Country:   "UK",
		},
	}
}
