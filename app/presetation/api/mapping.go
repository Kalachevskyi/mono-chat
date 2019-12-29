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

// Package api is an interface adapters of application
package api

import (
	"io"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// MappingUC - represents a usecase interface for processing category mapping business logic
type MappingUC interface {
	Validate(name string) error
	Parse(chatID int64, r io.Reader) error
	GetFile(url string) (io.ReadCloser, error)
}

// NewMapping - builds "NewMapping" internal handler
func NewMapping(uc MappingUC, botWrapper *BotWrapper) *Mapping {
	return &Mapping{uc: uc, BotWrapper: botWrapper}
}

// Mapping - represents an internal handler for processing category mapping
type Mapping struct {
	uc MappingUC
	*BotWrapper
}

// Handle - process category mapping, send the result to the user
func (m *Mapping) Handle(u tg.Update) {
	if err := m.uc.Validate(u.Message.Document.FileName); err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileTG, err := m.bot.GetFile(tg.FileConfig{FileID: u.Message.Document.FileID})
	if err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	file, err := m.uc.GetFile(fileTG.Link(m.bot.Token))
	if err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	if err := m.uc.Parse(u.Message.Chat.ID, file); err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}
	m.sendMSG(tg.NewMessage(u.Message.Chat.ID, "mapping successfully loaded"))
}
