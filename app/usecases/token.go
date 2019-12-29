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

// Package usecases is the business logic layer of the application.
package usecases

import "fmt"

//go:generate mockgen -destination=./token_mock_test.go -package=usecases_test -source=./token.go

// TokenRepo - represents Token repository interface
type TokenRepo interface {
	Set(key, token string) error
	Get(key string) (string, error)
}

// NewToken - builds Token report use-case
func NewToken(repo TokenRepo) *Token {
	return &Token{repo: repo}
}

// Token - represents Token use-case for processing token
type Token struct {
	repo TokenRepo
}

// Set - save token by key
func (c *Token) Set(chatID int64, token string) error {
	key := fmt.Sprintf("token_%v", chatID)
	return c.repo.Set(key, token)
}

// Get - return token by key
func (c *Token) Get(chatID int64) (string, error) {
	key := fmt.Sprintf("token_%v", chatID)
	return c.repo.Get(key)
}
