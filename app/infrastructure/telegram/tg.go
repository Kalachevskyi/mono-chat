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

// Package telegram is an data layer of application
package telegram

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// NewTelegram - builds Telegram repository
func NewTelegram(log Logger) *Telegram {
	return &Telegram{
		log: log,
	}
}

// Telegram - represents the Telegram repository for communication with Telegram api
type Telegram struct {
	log Logger
}

// GetFile -  get the file from Telegram REST API, makes HTTP call to telegram API
func (c *Telegram) GetFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url) //nolint
	if err != nil {
		return nil, errors.Errorf("can't get file by url: %s", url)
	}

	return resp.Body, nil
}
