// Copyright © 2019 Volodymyr Kalachevskyi <v.kalachevskyi@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package repository is an data layer of application
package repository

import (
	"encoding/json"

	"github.com/Kalachevskyi/mono-chat/entities"
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

// Get - return category mapping for chat key from redis
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
