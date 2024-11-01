package cryptos

import (
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type BcryptSetup struct {
	Cost   int
	Logger *slog.Logger
}

// Bcrypt defines logic to hash passwords using bcrypt library.
type Bcrypt struct {
	cost   int
	logger *slog.Logger
}

var errUnableToHashValue = errors.New("unable to hash password")

func NewBcrypt(bcryptSetup BcryptSetup) *Bcrypt {
	if bcryptSetup.Cost < bcrypt.MinCost {
		bcryptSetup.Cost = bcrypt.MinCost
	}

	newBcrypt := Bcrypt{
		cost:   bcryptSetup.Cost,
		logger: bcryptSetup.Logger,
	}

	return &newBcrypt
}

// Hash password using the bcrypt hashing algorithm.
func (b *Bcrypt) Hash(password string) ([]byte, error) {
	// Hash password with bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		b.logger.Error("generating hash from password", slog.Any("error", err.Error()))

		return nil, errUnableToHashValue
	}

	return hashedPasswordBytes, nil
}

// Check if two passwords match using Bcrypt's CompareHashAndPassword
// which return nil on success and an error on failure.
func (b *Bcrypt) DoPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))

	return err == nil
}
