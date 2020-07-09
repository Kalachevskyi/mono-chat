package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// NewToken - builds token repository
func NewToken(redisClient *redis.Client) *Token {
	return &Token{redisClient: redisClient}
}

// Token - represents the Token repository to interact with the token
type Token struct {
	redisClient *redis.Client
}

// Set - save the chat session token in redis
func (t *Token) Set(key, token string) error {
	if err := t.redisClient.Set(key, token, 0).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Get - return chat session token from redis
func (t *Token) Get(key string) (string, error) {
	val, err := t.redisClient.Get(key).Result()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return val, nil
}
