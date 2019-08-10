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

// Package handlers is an interface adapters of application
package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// TokenUC - represents a usecase interface for processing "Token" business logic
type TokenUC interface {
	Set(chatID int64, token string) error
	Get(chatID int64) (string, error)
}

// NewToken - builds "NewToken" internal handler
func NewToken(tokenUC TokenUC, botWrapper *BotWrapper) *Token {
	return &Token{
		tokenUC:    tokenUC,
		BotWrapper: botWrapper,
	}
}

// Token - represents an internal handler for processing "Token"
type Token struct {
	tokenUC TokenUC
	*BotWrapper
}

// Handle - process the "Token", send the result to the user
func (t *Token) Handle(u tg.Update) {
	if err := t.tokenUC.Set(u.Message.Chat.ID, u.Message.CommandArguments()); err != nil {
		t.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}
	t.sendMSG(tg.NewMessage(u.Message.Chat.ID, "successfully set token"))
}
