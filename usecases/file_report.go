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

	"github.com/pkg/errors"
)

const csvSuffix = ".csv"

const errWriteLine = "can't write line: line=%v err=%v"

// TelegramRepo - represents Telegram repository interface
type TelegramRepo interface {
	GetFile(url string) (io.ReadCloser, error)
}

type filter struct {
	from     time.Time
	to       time.Time
	truncate time.Duration
}

// NewFileReport - builds File report use-case
func NewFileReport(date Date, mappingRepo MappingRepo, log Logger, telegramRepo TelegramRepo) *FileReport {
	return &FileReport{
		date:         date,
		mappingRepo:  mappingRepo,
		log:          log,
		TelegramRepo: telegramRepo,
	}
}

// FileReport - represents File report use-case for processing file report
type FileReport struct {
	date        Date
	mappingRepo MappingRepo
	log         Logger
	TelegramRepo
}

// Validate - validate file name by suffix ".csv"
func (c *FileReport) Validate(name string) error {
	if !strings.HasSuffix(name, csvSuffix) {
		return errors.New(`chat can only be processed using the file "csv"`)
	}
	return nil
}

// Parse - parse MonoBank "csv" report, convert it to application format
func (c *FileReport) Parse(chatID int64, fileName string, r io.Reader) (io.Reader, error) {
	lines, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, errors.Errorf("can't read file: err=%s", err)
	}

	//Set header to file
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
		return nil, errors.Errorf(errWriteLine, header, err)
	}

	key := fmt.Sprintf("%s%s", strconv.Itoa(int(chatID)), mappingSufix)
	catMap, err := c.mappingRepo.Get(key) //Category mapping
	if err != nil {
		c.log.Error(err)
	}

	filter, err := c.date.getFilter(fileName)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if len(line) != 10 { //10 - default columns in a mono bank template
			return nil, errors.New("report template does not match, should be 10")
		}

		date, category, bankCategory, description, amount := line[0], line[2], line[2], line[1], line[3]
		description = strings.Replace(description, "\n", " ", -1)

		if filter != nil {
			dateTime, err := time.Parse(dateTimeReportPattern, date)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			dateTime = dateTime.Truncate(filter.truncate)
			if ok := c.applyFilter(dateTime, *filter); !ok {
				continue
			}
		}

		if catMap != nil {
			if m, ok := catMap[category+description]; ok {
				category = m.App
			} else if m, ok := catMap[category]; ok {
				category = m.App
			}
		}

		record := []string{date, description, category, bankCategory, amount}
		if err := wr.Write(record); err != nil {
			return nil, errors.Errorf(errWriteLine, record, err)
		}
	}
	wr.Flush()
	return buf, nil
}

func (c *FileReport) applyFilter(d time.Time, f filter) bool {
	if d.Equal(f.from) || d.Equal(f.to) {
		return true
	}

	if d.Before(f.from) {
		return false
	}

	if d.After(f.to) {
		return false
	}

	return true
}
