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
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Error messages
const defaultErrMSG = "Sorry, I can't process this message, view the logs or contact the owner of the service."

const dateTimePattern = "02.01.2006T15.04"

// Command messages
const (
	getCommand          = "get"
	todayCommand        = "today"
	currentMonthCommand = "month"
	tokenCommand        = "token"
)

// Logger - represents the application's logger interface
type Logger interface {
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}

// Handler - represents internal handler interface
type Handler interface {
	Handle(u tg.Update)
}

// NewChat - builds main chat handler
func NewChat(updates tg.UpdatesChannel, handlers map[HandlerKey]Handler, botWrapper *BotWrapper) *Chat {
	return &Chat{
		updates:    updates,
		handlers:   handlers,
		BotWrapper: botWrapper,
	}
}

// Chat - main chat handler
type Chat struct {
	updates  tg.UpdatesChannel
	handlers map[HandlerKey]Handler
	*BotWrapper
}

// Handle - routes between internal handlers depending on the type of message
func (c *Chat) Handle() {
	fmt.Println()
	for u := range c.updates {
		if u.Message == nil { // ignore any non-Message Updates
			continue
		}

		if u.Message.Document != nil {
			switch u.Message.Document.FileName {
			case "mapping.csv":
				if h, ok := c.handlers[MappingHandler]; ok {
					h.Handle(u)
					continue
				}
			default:
				if h, ok := c.handlers[FileReportHandler]; ok {
					h.Handle(u)
					continue
				}
			}
			continue
		}

		if u.Message.Command() != "" {
			switch u.Message.Command() {
			case getCommand, todayCommand, currentMonthCommand:
				if h, ok := c.handlers[TransactionsHandler]; ok {
					h.Handle(u)
					continue
				}
			case tokenCommand:
				if h, ok := c.handlers[TokenHandler]; ok {
					h.Handle(u)
					continue
				}
				return
			}
		}

		c.sendMSG(tg.NewMessage(u.Message.Chat.ID, defaultErrMSG))
	}
}
