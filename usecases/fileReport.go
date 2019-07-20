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

//timeDurationDay - time duration for days
const timeDurationDay = 24 * time.Hour

type TelegramRepo interface {
	GetFile(url string) (io.ReadCloser, error)
}

type filter struct {
	start    time.Time
	end      time.Time
	truncate time.Duration
}

func NewFileReport(date Date, mappingRepo MappingRepo, log Logger, telegramRepo TelegramRepo) *FileReport {
	return &FileReport{
		date:         date,
		mappingRepo:  mappingRepo,
		log:          log,
		TelegramRepo: telegramRepo,
	}
}

type FileReport struct {
	date        Date
	mappingRepo MappingRepo
	log         Logger
	TelegramRepo
}

func (c *FileReport) Validate(name string) error {
	if !strings.HasSuffix(name, csvSuffix) {
		return errors.New(`chat can only be processed using the file "csv"`)
	}
	return nil
}

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
		return nil, errors.Errorf("can't write line: line=%v err=%v", header, err)
	}

	key := fmt.Sprintf("%s%s", strconv.Itoa(int(chatID)), mappingSufix)
	catMap, err := c.mappingRepo.Get(key) //Category mapping
	if err != nil {
		c.log.Error(err)
	}

	filter, _ := c.date.getFilter(fileName)
	for _, line := range lines {
		if len(line) != 10 { //10 - default columns in a mono bank template
			return nil, errors.New("report template does not match, should be 10")
		}

		date, category, bankCategory, description, amount := line[0], line[2], line[2], line[1], line[3]

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
			msg := "can't write line: line=%v err=%v"
			return nil, errors.Errorf(msg, record, err)
		}
	}
	wr.Flush()
	return buf, nil
}

func (c *FileReport) applyFilter(d time.Time, f filter) bool {
	if d.Equal(f.start) || d.Equal(f.end) {
		return true
	}

	if d.Before(f.start) {
		return false
	}

	if d.After(f.end) {
		return false
	}

	return true
}
