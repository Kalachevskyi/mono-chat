package usecases

import "fmt"

type TokenRepo interface {
	Set(key, token string) error
	Get(key string) (string, error)
}

func NewToken(repo TokenRepo) *Token {
	return &Token{repo: repo}
}

type Token struct {
	repo TokenRepo
}

func (c *Token) Set(chatID int64, token string) error {
	key := fmt.Sprintf("token_%v", chatID)
	return c.repo.Set(key, token)
}

func (c *Token) Get(chatID int64) (string, error) {
	key := fmt.Sprintf("token_%v", chatID)
	return c.repo.Get(key)
}
