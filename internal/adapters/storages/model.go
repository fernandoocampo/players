package storages

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/fernandoocampo/players/internal/players"
	"github.com/google/uuid"
)

type dbPlayer struct {
	ID          uuid.UUID `db:"id"`
	FirstName   string    `db:"firstname"`
	LastName    string    `db:"lastname"`
	Nickname    string    `db:"nickname"`
	Email       string    `db:"email"`
	Password    []byte    `db:"usrpwd"`
	Country     string    `db:"country"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

type dbPlayerItem struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"firstname"`
	LastName  string    `db:"lastname"`
	Nickname  string    `db:"nickname"`
	Country   string    `db:"country"`
}

type filterBuilder struct {
	query          string
	countStatement string
	filters        []string
	queryArgs      []interface{}
	countArgs      []interface{}
}

// player columns.
const (
	countryColumn = "country"
)

const (
	equalsOperator = "="
	bsonInOperator = "@>"
	whereOperator  = "WHERE"
	andOperator    = "AND"
)

func (d *dbPlayer) toPlayer() players.Player {
	return players.Player{
		ID:          players.ToPlayerID(d.ID),
		FirstName:   d.FirstName,
		LastName:    d.LastName,
		Nickname:    d.Nickname,
		Email:       mail.Address{Address: d.Email},
		Password:    d.Password,
		Country:     d.Country,
		DateCreated: d.DateCreated.UTC(),
		DateUpdated: d.DateUpdated.UTC(),
	}
}

func (d dbPlayerItem) toPlayerItem() players.PlayerItem {
	return players.PlayerItem{
		ID:        *players.ToPlayerID(d.ID),
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Nickname:  d.Nickname,
		Country:   d.Country,
	}
}

func (f *filterBuilder) addCondition(field, operator string, value interface{}) *filterBuilder {
	isHint := false
	condition := whereOperator

	if len(f.filters) > 0 {
		condition = " " + andOperator
	}

	newStatement := fmt.Sprintf("%s %s %s", condition, field, operator)

	return f.addFilter(newStatement, value, isHint)
}

func (f *filterBuilder) addFilter(statement string, value interface{}, isHint bool) *filterBuilder {
	index := len(f.filters) + 1
	statement = fmt.Sprintf("%s $%d", statement, index)
	f.filters = append(f.filters, statement)

	if !isHint {
		f.countArgs = append(f.countArgs, value)
	}

	f.queryArgs = append(f.queryArgs, value)

	return f
}

func toDBPlayer(player *players.Player) dbPlayer {
	return dbPlayer{
		ID:          uuid.UUID(*player.ID),
		FirstName:   player.FirstName,
		LastName:    player.LastName,
		Nickname:    player.Nickname,
		Email:       player.Email.Address,
		Password:    player.Password,
		Country:     player.Country,
		DateCreated: player.DateCreated,
		DateUpdated: player.DateUpdated,
	}
}

func toPlayerItems(dbPlayerItems []dbPlayerItem) []players.PlayerItem {
	result := make([]players.PlayerItem, 0, len(dbPlayerItems))

	for _, v := range dbPlayerItems {
		result = append(result, v.toPlayerItem())
	}

	return result
}
