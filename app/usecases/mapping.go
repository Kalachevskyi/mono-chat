package usecases

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/Kalachevskyi/mono-chat/app/model"
)

const (
	mappingKey   = "mapping"
	mappingLines = 3
)

//go:generate mockgen -destination=./mapping_mock_test.go -package=usecases -source=./mapping.go

// MappingRepo - represents Mapping repository interface.
type MappingRepo interface {
	Set(key string, val map[string]model.CategoryMapping) error
	Get(key string) (map[string]model.CategoryMapping, error)
}

// NewMapping - builds mapping use-case.
func NewMapping(mappingRepo MappingRepo, telegramRepo TelegramRepo) *Mapping {
	return &Mapping{
		mappingRepo:  mappingRepo,
		TelegramRepo: telegramRepo,
	}
}

// Mapping - represents category mapping  use-case for processing category.
type Mapping struct {
	mappingRepo MappingRepo
	TelegramRepo
}

// Validate - validate file name by suffix ".csv".
func (c *Mapping) Validate(name string) error {
	if !strings.HasSuffix(name, csvSuffix) {
		return errors.New(`chat can only be processed using the file "csv"`)
	}

	return nil
}

// Parse - parse category mapping file, save it in repository.
func (c *Mapping) Parse(userID uuid.UUID, r io.Reader) error {
	lines, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return errors.Errorf("can't read file: err=%s", err)
	}

	mapping := make(map[string]model.CategoryMapping)
	for _, line := range lines {
		if len(line) != mappingLines {
			return errors.New("mapping should have 3 column")
		}

		categoryMapping := model.CategoryMapping{
			Mono:        line[0],
			Description: line[1],
			App:         line[2],
		}

		key := line[0] + line[1]

		mapping[key] = categoryMapping
	}

	key := fmt.Sprintf("%s_%s", mappingKey, userID)

	return c.mappingRepo.Set(key, mapping)
}
