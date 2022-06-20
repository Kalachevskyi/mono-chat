package usecases

import (
	"fmt"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=./token_mock_test.go -package=usecases_test -source=./token.go

// TokenRepo - represents Token repository interface.
type TokenRepo interface {
	Set(key, token string) error
	Get(key string) (string, error)
}

// NewToken - builds Token report use-case.
func NewToken(repo TokenRepo) *Token {
	return &Token{repo: repo}
}

// Token - represents Token use-case for processing token.
type Token struct {
	repo TokenRepo
}

// Set - save token by key.
func (c *Token) Set(userID uuid.UUID, token string) error {
	key := fmt.Sprintf("token_%v", userID)

	return c.repo.Set(key, token)
}

// Get - return token by key.
func (c *Token) Get(userID uuid.UUID) (string, error) {
	key := fmt.Sprintf("token_%v", userID)

	return c.repo.Get(key)
}
