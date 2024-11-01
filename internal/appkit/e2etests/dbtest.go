package e2etests

import (
	"flag"
	"fmt"
	"math/rand"
	"net/mail"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/fernandoocampo/players/internal/adapters/storages"
	"github.com/fernandoocampo/players/internal/players"
	"github.com/google/uuid"
)

var E2ETest = flag.Bool("e2e-test", false, "this flags indicates if this is an e2e test")

func DatabaseParametersFixture(t *testing.T) storages.Parameters {
	t.Helper()

	dbParameters := storages.Parameters{
		DBName:   getStringEnvVar(t, "PLAYERS_POSTGRES_DB", "playersdb"),
		Host:     getStringEnvVar(t, "PLAYERS_POSTGRES_HOST", "localhost"),
		Player:   getStringEnvVar(t, "PLAYERS_POSTGRES_PLAYER", "playersdb"),
		Password: getStringEnvVar(t, "PLAYERS_POSTGRES_PASSWORD", ""),
		Port:     getIntEnvVar(t, "PLAYERS_POSTGRES_PORT", 5432),
	}

	return dbParameters
}

func getStringEnvVar(t *testing.T, key, defaultValue string) string {
	t.Helper()

	value := os.Getenv(key)
	if value == "" {
		t.Logf("env var %s is empty, using default %q", key, defaultValue)

		return defaultValue
	}

	return value
}

func getIntEnvVar(t *testing.T, key string, defaultValue int) int {
	t.Helper()

	value := os.Getenv(key)
	if value == "" {
		t.Logf("env var %s is empty, using default %q", key, defaultValue)

		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		t.Logf("env var %s has an invalid value, using default %q", key, defaultValue)

		return defaultValue
	}

	return intValue
}

func RandomPlayerFixture() players.Player {
	newPlayerID := players.PlayerID(uuid.New())
	newNickname := randomString(10)
	newEmail := mail.Address{
		Address: fmt.Sprintf("%s@e2etest.com", newNickname),
	}
	anyPassword := []byte("$2a$04$zedC.nDTul7ks4kELCsb4OldjunQoDkeMisEk822pY6XqtAJo1uo6")
	anyUTCDate := time.Now().UTC()

	return players.Player{
		ID:          &newPlayerID,
		FirstName:   "e2e",
		LastName:    "dbtest",
		Nickname:    newNickname,
		Email:       newEmail,
		Password:    anyPassword,
		Country:     randomCountry(),
		DateCreated: anyUTCDate,
		DateUpdated: anyUTCDate,
	}
}

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func randomCountry() string {
	countries := []string{"UK", "Spain", "Colombia", "France", "Germany"}

	return countries[rand.Intn(4)]
}
