package repository

import "github.com/go-redis/redis"

func NewToken(redisClient *redis.Client) *Token {
	return &Token{redisClient: redisClient}
}

type Token struct {
	redisClient *redis.Client
}

func (t *Token) SaveToken(key, token string) error {
	if err := t.redisClient.Set(key, token, 0).Err(); err != nil {
		return err
	}
	return nil
}
