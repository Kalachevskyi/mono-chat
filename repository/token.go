package repository

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func NewToken(redisClient *redis.Client) *Token {
	return &Token{redisClient: redisClient}
}

type Token struct {
	redisClient *redis.Client
}

func (t *Token) Set(key, token string) error {
	if err := t.redisClient.Set(key, token, 0).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (t *Token) Get(key string) (string, error) {
	val, err := t.redisClient.Get(key).Result()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return val, nil
}
