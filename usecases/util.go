package usecases

import (
	"sync"

	"gitlab.com/Kalachevskyi/mono-chat/entities"
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
