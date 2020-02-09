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

// Package di is an application layer for initializing app components
package di

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client //nolint:gochecknoglobals
	redisOnce   sync.Once     //nolint:gochecknoglobals
)

// RedisClient - initialize the instance of redis client
func RedisClient(url string) (*redis.Client, error) {
	var err error

	redisOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     url,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		if _, err = client.Ping().Result(); err != nil {
			err = fmt.Errorf("can't initialize Redis client: err=%s", err.Error())
			return
		}
		redisClient = client
	})

	return redisClient, err
}
