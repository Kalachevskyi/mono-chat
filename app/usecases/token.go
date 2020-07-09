package usecases

import "fmt"

//go:generate mockgen -destination=./token_mock_test.go -package=usecases_test -source=./token.go

// TokenRepo - represents Token repository interface
type TokenRepo interface {
	Set(key, token string) error
	Get(key string) (string, error)
}

// NewToken - builds Token report use-case
func NewToken(repo TokenRepo) *Token {
	return &Token{repo: repo}
}

// Token - represents Token use-case for processing token
type Token struct {
	repo TokenRepo
}

// Set - save token by key
func (c *Token) Set(chatID int64, token string) error {
	key := fmt.Sprintf("token_%v", chatID)
	return c.repo.Set(key, token)
}

// Get - return token by key
func (c *Token) Get(chatID int64) (string, error) {
	key := fmt.Sprintf("token_%v", chatID)
	return c.repo.Get(key)
}
