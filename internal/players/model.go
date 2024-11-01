package players

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// Storage defines behaviour for player repositories.
type Storage interface {
	// Save persists a new player in the player repository.
	Save(ctx context.Context, player Player) error
	// Update player in the player repository.
	Update(ctx context.Context, player Player) error
	// Delete player in the repository.
	Delete(ctx context.Context, playerID PlayerID) error
	// GetByID get a player with the given id.
	GetByID(ctx context.Context, id PlayerID) (*Player, error)
	// GetPlayersWithEmailOrNickName get players with given email or nickname.
	GetPlayersWithEmailOrNickName(ctx context.Context, filter PlayerFilter) (*PlayerExistResult, error)
	// Search looks up players that match the given filter criteria.
	Search(ctx context.Context, searchCriteria SearchCriteria) (*SearchResult, error)
}

// Hasher defines behaviour for crypto mechanisms.
type Hasher interface {
	Hash(password string) ([]byte, error)
}

// Notifier defines behavior to notify about player events.
type Notifier interface {
	// Notify notifies a new event.
	Notify(event NewEvent)
}

// PlayerID defines player id.
type PlayerID uuid.UUID

// Player contains player data.
type Player struct {
	ID          *PlayerID
	FirstName   string
	LastName    string
	Nickname    string
	Email       mail.Address
	Password    []byte
	Country     string
	DateCreated time.Time
	DateUpdated time.Time
}

// NewPlayer contains data required to create a new player.
type NewPlayer struct {
	FirstName string
	LastName  string
	Nickname  string
	Email     mail.Address
	Password  string
	Country   string
}

// UpdatePlayer contains data required to update player.
type UpdatePlayer struct {
	ID        PlayerID
	FirstName *string
	LastName  *string
	Nickname  *string
	Email     *mail.Address
	Password  *string
	Country   *string
}

// PlayerFilter email or nicknames fields.
type PlayerFilter struct {
	Email    string
	Nickname string
	IgnoreID *PlayerID
}

// PlayerExistResult indicates whether a player already exists with the given email or nickname.
type PlayerExistResult struct {
	EmailExist    bool
	NicknameExist bool
}

// PlayerAlreadyExistsError defines error for trying to create a player
// that already exists with the given email or nickname.
type PlayerAlreadyExistsError struct {
	wasEmail    bool
	wasNickname bool
}

// SearchCriteria criteria data to search players.
type SearchCriteria struct {
	Country *string
	// determines the number of rows.
	Limit uint16 `json:"limit"`
	// skips the offset rows before beginning to return the rows.
	Offset uint16 `json:"offset"`
}

// PlayerItem contains few data about a player.
type PlayerItem struct {
	ID        PlayerID
	FirstName string
	LastName  string
	Nickname  string
	Country   string
}

// SearchResult contains data about the result of a specific player search.
type SearchResult struct {
	Items []PlayerItem `json:"items"`
	// Total total number of players that match the search criteria.
	Total int `json:"total"`
	// determines the number of rows.
	Limit uint16 `json:"limit"`
	// skips the offset rows before beginning to return the rows.
	Offset uint16 `json:"offset"`
}

// CreatePlayerResult standard response for creating a Player.
type CreatePlayerResult struct {
	ID  *PlayerID
	Err string
}

// UpdatePlayerResult standard response for updating a Player.
type UpdatePlayerResult struct {
	Err string
}

// DeletePlayerResult standard response for deleting a Player.
type DeletePlayerResult struct {
	Err string
}

// SearchPlayersDataResult standard response for search players.
type SearchPlayersDataResult struct {
	SearchResult *SearchResult
	Err          string
}

// NewEvent contains data about player events.
type NewEvent struct {
	PlayerID string
	Event    string
}

// updateToPlayerResult result after converting UpdatePlayer to Player.
type updateToPlayerResult struct {
	player  *Player
	changes bool
}

const (
	redactedValue = "[REDACTED]"
)

var (
	ErrInvalidPlayerID    = errors.New("invalid player id")
	ErrPlayerDoesNotExist = errors.New("player doesn't exist")
	errEmptyFirstName     = errors.New("first name is empty")
	errEmptyLastName      = errors.New("last name is empty")
	errEmptyNickname      = errors.New("nickname is empty")
	errEmptyCountry       = errors.New("country is empty")
	errEmptyPassword      = errors.New("password is empty")
	errEmptyEmail         = errors.New("email is empty")
)

func (u Player) obfuscate() Player {
	obfuscated := u
	obfuscated.Password = []byte(redactedValue)

	return obfuscated
}

func (u *PlayerAlreadyExistsError) Error() string {
	if u.wasEmail && u.wasNickname {
		return "player with the given email or nickname already exists"
	}

	if u.wasEmail {
		return "player with the given email already exists"
	}

	return "player with the given nickname already exists"
}

func (u PlayerExistResult) toPlayerAlreadyExistsError() *PlayerAlreadyExistsError {
	newError := PlayerAlreadyExistsError{
		wasEmail:    u.EmailExist,
		wasNickname: u.NicknameExist,
	}

	return &newError
}

func (u PlayerExistResult) Exist() bool {
	return u.EmailExist || u.NicknameExist
}

// newPlayerID creates a new player id value.
func newPlayerID() PlayerID {
	return PlayerID(uuid.New())
}

func (u PlayerID) String() string {
	return uuid.UUID(u).String()
}

func (n NewPlayer) toPlayer(hashedPassword []byte) Player {
	newPlayerID := newPlayerID()

	now := time.Now().UTC()

	return Player{
		ID:          &newPlayerID,
		FirstName:   n.FirstName,
		LastName:    n.LastName,
		Nickname:    n.Nickname,
		Email:       n.Email,
		Password:    hashedPassword,
		Country:     n.Country,
		DateCreated: now,
		DateUpdated: now,
	}
}

func (n NewPlayer) Validate() error {
	var err error

	if n.FirstName == "" {
		err = errors.Join(err, errEmptyFirstName)
	}

	if n.LastName == "" {
		err = errors.Join(err, errEmptyLastName)
	}

	if n.Nickname == "" {
		err = errors.Join(err, errEmptyNickname)
	}

	if n.Email.Address == "" {
		err = errors.Join(err, errEmptyEmail)
	}

	if n.Country == "" {
		err = errors.Join(err, errEmptyCountry)
	}

	if n.Password == "" {
		err = errors.Join(err, errEmptyPassword)
	}

	return err
}

func (n NewPlayer) obfuscate() NewPlayer {
	obfuscated := n
	obfuscated.Password = redactedValue

	return obfuscated
}

func (n NewPlayer) toPlayerFilter() PlayerFilter {
	return PlayerFilter{
		Email:    n.Email.Address,
		Nickname: n.Nickname,
		IgnoreID: nil,
	}
}

func (u UpdatePlayer) Validate() error {
	var err error

	if isStringEmpty(u.FirstName) {
		err = errors.Join(err, errEmptyFirstName)
	}

	if isStringEmpty(u.LastName) {
		err = errors.Join(err, errEmptyLastName)
	}

	if isStringEmpty(u.Nickname) {
		err = errors.Join(err, errEmptyNickname)
	}

	if u.Email != nil && u.Email.Address == "" {
		err = errors.Join(err, errEmptyEmail)
	}

	if isStringEmpty(u.Country) {
		err = errors.Join(err, errEmptyCountry)
	}

	if isStringEmpty(u.Password) {
		err = errors.Join(err, errEmptyPassword)
	}

	return err
}

func (u UpdatePlayer) obfuscate() UpdatePlayer {
	obfuscated := u
	passwordObfuscated := redactedValue
	obfuscated.Password = &passwordObfuscated

	return obfuscated
}

func (u UpdatePlayer) toPlayerFilter() PlayerFilter {
	var newPlayerFilter PlayerFilter

	newPlayerFilter.IgnoreID = &u.ID

	if u.Email != nil {
		newPlayerFilter.Email = u.Email.Address
	}

	if u.Nickname != nil {
		newPlayerFilter.Nickname = *u.Nickname
	}

	return newPlayerFilter
}

// toPlayer updates the data for the given player with the data in updateplayer.
func (u UpdatePlayer) toPlayer(player Player, hasher Hasher) (*updateToPlayerResult, error) {
	var result updateToPlayerResult
	if areStringDifferent(u.FirstName, player.FirstName) {
		result.changes = true
		player.FirstName = *u.FirstName
	}

	if areStringDifferent(u.LastName, player.LastName) {
		result.changes = true
		player.LastName = *u.LastName
	}

	if areStringDifferent(u.Nickname, player.Nickname) {
		result.changes = true
		player.Nickname = *u.Nickname
	}

	if areStringDifferent(u.Country, player.Country) {
		result.changes = true
		player.Country = *u.Country
	}

	if u.Email != nil && u.Email.Address != player.Email.Address {
		result.changes = true
		player.Email = *u.Email
	}

	if u.Password != nil {
		result.changes = true

		hashedPassword, err := hasher.Hash(*u.Password)
		if err != nil {
			return nil, fmt.Errorf("unable to hash password: %w", err)
		}

		player.Password = hashedPassword
	}

	if result.changes {
		player.DateUpdated = time.Now().UTC()
		result.player = &player
	}

	return &result, nil
}

func (u UpdatePlayer) updateKeyValues() bool {
	return u.Email != nil || u.Nickname != nil
}

func (s SearchCriteria) isEmpty() bool {
	empty := true

	if !isStringEmpty(s.Country) {
		empty = false
	}

	return empty
}

func (s *SearchCriteria) setDefaultPaginationIfEmpty() {
	if s.Limit == 0 {
		s.Limit = 5
	}
}

func newEmptySearchResult() *SearchResult {
	newSearchResult := SearchResult{
		Items:  make([]PlayerItem, 0),
		Total:  0,
		Limit:  0,
		Offset: 0,
	}

	return &newSearchResult
}

// newCreatePlayerResult create a new CreatePlayerResponse.
func newCreatePlayerResult(playerID *PlayerID, err error) CreatePlayerResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}

	return CreatePlayerResult{
		ID:  playerID,
		Err: errmessage,
	}
}

// newUpdatePlayerResult udpate a new UpdatePlayerResponse.
func newUpdatePlayerResult(err error) UpdatePlayerResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}

	return UpdatePlayerResult{
		Err: errmessage,
	}
}

// newDeletePlayerResult udpate a new DeletePlayerResponse.
func newDeletePlayerResult(err error) DeletePlayerResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}

	return DeletePlayerResult{
		Err: errmessage,
	}
}

// newSearchPlayersResult create a new SearchPlayersResult.
func newSearchPlayersDataResult(result *SearchResult, err error) SearchPlayersDataResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}

	return SearchPlayersDataResult{
		SearchResult: result,
		Err:          errmessage,
	}
}

func areStringDifferent(newValue *string, oldValue string) bool {
	if newValue == nil || *newValue == "" {
		return false
	}

	return *newValue != oldValue
}

func isStringEmpty(value *string) bool {
	return value != nil && *value == ""
}

func NewString(value string) *string {
	return &value
}

func StringToPlayerID(playerID string) (*PlayerID, error) {
	uuidValue, err := uuid.Parse(playerID)
	if err != nil {
		return nil, ErrInvalidPlayerID
	}

	return ToPlayerID(uuidValue), nil
}

func ToPlayerID(uuidValue uuid.UUID) *PlayerID {
	newPlayerID := PlayerID(uuidValue)

	return &newPlayerID
}

func FromString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func newCreatePlayerEvent(playerID *PlayerID) NewEvent {
	return NewEvent{
		PlayerID: playerID.String(),
		Event:    "new player was created",
	}
}

func newUpdatePlayerEvent(playerID *PlayerID) NewEvent {
	return NewEvent{
		PlayerID: playerID.String(),
		Event:    "player was updated",
	}
}

func newDeletePlayerEvent(playerID PlayerID) NewEvent {
	return NewEvent{
		PlayerID: playerID.String(),
		Event:    "player was deleted",
	}
}
