package usecases

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	csvSuffix    = ".csv"
	errWriteLine = "can't write line: line=%v err=%v"
	reportLines  = 10
)

// TelegramRepo - represents Telegram repository interface.
type TelegramRepo interface {
	GetFile(u *url.URL) (io.ReadCloser, error)
}

type filter struct {
	from     time.Time
	to       time.Time
	truncate time.Duration
}

// NewFileReport - builds File report use-case.
func NewFileReport(date *Date, mappingRepo MappingRepo, log Logger, telegramRepo TelegramRepo) *FileReport {
	return &FileReport{
		date:         date,
		mappingRepo:  mappingRepo,
		log:          log,
		TelegramRepo: telegramRepo,
	}
}

// FileReport - represents File report use-case for processing file report.
type FileReport struct {
	date        *Date
	mappingRepo MappingRepo
	log         Logger
	TelegramRepo
}

// Validate - validate file name by suffix ".csv".
func (c *FileReport) Validate(name string) error {
	if !strings.HasSuffix(name, csvSuffix) {
		return errors.New(`chat can only be processed using the file "csv"`)
	}

	return nil
}

// Parse - parse MonoBank "csv" report, convert it to application format.
func (c *FileReport) Parse(userID uuid.UUID, fileName string, r io.Reader) (io.Reader, error) {
	lines, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, errors.Errorf("can't read file: err=%s", err)
	}

	// Set header to file
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

	key := fmt.Sprintf("%s_%s", mappingKey, userID)

	catMap, err := c.mappingRepo.Get(key) // Category mapping
	if err != nil {
		c.log.Error(err)
	}

	filter, err := c.date.getFilter(fileName)
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		if len(line) != reportLines {
			return nil, errors.New("report template does not match, should be 10")
		}

		date, category, bankCategory, description, amount := line[0], line[2], line[2], line[1], line[3]
		description = strings.ReplaceAll(description, "\n", " ")

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
