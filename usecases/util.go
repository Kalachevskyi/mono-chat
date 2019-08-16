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

// ReportHeader - report header enum
type ReportHeader int

func (r ReportHeader) Str() string { return months[r-1] }

const (
	DateHeader ReportHeader = 1 + iota
	DescriptionHeader
	CategoryHeader
	BankCategoryHeader
	AmountHeader
)

var months = [...]string{ //nolint:gochecknoglobals
	"Date",
	"Description",
	"Category",
	"Bank category",
	"Amount",
}
