package players

import (
	"context"
	"fmt"
	"log/slog"
)

// ServiceSetup encapsulates service parameters.
type ServiceSetup struct {
	// Storage player repository
	Storage Storage
	// Notifier players event notifier
	Notifier Notifier
	// password hasher
	Hasher Hasher
	Logger *slog.Logger
}

// Service defines business logic for this service.
type Service struct {
	// storage player repository
	storage Storage
	// password hasher
	hasher Hasher
	// player events notifier
	notifier Notifier
	logger   *slog.Logger
}

// NewService create a new service instance.
func NewService(setup *ServiceSetup) *Service {
	newService := Service{
		storage:  setup.Storage,
		hasher:   setup.Hasher,
		notifier: setup.Notifier,
		logger:   setup.Logger,
	}

	return &newService
}

// Create creates a new player with the given data. Validates new player information and verify
// that other players with same email or nickname already exist.
func (s *Service) Create(ctx context.Context, newPlayer NewPlayer) (*Player, error) {
	s.logger.Debug("starting to create a new player", slog.Any("new_player", newPlayer.obfuscate()))

	err := newPlayer.Validate()
	if err != nil {
		return nil, fmt.Errorf("unable to create player: %w", err)
	}

	hashedPassword, err := s.hasher.Hash(newPlayer.Password)
	if err != nil {
		s.logger.Error("hashing password", "error", err.Error())

		return nil, fmt.Errorf("unable to hash password: %w", err)
	}

	checkExistResult, err := s.doesThePlayerAlreadyExist(ctx, newPlayer.toPlayerFilter())
	if err != nil {
		return nil, fmt.Errorf("unable to create player: %w", err)
	}

	if checkExistResult.Exist() {
		s.logger.Debug(
			"player with the given email or nickname already exists",
			slog.String("email", newPlayer.Email.Address),
			slog.String("nickname", newPlayer.Nickname),
		)

		return nil, checkExistResult.toPlayerAlreadyExistsError()
	}

	player := newPlayer.toPlayer(hashedPassword)

	err = s.storage.Save(ctx, player)
	if err != nil {
		s.logger.Error("creating player", "error", err)

		return nil, fmt.Errorf("unable to create player: %w", err)
	}

	s.logger.Debug("new player was created", slog.Any("id", player.ID))

	s.notifier.Notify(newCreatePlayerEvent(player.ID))

	return &player, nil
}

// Update updates existing player with the given data. Validates the new player data and verify
// that other players with same email or nickname already exist.
func (s *Service) Update(ctx context.Context, updatePlayer UpdatePlayer) (*Player, error) {
	s.logger.Debug("starting to update player", slog.Any("player", updatePlayer.obfuscate()))

	err := updatePlayer.Validate()
	if err != nil {
		return nil, fmt.Errorf("unable to update player: %w", err)
	}

	if updatePlayer.updateKeyValues() {
		playerExistResult, err := s.doesThePlayerAlreadyExist(ctx, updatePlayer.toPlayerFilter())
		if err != nil {
			return nil, fmt.Errorf("unable to update player: %w", err)
		}

		if playerExistResult.Exist() {
			return nil, playerExistResult.toPlayerAlreadyExistsError()
		}
	}

	player, err := s.storage.GetByID(ctx, updatePlayer.ID)
	if err != nil {
		s.logger.Error("getting player by id", slog.String("id", updatePlayer.ID.String()), slog.String("error", err.Error()))

		return nil, fmt.Errorf("unable to update player: %w", err)
	}

	if player == nil {
		return nil, ErrPlayerDoesNotExist
	}

	playerToUpdate, err := updatePlayer.toPlayer(*player, s.hasher)
	if err != nil {
		s.logger.Error("converting updateplayer to player", slog.String("id", updatePlayer.ID.String()), slog.String("error", err.Error()))

		return nil, fmt.Errorf("unable to update player: %w", err)
	}

	if !playerToUpdate.changes {
		s.logger.Debug("there is nothing to update in the player", slog.Any("player", player.obfuscate()))

		return player, nil
	}

	err = s.storage.Update(ctx, *playerToUpdate.player)
	if err != nil {
		s.logger.Error("updating player", "error", err)

		return nil, fmt.Errorf("unable to update player: %w", err)
	}

	s.logger.Debug("player was updated", slog.Any("id", player.ID))

	s.notifier.Notify(newUpdatePlayerEvent(player.ID))

	return playerToUpdate.player, nil
}

func (s *Service) Delete(ctx context.Context, playerID PlayerID) error {
	s.logger.Debug("starting to delete player", slog.Any("player_id", playerID.String()))

	player, err := s.storage.GetByID(ctx, playerID)
	if err != nil {
		s.logger.Error("getting player by id", slog.String("id", playerID.String()), slog.String("error", err.Error()))

		return fmt.Errorf("unable to delete player: %w", err)
	}

	if player == nil {
		return ErrPlayerDoesNotExist
	}

	err = s.storage.Delete(ctx, playerID)
	if err != nil {
		s.logger.Error("deleting player", "error", err)

		return fmt.Errorf("unable to delete player: %w", err)
	}

	s.logger.Debug("player was deleted", slog.Any("id", playerID))

	s.notifier.Notify(newDeletePlayerEvent(playerID))

	return nil
}

func (s *Service) doesThePlayerAlreadyExist(ctx context.Context, playerFilter PlayerFilter) (*PlayerExistResult, error) {
	s.logger.Debug(
		"checking if a player with the given email and nickname already exists",
		slog.String("email", playerFilter.Email),
		slog.String("nickname", playerFilter.Nickname),
	)

	playerWithEmailOrNickname, err := s.storage.GetPlayersWithEmailOrNickName(ctx, playerFilter)
	if err != nil {
		s.logger.Error(
			"checking if a player with the given email or nickname already exists",
			slog.String("email", playerFilter.Email),
			slog.String("nickname", playerFilter.Nickname),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("unable to check if player already exists: %w", err)
	}

	return playerWithEmailOrNickname, nil
}

func (s *Service) List(ctx context.Context, searchCriteria SearchCriteria) (*SearchResult, error) {
	s.logger.Debug("starting to search players", slog.Any("criteria", searchCriteria))

	if searchCriteria.isEmpty() {
		return newEmptySearchResult(), nil
	}

	searchCriteria.setDefaultPaginationIfEmpty()

	result, err := s.storage.Search(ctx, searchCriteria)
	if err != nil {
		s.logger.Error("searching players", slog.String("error", err.Error()))

		return nil, fmt.Errorf("unable to list players: %w", err)
	}

	return result, nil
}
