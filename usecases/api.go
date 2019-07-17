package usecases

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/Kalachevskyi/mono-chat/entities"
	"github.com/pkg/errors"
)

const timeLocation = "Europe/Kiev"

type ApiRepo interface {
	GetTransactions(token string, from time.Time, to time.Time) ([]entities.Transaction, error)
}

func NewApi(apiRepo ApiRepo, dateRegexp Date) *Api {
	return &Api{apiRepo: apiRepo, Date: dateRegexp}
}

type Api struct {
	apiRepo ApiRepo
	Date
}

func (a *Api) GetTransactions(token string, chatID int64, from time.Time, to time.Time) (io.Reader, error) {
	transactions, err := a.apiRepo.GetTransactions(token, from, to)
	if err != nil {
		return nil, err
	}

	header := []string{
		DateHeader.Str(),
		DescriptionHeader.Str(),
		CategoryHeader.Str(),
		BankCategoryHeader.Str(),
		AmountHeader.Str(),
	}
	buf := &bytes.Buffer{}
	wr := csv.NewWriter(buf)
	if err := wr.Write(header); err != nil {
		return nil, errors.Errorf("can't write line: line=%v err=%v", header, err)
	}

	for _, tr := range transactions {
		description := tr.Description
		category := strconv.Itoa(tr.Mcc)
		bankCategory := strconv.Itoa(tr.Mcc)

		categoryMapping.Lock()
		if mapping := categoryMapping.v[chatID]; mapping != nil {
			if m, ok := mapping[category+description]; ok {
				category = m.App
			} else if m, ok := mapping[category]; ok {
				category = m.App
			}
		}
		categoryMapping.Unlock()

		amount := fmt.Sprintf("%.2f", float64(tr.Amount)/100)
		unixTime := time.Unix(int64(tr.Time), 0).In(a.loc)
		date := unixTime.Format(dateTimeReportPattern)

		record := []string{date, description, category, bankCategory, amount}
		if err := wr.Write(record); err != nil {
			msg := "can't write line: line=%v err=%v"
			return nil, errors.Errorf(msg, record, err)
		}
	}
	wr.Flush()

	return buf, nil
}

func (a *Api) ParseDate(period string) (from time.Time, to time.Time, err error) {
	filter, err := a.getFilter(period)
	if err != nil {
		return
	}

	return filter.start, filter.end, nil
}
