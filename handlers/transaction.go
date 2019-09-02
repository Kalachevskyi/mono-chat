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
	"io"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"
)

// ApiUC - represents a use-case interface for processing business logic of "MonoBank" transactions API
type ApiUC interface {
	GetTransactions(token string, chatID int64, from time.Time, to time.Time) (io.Reader, error)
	ParseDate(period string) (from time.Time, to time.Time, err error)
	Locale() *time.Location
}

// NewTransaction - builds "NewTransaction" internal handler
func NewTransaction(tokenUC TokenUC, apiUC ApiUC, botWrapper *BotWrapper) *Transaction {
	return &Transaction{tokenUC: tokenUC, apiUC: apiUC, BotWrapper: botWrapper}
}

// Transaction - represents an internal handler for processing "MonoBank" transactions API
type Transaction struct {
	tokenUC TokenUC
	apiUC   ApiUC
	*BotWrapper
}

// Handle - process the "MonoBank" transactions API, send the result to the user
func (t *Transaction) Handle(u tg.Update) {
	var (
		from, to time.Time
		timeNow  = now.New(time.Now().In(t.apiUC.Locale()))
	)
	switch u.Message.Command() {
	case getCommand:
		fromTime, toTime, err := t.apiUC.ParseDate(u.Message.CommandArguments())
		if err != nil {
			t.sendDefaultErr(u.Message.Chat.ID, err)
			return
		}
		from, to = fromTime, toTime
	case todayCommand:
		from, to = timeNow.BeginningOfDay(), timeNow.EndOfDay()
	case currentMonthCommand:
		from, to = timeNow.BeginningOfMonth(), timeNow.EndOfMonth()
	default:
		t.sendDefaultErr(u.Message.Chat.ID, errors.New("can't detect command"))
		return
	}

	chatID := u.Message.Chat.ID
	token, err := t.tokenUC.Get(chatID)
	if err != nil {
		t.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileResp, err := t.apiUC.GetTransactions(token, chatID, from, to)
	if err != nil {
		t.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}
	reader := tg.FileReader{
		Name:   fmt.Sprintf("%s-%s%s", from.Format(dateTimePattern), to.Format(dateTimePattern), ".csv"),
		Reader: fileResp,
		Size:   -1,
	}
	msg := tg.NewDocumentUpload(u.Message.Chat.ID, reader)
	t.sendMSG(msg)
}
