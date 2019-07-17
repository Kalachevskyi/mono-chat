package usecases

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/Kalachevskyi/mono-chat/entities"
	"github.com/pkg/errors"
)

func NewMapping(telegramRepo TelegramRepo) *Mapping {
	return &Mapping{TelegramRepo: telegramRepo}
}

type Mapping struct {
	TelegramRepo
}

func (c *Mapping) Validate(name string) error {
	if !strings.HasSuffix(name, csvSuffix) {
		return errors.New(`chat can only be processed using the file "csv"`)
	}
	return nil
}

func (c *Mapping) Parse(chatID int64, r io.Reader) error {
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
