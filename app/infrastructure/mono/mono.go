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

// Package mono is an data layer of application
package mono

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kalachevskyi/mono-chat/app/infrastructure/telegram"

	"github.com/pkg/errors"

	"github.com/Kalachevskyi/mono-chat/app/model"
)

const (
	domainMono   = "https://api.monobank.ua"
	tokenMonoKey = "X-Token"
)

// NewMono - builds Mono repository
func NewMono(log telegram.Logger) *Mono {
	return &Mono{log: log}
}

// Mono - represents the Mono repository for getting transaction from MonoBank api
type Mono struct {
	log Logger
}

// GetTransactions - return Transactions from MonoBank
func (m *Mono) GetTransactions(token, account string, from, to time.Time) ([]model.Transaction, error) {
	url := fmt.Sprintf("%s/personal/statement/%s/%d/%d", domainMono, account, from.Unix(), to.Unix())

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set(tokenMonoKey, token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer closeBody(resp.Body, m.log)

	transactions := make([]model.Transaction, 0)
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, errors.WithStack(err)
	}

	return transactions, nil
}

// GetClientInfo - returns information about accounts (card, currency)
func (m Mono) GetClientInfo(token string) (c model.ClientInfo, err error) {
	url := fmt.Sprintf("%s/personal/client-info", domainMono)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return c, errors.WithStack(err)
	}

	req.Header.Set(tokenMonoKey, token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c, errors.WithStack(err)
	}

	defer closeBody(resp.Body, m.log)

	clientInfo := model.ClientInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&clientInfo); err != nil {
		return c, errors.WithStack(err)
	}

	return clientInfo, nil
}
