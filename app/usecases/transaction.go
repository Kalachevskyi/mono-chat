package usecases

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Kalachevskyi/mono-chat/app/model"
	"github.com/pkg/errors"
)

const accuracy = 100

//go:generate mockgen -destination=./mono_mock_test.go -package=usecases -source=./transaction.go

// Logger - represents the application's logger interface
type Logger interface {
	Error(args ...interface{})
}

// MonoRepo - represents Transaction repository interface
type MonoRepo interface {
	GetTransactions(token, account string, from, to time.Time) ([]model.Transaction, error)
}

type categoryMapping map[string]model.CategoryMapping

// NewTransaction - builds Transaction report use-case
func NewTransaction(trRepo MonoRepo, mapRepo MappingRepo, log Logger, date Date) *Transaction {
	return &Transaction{
		apiRepo:     trRepo,
		mappingRepo: mapRepo,
		log:         log,
		Date:        date,
	}
}

// Transaction - represents Transaction  use-case for processing bank Transactions
type Transaction struct {
	apiRepo     MonoRepo
	mappingRepo MappingRepo
	log         Logger
	Date
}

// GetTransactions - get bank transactions, convert it to app csv report
func (a *Transaction) GetTransactions(token, account string, chatID int64, from, to time.Time) (io.Reader, error) {
	transactions, err := a.apiRepo.GetTransactions(token, account, from, to)
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
		amount := fmt.Sprintf("%.2f", float64(tr.Amount)/accuracy)
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

// Locale - return the transaction localeapi/chat.go:89
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
