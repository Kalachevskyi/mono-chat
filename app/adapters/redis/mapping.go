package redis

import (
	"encoding/json"

	"github.com/Kalachevskyi/mono-chat/app/model"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// NewMapping - builds mapping repository
func NewMapping(redisClient *redis.Client) *Mapping {
	return &Mapping{redisClient: redisClient}
}

// Mapping - represents an mapping repository for mapping app category with MonoBank category
type Mapping struct {
	redisClient *redis.Client
}

// Set - save category mapping for chat key in redis
func (t *Mapping) Set(key string, val map[string]model.CategoryMapping) error {
	mapping, err := json.Marshal(val)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := t.redisClient.Set(key, string(mapping), 0).Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Get - return category mapping for chat key from redis
func (t *Mapping) Get(key string) (map[string]model.CategoryMapping, error) {
	val, err := t.redisClient.Get(key).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var mapping map[string]model.CategoryMapping
	if err := json.Unmarshal([]byte(val), &mapping); err != nil {
		return nil, errors.WithStack(err)
	}

	return mapping, nil
}
