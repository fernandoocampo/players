package grpc

import (
	"context"
	"log/slog"

	"github.com/fernandoocampo/players/internal/players"
	pb "github.com/fernandoocampo/players/pkg/pb/players"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlayerService interface {
	Create(ctx context.Context, newPlayer players.NewPlayer) (*players.Player, error)
	Update(ctx context.Context, updatePlayer players.UpdatePlayer) (*players.Player, error)
	Delete(ctx context.Context, playerID players.PlayerID) error
	List(ctx context.Context, searchCriteria players.SearchCriteria) (*players.SearchResult, error)
}

type HandlerSetup struct {
	Service PlayerService
	Logger  *slog.Logger
}

type Handler struct {
	pb.UnimplementedPlayerHandlerServer
	service PlayerService
	logger  *slog.Logger
}

func NewHandler(setup HandlerSetup) *Handler {
	newHandler := Handler{
		service: setup.Service,
		logger:  setup.Logger,
	}

	return &newHandler
}

// CreatePlayer creates a player.
func (s *Handler) CreatePlayer(ctx context.Context, request *pb.CreatePlayerRequest) (*pb.CreatePlayerReply, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	players, err := s.service.Create(ctx, toNewPlayer(request))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toCreatePlayerReply(players.ID), nil
}

// UpdatePlayer updates a player.
func (s *Handler) UpdatePlayer(ctx context.Context, request *pb.UpdatePlayerRequest) (*pb.UpdatePlayerReply, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be empty")
	}

	if request.GetPlayerId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "request missing required field: player id")
	}

	playerID, err := players.StringToPlayerID(request.GetPlayerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "request has invalid player id, it must be a uuid")
	}

	_, err = s.service.Update(ctx, toUpdatePlayer(request, playerID))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return updatePlayerReplyOK(), nil
}

// DeletePlayer deletes a player.
func (s *Handler) DeletePlayer(ctx context.Context, request *pb.DeletePlayerRequest) (*pb.DeletePlayerReply, error) {
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be empty")
	}

	playerID, err := players.StringToPlayerID(request.GetPlayerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "request has invalid player id, it must be a uuid")
	}

	err = s.service.Delete(ctx, *playerID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return deletePlayerReplyOK(), nil
}

// SearchPlayers searches a player.
func (s *Handler) SearchPlayers(ctx context.Context, request *pb.SearchPlayersRequest) (*pb.SearchPlayersReply, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	result, err := s.service.List(ctx, toSearchCriteria(request))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toSearchPlayerReply(result), nil
}
