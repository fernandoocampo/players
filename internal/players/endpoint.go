package players

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type CreatePlayerEndpoint struct {
	service *Service
	logger  *slog.Logger
}

type UpdatePlayerEndpoint struct {
	service *Service
	logger  *slog.Logger
}

type DeletePlayerEndpoint struct {
	service *Service
	logger  *slog.Logger
}

type SearchPlayersEndpoint struct {
	service *Service
	logger  *slog.Logger
}

// Endpoints is a wrapper for endpoints.
type Endpoints struct {
	CreatePlayerEndpoint  *CreatePlayerEndpoint
	UpdatePlayerEndpoint  *UpdatePlayerEndpoint
	DeletePlayerEndpoint  *DeletePlayerEndpoint
	SearchPlayersEndpoint *SearchPlayersEndpoint
}

var (
	errInvalidNewPlayerType      = errors.New("invalid new player type")
	errInvalidUpdatePlayerType   = errors.New("invalid update player type")
	errInvalidPlayerIDType       = errors.New("invalid player id type")
	errInvalidSearchCriteriaType = errors.New("invalid search players type")
)

// NewEndpoints Create the endpoints for player application.
func NewEndpoints(service *Service, logger *slog.Logger) Endpoints {
	return Endpoints{
		CreatePlayerEndpoint:  MakeCreatePlayerEndpoint(service, logger),
		UpdatePlayerEndpoint:  MakeUpdatePlayerEndpoint(service, logger),
		DeletePlayerEndpoint:  MakeDeletePlayerEndpoint(service, logger),
		SearchPlayersEndpoint: MakeSearchPlayersEndpoint(service, logger),
	}
}

// MakeCreatePlayerEndpoint create endpoint for create player service.
func MakeCreatePlayerEndpoint(srv *Service, logger *slog.Logger) *CreatePlayerEndpoint {
	newNewEndpoint := CreatePlayerEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeUpdatePlayerEndpoint create endpoint for update player service.
func MakeUpdatePlayerEndpoint(srv *Service, logger *slog.Logger) *UpdatePlayerEndpoint {
	newNewEndpoint := UpdatePlayerEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeDeletePlayerEndpoint create endpoint for the delete player service.
func MakeDeletePlayerEndpoint(srv *Service, logger *slog.Logger) *DeletePlayerEndpoint {
	newNewEndpoint := DeletePlayerEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeSearchPlayersEndpoint player endpoint to search players with filters.
func MakeSearchPlayersEndpoint(srv *Service, logger *slog.Logger) *SearchPlayersEndpoint {
	newNewEndpoint := SearchPlayersEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

func (c *CreatePlayerEndpoint) Do(ctx context.Context, request any) (any, error) {
	newPlayer, ok := request.(*NewPlayer)
	if !ok {
		c.logger.Error("invalid new player type", slog.String("request", fmt.Sprintf("%t", request)))

		return nil, errInvalidNewPlayerType
	}

	player, err := c.service.Create(ctx, *newPlayer)
	if err != nil {
		c.logger.Error(
			"creating player",
			slog.Any("new_player", newPlayer.obfuscate()),
			slog.String("error", err.Error()),
		)

		return newCreatePlayerResult(nil, err), nil
	}

	return newCreatePlayerResult(player.ID, err), nil
}

func (u *UpdatePlayerEndpoint) Do(ctx context.Context, request any) (any, error) {
	updatePlayer, ok := request.(*UpdatePlayer)
	if !ok {
		u.logger.Error("invalid update player type", slog.String("request", fmt.Sprintf("%t", request)))

		return nil, errInvalidUpdatePlayerType
	}

	_, err := u.service.Update(ctx, *updatePlayer)
	if err != nil {
		u.logger.Error(
			"updating a player",
			slog.Any("player", updatePlayer.obfuscate()),
			slog.String("error", err.Error()),
		)
	}

	return newUpdatePlayerResult(err), nil
}

func (d *DeletePlayerEndpoint) Do(ctx context.Context, request any) (any, error) {
	playerID, ok := request.(PlayerID)
	if !ok {
		d.logger.Error("invalid delete player type", slog.String("received", fmt.Sprintf("%t", request)))

		return nil, errInvalidPlayerIDType
	}

	err := d.service.Delete(ctx, playerID)
	if err != nil {
		d.logger.Error(
			"deleting player with the given id",
			slog.String("id", playerID.String()),
			slog.String("error", err.Error()),
		)
	}

	return newDeletePlayerResult(err), nil
}

func (s *SearchPlayersEndpoint) Do(ctx context.Context, request any) (any, error) {
	searchCriteria, ok := request.(SearchCriteria)
	if !ok {
		s.logger.Error("invalid search criteria request", slog.String("received", fmt.Sprintf("%t", request)))

		return nil, errInvalidSearchCriteriaType
	}

	searchResult, err := s.service.List(ctx, searchCriteria)
	if err != nil {
		s.logger.Error(
			"querying players with the given filter",
			slog.Any("filters", searchCriteria),
			slog.String("error", err.Error()),
		)
	}

	return newSearchPlayersDataResult(searchResult, err), nil
}
