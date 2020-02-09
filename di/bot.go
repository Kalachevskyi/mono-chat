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

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	tgBot  *tg.BotAPI //nolint:gochecknoglobals
	tgOnce sync.Once  //nolint:gochecknoglobals
)

// TGBot - initialize the instance of Telegram client
func TGBot(token string) (*tg.BotAPI, error) {
	var err error

	tgOnce.Do(func() {
		tgBot, err = tg.NewBotAPI(token)
		if err != nil {
			err = fmt.Errorf("can't initialize Telegram client: err=%s", err.Error())
			return
		}
	})

	return tgBot, err
}
