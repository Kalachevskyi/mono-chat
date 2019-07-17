package usecases

import (
	"strings"
	"sync"

	"github.com/Kalachevskyi/mono-chat/entities"
	"github.com/pkg/errors"
)

//Save user category mapping in memory
type categoryMappingT struct {
	v map[int64]map[string]entities.CategoryMapping
	sync.Mutex
}

var categoryMapping = categoryMappingT{v: make(map[int64]map[string]entities.CategoryMapping)}

//ReportHeader - report header enum
type ReportHeader int

func (r ReportHeader) Str() string { return months[r-1] }

const (
	DateHeader ReportHeader = 1 + iota
	DescriptionHeader
	CategoryHeader
	BankCategoryHeader
	AmountHeader
)

var months = [...]string{
	"Date",
	"Description",
	"Category",
	"Bank category",
	"Amount",
}

func Validate(name string) error {
	if !strings.HasSuffix(name, csvSuffix) {
		return errors.New(`chat can only be processed using the file "csv"`)
	}
	return nil
}
