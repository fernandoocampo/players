package storages

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fernandoocampo/players/internal/players"
	"github.com/google/uuid"
)

type StorageSetup struct {
	DB     *sql.DB
	Logger *slog.Logger
}

// Storage is the repository handler for this application in a relational db.
type Storage struct {
	db     *sql.DB
	logger *slog.Logger
}

// Queries.
const (
	createPlayerSQL = `INSERT INTO players(
	id,firstname,lastname,nickname,email,usrpwd,country,date_created,date_updated) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	updatePlayerSQL = `UPDATE players 
	SET firstname = $1,
	lastname = $2,
	nickname = $3,
	email = $4,
	country = $5, 
	usrpwd = $6,
	date_updated = $7 
	WHERE id = $8`
	deletePlayerSQL = "DELETE FROM players WHERE id = $1"
	selectByIDSQL   = `SELECT id,firstname,lastname,nickname,email,usrpwd,country,date_created,date_updated 
	FROM players 
	WHERE id = $1`
	selectByNicknameAndEmail = `SELECT (SELECT COUNT(*) FROM players WHERE nickname = $1) AS count_nickname, 
	(SELECT COUNT(*) FROM players WHERE email = $2) AS count_email`
	selectByNicknameAndEmailAndID = `SELECT (SELECT COUNT(*) FROM players WHERE nickname = $1 AND id <> $3) AS count_nickname, 
	(SELECT COUNT(*) FROM players WHERE email = $2 AND id <> $3) AS count_email`
	selectByFilterSQL = "SELECT id, firstname, lastname, nickname, country FROM players %s;"
	countByFilterSQL  = "SELECT COUNT(id) FROM players %s;"
)

// Error messages.
var (
	errPlayerCannotBeStored  = errors.New("player cannot be stored")
	errPlayerCannotBeUpdated = errors.New("player cannot be updated")
	errPlayerCannotBeDeleted = errors.New("player cannot be deleted")
	errPlayerCannotBeRead    = errors.New("player cannot be read in the database")
	errPlayersCannotBeRead   = errors.New("players cannot be read in the database")
	errUnableToSearchPlayers = errors.New("unable to search players")
)

// NewPlayerRepository creates a new player repository that will use a rdb.
func NewPlayerRepository(setup StorageSetup) *Storage {
	newStorage := Storage{
		db:     setup.DB,
		logger: setup.Logger,
	}

	return &newStorage
}

// Save persists a new player in the player repository.
func (s *Storage) Save(ctx context.Context, newPlayer players.Player) error {
	s.logger.Debug("storing player", slog.String("player_id", newPlayer.ID.String()))

	stmt, err := s.db.Prepare(createPlayerSQL)
	if err != nil {
		s.logger.Error("building query",
			slog.String("player_id", newPlayer.ID.String()),
			slog.String("error", err.Error()))

		return errPlayerCannotBeStored
	}

	defer stmt.Close()

	player := toDBPlayer(&newPlayer)

	_, err = stmt.ExecContext(ctx,
		player.ID.String(), player.FirstName, player.LastName,
		player.Nickname, player.Email, player.Password, player.Country,
		player.DateCreated, player.DateUpdated,
	)
	if err != nil {
		s.logger.Error("executing insert to store player",
			slog.String("player_id", player.ID.String()),
			slog.String("error", err.Error()))

		return errPlayerCannotBeStored
	}

	return nil
}

// Update player in the player repository.
func (s *Storage) Update(ctx context.Context, player players.Player) error {
	s.logger.Debug("updating player", slog.String("player_id", player.ID.String()))

	stmt, err := s.db.Prepare(updatePlayerSQL)
	if err != nil {
		s.logger.Error("building query",
			slog.String("player_id", player.ID.String()),
			slog.String("error", err.Error()))

		return errPlayerCannotBeUpdated
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		player.FirstName, player.LastName,
		player.Nickname, player.Email.Address, player.Country, player.Password,
		player.DateUpdated, player.ID.String(),
	)
	if err != nil {
		s.logger.Error("executing to update player",
			slog.String("player_id", player.ID.String()),
			slog.String("error", err.Error()))

		return errPlayerCannotBeUpdated
	}

	return nil
}

// Delete player in the repository.
func (s *Storage) Delete(ctx context.Context, playerID players.PlayerID) error {
	s.logger.Debug("deleting player", slog.String("player_id", playerID.String()))

	stmt, err := s.db.Prepare(deletePlayerSQL)
	if err != nil {
		s.logger.Error("building query",
			slog.String("player_id", playerID.String()),
			slog.String("error", err.Error()))

		return errPlayerCannotBeDeleted
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		playerID.String(),
	)
	if err != nil {
		s.logger.Error("executing to delete player",
			slog.String("player_id", playerID.String()),
			slog.String("error", err.Error()))

		return errPlayerCannotBeDeleted
	}

	return nil
}

// GetByID get a player with the given id.
func (s *Storage) GetByID(ctx context.Context, playerID players.PlayerID) (*players.Player, error) {
	s.logger.Debug("get player by id", slog.String("player id", playerID.String()))

	var player dbPlayer
	// id,firstname,lastname,nickname,usrpwd,country,date_created,date_updated
	err := s.db.QueryRowContext(ctx, selectByIDSQL, uuid.UUID(playerID)).
		Scan(
			&player.ID, &player.FirstName,
			&player.LastName, &player.Nickname,
			&player.Email, &player.Password,
			&player.Country, &player.DateCreated,
			&player.DateUpdated,
		)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("getting player by id",
			slog.String("player_id", playerID.String()),
			slog.String("error", err.Error()))

		return nil, errPlayerCannotBeRead
	}

	got := player.toPlayer()

	return &got, nil
}

// GetPlayersWithEmailOrNickName get players with given email or nickname.
func (s *Storage) GetPlayersWithEmailOrNickName(ctx context.Context, filter players.PlayerFilter) (*players.PlayerExistResult, error) {
	s.logger.Debug("get players by nickname or email", slog.Any("filter", filter))

	var countEmail, countNickName int

	var queryRow *sql.Row

	if filter.IgnoreID == nil {
		queryRow = s.db.QueryRowContext(ctx, selectByNicknameAndEmail,
			filter.Nickname,
			filter.Email,
		)
	} else {
		queryRow = s.db.QueryRowContext(ctx, selectByNicknameAndEmailAndID,
			filter.Nickname,
			filter.Email,
			filter.IgnoreID.String(),
		)
	}

	err := queryRow.Scan(&countNickName, &countEmail)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("getting players with given nickname or email",
			slog.Any("filter", filter),
			slog.String("error", err.Error()))

		return nil, errPlayersCannotBeRead
	}

	result := players.PlayerExistResult{
		EmailExist:    countEmail > 0,
		NicknameExist: countNickName > 0,
	}

	return &result, nil
}

// Search looks up players that match the given filter criteria.
func (s *Storage) Search(ctx context.Context, searchCriteria players.SearchCriteria) (*players.SearchResult, error) {
	s.logger.Debug("searching for players with search criteria", slog.Any("criteria", searchCriteria))

	result := players.SearchResult{
		Limit:  searchCriteria.Limit,
		Offset: searchCriteria.Offset,
	}

	searchFilters := buildSQLFilters(searchCriteria)

	count, err := s.queryCount(ctx, searchFilters)
	if err != nil {
		s.logger.Error("running count query of players that match given criteria",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.countStatement),
			slog.String("error", err.Error()),
		)

		return nil, errUnableToSearchPlayers
	}

	result.Total = count

	s.logger.Debug(
		"search players with filters",
		slog.String("query", searchFilters.query),
		slog.Any("filter", searchFilters),
	)

	playersFound, err := s.queryPlayers(ctx, searchFilters)
	if err != nil {
		s.logger.Error("checking if rows results has an error",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.query),
			slog.String("error", err.Error()),
		)

		return nil, errUnableToSearchPlayers
	}

	result.Items = toPlayerItems(playersFound)

	return &result, nil
}

func (s *Storage) queryCount(ctx context.Context, searchFilters *filterBuilder) (int, error) {
	var count int

	countStmt, err := s.db.Prepare(searchFilters.countStatement)
	if err != nil {
		s.logger.Error("building count players prepared statement",
			slog.Any("filter", searchFilters),
			slog.String("error", err.Error()),
		)

		return -1, fmt.Errorf("unable to build query to count players: %w", err)
	}

	defer countStmt.Close()

	row := countStmt.QueryRowContext(ctx, searchFilters.countArgs...)

	err = row.Scan(&count)
	if err != nil {
		s.logger.Error("scanning count of players found",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.countStatement),
			slog.String("error", err.Error()),
		)

		return -1, fmt.Errorf("unable to scanning count of players found: %w", err)
	}

	return count, nil
}

func (s *Storage) queryPlayers(ctx context.Context, searchFilters *filterBuilder) ([]dbPlayerItem, error) {
	rows, err := s.db.QueryContext(ctx, searchFilters.query, searchFilters.queryArgs...)
	if err != nil {
		s.logger.Error("running query to find players with given criteria",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.query),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("unable to query players: %w", err)
	}

	defer rows.Close()

	playersFound := make([]dbPlayerItem, 0)

	for rows.Next() {
		player := new(dbPlayerItem)
		// id, firstname, lastname, nickname, country
		rowErr := rows.Scan(&player.ID, &player.FirstName, &player.LastName, &player.Nickname, &player.Country)
		if rowErr != nil {
			s.logger.Error("scanning rows for searching players with search criteria",
				slog.Any("filter", searchFilters),
				slog.String("query", searchFilters.query),
				slog.String("error", rowErr.Error()),
			)

			return nil, fmt.Errorf("unable to scan player rows: %w", rowErr)
		}

		playersFound = append(playersFound, *player)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error("checking if rows results has an error",
			slog.Any("filter", searchFilters),
			slog.String("query", searchFilters.query),
			slog.String("error", err.Error()),
		)

		return nil, fmt.Errorf("player search query had some errors: %w", err)
	}

	return playersFound, nil
}

func (s *Storage) Health() (string, error) {
	err := s.db.Ping()
	if err != nil {
		s.logger.Error("doing ping to database", slog.String("error", err.Error()))

		return "storage", fmt.Errorf("unable to ping db: %w", err)
	}

	return "storage", nil
}

func buildSQLFilters(filters players.SearchCriteria) *filterBuilder {
	newFilterBuilder := &filterBuilder{
		filters:   make([]string, 0),
		countArgs: make([]interface{}, 0),
		queryArgs: make([]interface{}, 0),
	}

	if filters.Country != nil && *filters.Country != "" {
		newFilterBuilder.addCondition(countryColumn, equalsOperator, filters.Country)
	}

	var countWhereClause string
	for _, v := range newFilterBuilder.filters {
		countWhereClause += v
	}

	countStatement := fmt.Sprintf(countByFilterSQL, countWhereClause)
	newFilterBuilder.countStatement = countStatement

	newFilterBuilder.addFilter(" LIMIT", filters.Limit, true)
	newFilterBuilder.addFilter(" OFFSET", filters.Offset, true)

	var whereClause string
	for _, v := range newFilterBuilder.filters {
		whereClause += v
	}

	queryStatement := fmt.Sprintf(selectByFilterSQL, whereClause)
	newFilterBuilder.query = queryStatement

	return newFilterBuilder
}
