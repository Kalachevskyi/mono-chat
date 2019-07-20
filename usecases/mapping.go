package usecases

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Kalachevskyi/mono-chat/entities"
	"github.com/pkg/errors"
)

const mappingSufix = "_mapping"

type MappingRepo interface {
	Set(key string, val map[string]entities.CategoryMapping) error
	Get(key string) (map[string]entities.CategoryMapping, error)
}

func NewMapping(mappingRepo MappingRepo, telegramRepo TelegramRepo) *Mapping {
	return &Mapping{
		mappingRepo:  mappingRepo,
		TelegramRepo: telegramRepo,
	}
}

type Mapping struct {
	mappingRepo MappingRepo
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

	key := fmt.Sprintf("%s%s", strconv.Itoa(int(chatID)), mappingSufix)
	return c.mappingRepo.Set(key, mapping)
}
