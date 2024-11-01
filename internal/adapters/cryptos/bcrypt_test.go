package cryptos_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/fernandoocampo/players/internal/adapters/cryptos"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	// Given
	bcryptSettings := cryptos.BcryptSetup{
		Cost:   bcrypt.MinCost,
		Logger: newLogger(),
	}
	aBcrypt := cryptos.NewBcrypt(bcryptSettings)
	rawPassword := "password1"

	// When
	hashedPassword, err := aBcrypt.Hash(rawPassword)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.True(t, aBcrypt.DoPasswordsMatch(string(hashedPassword), rawPassword))
}

func newLogger() *slog.Logger {
	handlerOptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	loggerHandler := slog.NewJSONHandler(os.Stdout, handlerOptions)
	return slog.New(loggerHandler)
}
