package repository

import (
	"encoding/json"

	"github.com/Kalachevskyi/mono-chat/entities"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

func NewMapping(redisClient *redis.Client) *Mapping {
	return &Mapping{redisClient: redisClient}
}

type Mapping struct {
	redisClient *redis.Client
}

func (t *Mapping) Set(key string, val map[string]entities.CategoryMapping) error {
	mapping, err := json.Marshal(val)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := t.redisClient.Set(key, string(mapping), 0).Err(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (t *Mapping) Get(key string) (map[string]entities.CategoryMapping, error) {
	val, err := t.redisClient.Get(key).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var mapping map[string]entities.CategoryMapping
	if err := json.Unmarshal([]byte(val), &mapping); err != nil {
		return nil, errors.WithStack(err)
	}
	return mapping, nil
}
