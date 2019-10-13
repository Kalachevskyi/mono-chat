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

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Kalachevskyi/mono-chat/entities"
	"github.com/pkg/errors"
)

//go:generate mockgen -destination=./transaction_mock_test.go -package=usecases -source=./transaction.go

// Logger - represents the application's logger interface
type Logger interface {
	Error(args ...interface{})
}

// TransactionRepo - represents Transaction repository interface
type TransactionRepo interface {
	GetTransactions(token string, from time.Time, to time.Time) ([]entities.Transaction, error)
}

type categoryMapping map[string]entities.CategoryMapping

// NewTransaction - builds Transaction report use-case
func NewTransaction(trRepo TransactionRepo, mapRepo MappingRepo, log Logger, date Date) *Transaction {
	return &Transaction{
		apiRepo:     trRepo,
		mappingRepo: mapRepo,
		log:         log,
		Date:        date,
	}
}

// Transaction - represents Transaction  use-case for processing bank Transactions
type Transaction struct {
	apiRepo     TransactionRepo
	mappingRepo MappingRepo
	log         Logger
	Date
}

// GetTransactions - get bank transactions, convert it to app csv report
func (a *Transaction) GetTransactions(token string, chatID int64, from time.Time, to time.Time) (io.Reader, error) {
	transactions, err := a.apiRepo.GetTransactions(token, from, to)
	if err != nil {
		return nil, err
	}

	catMap := a.getCategoryMapping(chatID)
	records := [][]string{
		{
			DateHeader.Str(),
			DescriptionHeader.Str(),
			CategoryHeader.Str(),
			BankCategoryHeader.Str(),
			AmountHeader.Str(),
		},
	}

	for _, tr := range transactions {
		description := strings.Replace(tr.Description, "\n", " ", -1)
		category := strconv.Itoa(tr.Mcc)
		bankCategory := strconv.Itoa(tr.Mcc)
		amount := fmt.Sprintf("%.2f", float64(tr.Amount)/100)
		unixTime := time.Unix(int64(tr.Time), 0).In(a.loc)
		date := unixTime.Format(dateTimeReportPattern)

		if c, err := a.mapCategory(catMap, category, description); err == nil {
			category = c
		}
		record := []string{date, description, category, bankCategory, amount}
		records = append(records, record)
	}

	buf := &bytes.Buffer{}
	wr := csv.NewWriter(buf)

	return a.writeRecords(buf, wr, records)
}

func (a *Transaction) writeRecords(r io.Reader, w *csv.Writer, record [][]string) (io.Reader, error) {
	if err := w.WriteAll(record); err != nil {
		return nil, errors.Errorf("can't write lines: lines=%v err=%v", record, err)
	}
	return r, nil
}

func (a Transaction) mapCategory(m categoryMapping, category, description string) (string, error) {
	if m == nil {
		return "", errors.New("empty mapping")
	}

	if m, ok := m[category+description]; ok {
		return m.App, nil
	}

	if m, ok := m[category]; ok {
		return m.App, nil
	}

	return "", errors.New("can't find mapping")
}

func (a *Transaction) getCategoryMapping(chatID int64) categoryMapping {
	key := fmt.Sprintf("%s%s", strconv.Itoa(int(chatID)), mappingSufix)
	categoryMapping, err := a.mappingRepo.Get(key) //Category mapping
	if err != nil {
		a.log.Error(err)
	}
	return categoryMapping
}

// Locale - return the transaction locale
func (a *Transaction) Locale() *time.Location {
	return a.loc
}

// ParseDate - parse date from string
func (a *Transaction) ParseDate(period string) (from time.Time, to time.Time, err error) {
	filter, err := a.getFilter(period)
	if err != nil {
		return
	}

	return filter.from, filter.to, nil
}
