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

// Package repository is an data layer of application
package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/Kalachevskyi/mono-chat/app/model"
)

const monoDomain = "https://api.monobank.ua"

// NewTransaction - builds Transaction repository
func NewTransaction(log Logger) *Transaction {
	return &Transaction{log: log}
}

// Transaction - represents the Transaction repository for getting transaction from MonoBank api
type Transaction struct {
	log Logger
}

// GetTransactions - return Transactions from MonoBank
func (m *Transaction) GetTransactions(token string, from, to time.Time) ([]model.Transaction, error) {
	url := fmt.Sprintf("%s/personal/statement/0/%d/%d", monoDomain, from.Unix(), to.Unix())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set("X-Token", token)

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
