// Package redis is an data layer of application
package redis

import (
	"github.com/Kalachevskyi/mono-chat/app/model"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// NewAccount - create user account constructor
func NewAccount(redisClient *redis.Client) *Account {
	return &Account{redisClient: redisClient}
}

// Account - represents user Account
type Account struct {
	redisClient *redis.Client
}

// Set - save the chosen account in redis
func (a *Account) Set(key, account string) error {
	if err := a.redisClient.Set(key, account, 0).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Get - get chosen account from redis
func (a Account) Get(key string) (string, error) {
	val, err := a.redisClient.Get(key).Result()
	if err != nil && err.Error() != string(redis.Nil) {
		return "", errors.WithStack(err)
	}

	if err != nil && err.Error() == string(redis.Nil) {
		return "", model.ErrNil
	}

	return val, nil
}
