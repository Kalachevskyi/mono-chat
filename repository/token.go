package repository

import "github.com/go-redis/redis"

func NewToken(redisClient *redis.Client) *Token {
	return &Token{redisClient: redisClient}
}

type Token struct {
	redisClient *redis.Client
}

func (t *Token) Set(key, token string) error {
	return t.redisClient.Set(key, token, 0).Err()
}

func (t *Token) Get(key string) (string, error) {
	return t.redisClient.Get(key).Result()
}
