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
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Kalachevskyi/mono-chat/entities"
	"github.com/pkg/errors"
)

const mappingSufix = "_mapping"

// MappingRepo - represents Mapping repository interface
type MappingRepo interface {
	Set(key string, val map[string]entities.CategoryMapping) error
	Get(key string) (map[string]entities.CategoryMapping, error)
}

// NewMapping - builds mapping use-case
func NewMapping(mappingRepo MappingRepo, telegramRepo TelegramRepo) *Mapping {
	return &Mapping{
		mappingRepo:  mappingRepo,
		TelegramRepo: telegramRepo,
	}
}

// Mapping - represents category mapping  use-case for processing category
type Mapping struct {
	mappingRepo MappingRepo
	TelegramRepo
}

// Validate - validate file name by suffix ".csv"
func (c *Mapping) Validate(name string) error {
	if !strings.HasSuffix(name, csvSuffix) {
		return errors.New(`chat can only be processed using the file "csv"`)
	}
	return nil
}

// Parse - parse category mapping file, save it in repository
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
