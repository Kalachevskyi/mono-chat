package usecases

import "fmt"

type TokenRepo interface {
	SaveToken(key, token string) error
}

func NewToken(repo TokenRepo) *Token {
	return &Token{repo: repo}
}

type Token struct {
	repo TokenRepo
}

func (c *Token) SaveToken(chatID int64, token string) error {
	key := fmt.Sprintf("token_%v", chatID)
	return c.repo.SaveToken(key, token)
}
