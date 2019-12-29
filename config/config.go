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

// Package config implements the application configuration
package config

import (
	"github.com/pkg/errors"
)

// Config - app configuration
type Config struct {
	Token       string // Telegram token
	Debug       bool   // Debug mod
	Offset      int
	Timeout     int
	EncodingLog string // Valid values are "json" and "console",
	RedisURL    string // Example localhost:6379
}

// Validate - verify app configuration
func (c *Config) Validate() error {
	if c.Token == "" {
		return errors.New(`config parameter "token" can't be empty`)
	}

	if c.Timeout == 0 {
		return errors.New(`config parameter "timeout" can't be empty`)
	}

	if c.RedisURL == "" {
		return errors.New(`config parameter "redis_url" can't be empty`)
	}

	return nil
}
