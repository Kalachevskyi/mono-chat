package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// NewUser constructor for User repository.
func NewUser(redisClient *redis.Client) *User {
	return &User{redisClient: redisClient}
}

// User represents User redis repository.
type User struct {
	redisClient *redis.Client
}

// Set - set user ID to bool value.
func (u User) Set(key string) error {
	if err := u.redisClient.Set(key, true, 0).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// CheckUser check if User exists.
func (u User) CheckUser(key string) (bool, error) {
	val, err := u.redisClient.Do("get", key).Bool()
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, errors.WithStack(err)
	}

	return val, nil
}
