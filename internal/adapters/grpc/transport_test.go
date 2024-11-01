package grpc_test

import (
	"context"
	"testing"

	"github.com/fernandoocampo/players/internal/appkit/e2etests"
	pb "github.com/fernandoocampo/players/pkg/pb/players"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestE2ECreatePlayer(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	request := e2etests.RandomPBCreatePlayerFixture()

	grpcClient := createPlayerClient(t)
	defer grpcClient.Close(t)

	// when
	reply, err := grpcClient.client.CreatePlayer(ctx, request)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, reply)
	assert.NotEmpty(t, reply.PlayerId)
}

// UpdatePlayer updates a player.
func TestE2EUpdatePlayer(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	createRequest := e2etests.RandomPBCreatePlayerFixture()

	grpcClient := createPlayerClient(t)
	defer grpcClient.Close(t)

	createReply, err := grpcClient.client.CreatePlayer(ctx, createRequest)
	require.NoError(t, err)
	require.NotEmpty(t, createReply.PlayerId)

	t.Logf("updating player: %s", createReply.GetPlayerId())

	updateRequest := pb.UpdatePlayerRequest{
		PlayerId: createReply.PlayerId,
		Lastname: "e2eupdated",
	}

	// when
	reply, err := grpcClient.client.UpdatePlayer(ctx, &updateRequest)

	// Then
	assert.NoError(t, err)
	assert.True(t, reply.Ok)
}

// DeletePlayer deletes a player.
func TestE2EDeletePlayer(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	createRequest := e2etests.RandomPBCreatePlayerFixture()

	grpcClient := createPlayerClient(t)
	defer grpcClient.Close(t)

	createReply, err := grpcClient.client.CreatePlayer(ctx, createRequest)
	require.NoError(t, err)

	deleteRequest := pb.DeletePlayerRequest{
		PlayerId: createReply.PlayerId,
	}

	// when
	reply, err := grpcClient.client.DeletePlayer(ctx, &deleteRequest)

	// Then
	assert.NoError(t, err)
	assert.True(t, reply.Ok)
}

// SearchPlayers searches a player.
func TestE2ESearchPlayers(t *testing.T) {
	if !*e2etests.E2ETest {
		t.Skip("this is an e2e test to verify database calls, to execute this test send e2e-test flag to true")
	}

	// Given
	ctx := context.TODO()

	requests := []*pb.CreatePlayerRequest{
		e2etests.RandomPBCreatePlayerFixture(),
		e2etests.RandomPBCreatePlayerFixture(),
		e2etests.RandomPBCreatePlayerFixture(),
	}

	grpcClient := createPlayerClient(t)
	defer grpcClient.Close(t)

	for _, request := range requests {
		request.Country = "UK"
		_, err := grpcClient.client.CreatePlayer(ctx, request)
		require.NoError(t, err)
	}

	searchRequest := pb.SearchPlayersRequest{
		Country: "UK",
		Limit:   3,
		Offset:  0,
	}

	// when
	reply, err := grpcClient.client.SearchPlayers(ctx, &searchRequest)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, reply.GetPlayerItems())
	assert.Equal(t, uint32(3), reply.GetLimit())
	assert.Equal(t, 3, len(reply.PlayerItems))
}

func createPlayerClient(t *testing.T) *playerClient {
	t.Helper()

	setup := e2etests.GRPCServerParametersFixture(t)
	// Set up a connection to the server.
	conn, err := grpc.NewClient(setup.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := pb.NewPlayerHandlerClient(conn)

	newPlayerClient := playerClient{
		conn:   conn,
		client: client,
	}

	return &newPlayerClient
}

func (u *playerClient) Close(t *testing.T) {
	t.Helper()

	err := u.conn.Close()
	if err != nil {
		t.Logf("unable to close grpc connection to server: %s", err)
	}
}

type PlayerClientSetup struct {
	Address string
}

type playerClient struct {
	conn   *grpc.ClientConn
	client pb.PlayerHandlerClient
}
