package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	"github.com/Kalachevskyi/mono-chat/app/model"
)

// NewGeneric - create generic constructor.
func NewGeneric(redisClient *redis.Client) *Generic {
	return &Generic{redisClient: redisClient}
}

// Generic - represents Generic redis repository.
type Generic struct {
	redisClient *redis.Client
}

// Set - save key, val in redis.
func (a *Generic) Set(key, val string) error {
	if err := a.redisClient.Set(key, val, 0).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Get - get value from redis.
func (a Generic) Get(key string) (string, error) {
	val, err := a.redisClient.Get(key).Result()
	if err != nil && err.Error() != string(redis.Nil) {
		return "", errors.WithStack(err)
	}

	if err != nil && err.Error() == string(redis.Nil) {
		return "", model.ErrNil
	}

	return val, nil
}
