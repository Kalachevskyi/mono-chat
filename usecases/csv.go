package usecases

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/Kalachevskyi/mono-chat/entities"
)

const csvSuffix = ".csv"

//timeDurationDay - time duration for days
const timeDurationDay = 24 * time.Hour

type ChatRepo interface {
	GetFile(url string) (io.ReadCloser, error)
}

type filter struct {
	start    time.Time
	end      time.Time
	truncate time.Duration
}

func NewChat(repo ChatRepo, d Date) *Chat {
	return &Chat{
		repo: repo,
		date: d,
	}
}

type Chat struct {
	repo ChatRepo
	date Date
}

func (c *Chat) Validate(name string) error {
	if !strings.HasSuffix(name, csvSuffix) {
		return errors.New(`chat can only be processed using the file "csv"`)
	}
	return nil
}

func (c *Chat) GetFile(url string) (io.ReadCloser, error) {
	return c.repo.GetFile(url)
}

func (c *Chat) ParseMapping(chatID int64, r io.Reader) error {
	lines, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return errors.Errorf("can't read file: err=%s", err)
	}

	mapping := make(map[string]entities.CategoryMapping)
	for _, line := range lines {
		if len(line) != 3 {
			return errors.New("mapping should have 3 column")
		}

		categoryMapping := entities.CategoryMapping{
			Mono:        line[0],
			Description: line[1],
			App:         line[2],
		}

		key := line[0] + line[1]

		mapping[key] = categoryMapping
	}

	categoryMapping.Lock()
	categoryMapping.v[chatID] = mapping
	categoryMapping.Unlock()
	return nil
}

func (c *Chat) ParseReport(chatID int64, fileName string, r io.Reader) (io.Reader, error) {
	lines, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, errors.Errorf("can't read file: err=%s", err)
	}

	//Set header to file
	header := []string{DateHeader.Str(), DescriptionHeader.Str(), CategoryHeader.Str(), AmountHeader.Str()}
	buf := &bytes.Buffer{}
	wr := csv.NewWriter(buf)
	if err := wr.Write(header); err != nil {
		return nil, errors.Errorf("can't write line: line=%v err=%v", header, err)
	}

	filter, _ := c.date.getFilter(fileName)
	for _, line := range lines {
		if len(line) != 10 { //10 - default columns in a mono bank template
			return nil, errors.New("report template does not match, should be 10")
		}

		date, category, description, amount := line[0], line[2], line[1], line[3]

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

		categoryMapping.Lock()
		if mapping := categoryMapping.v[chatID]; mapping != nil {
			if m, ok := mapping[category+description]; ok {
				category = m.App
			} else if m, ok := mapping[category]; ok {
				category = m.App
			}
		}
		categoryMapping.Unlock()

		record := []string{date, description, category, amount}
		if err := wr.Write(record); err != nil {
			msg := "can't write line: line=%v err=%v"
			return nil, errors.Errorf(msg, record, err)
		}
	}
	wr.Flush()
	return buf, nil
}

func (c *Chat) applyFilter(d time.Time, f filter) bool {
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
