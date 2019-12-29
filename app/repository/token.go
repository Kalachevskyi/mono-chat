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
