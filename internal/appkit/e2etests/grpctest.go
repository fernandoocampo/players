package e2etests

import (
	"fmt"
	"testing"

	pb "github.com/fernandoocampo/players/pkg/pb/players"
)

type GRPCServerParameters struct {
	Address string
}

func GRPCServerParametersFixture(t *testing.T) GRPCServerParameters {
	t.Helper()

	parameters := GRPCServerParameters{
		Address: getStringEnvVar(t, "PLAYERS_E2E_GRPC_SERVER_ADDRESS", "localhost:50051"),
	}

	return parameters
}

func RandomPBCreatePlayerFixture() *pb.CreatePlayerRequest {
	anyPassword := "th1s1s@dummypwd"
	newNickname := randomString(10)
	newEmail := fmt.Sprintf("%s@e2etest.com", newNickname)

	return &pb.CreatePlayerRequest{
		Firstname: "e2e",
		Lastname:  "grpctest",
		Nickname:  newNickname,
		Email:     newEmail,
		Password:  anyPassword,
		Country:   randomCountry(),
	}
}
