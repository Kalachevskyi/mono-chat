package usecases

// ReportHeader - report header enum.
type ReportHeader int

const dateTimeReportPattern = "02.01.2006 15:04:05"

// Str - returns a report header string.
func (r ReportHeader) Str() string { return months[r-1] }

// Headers for csv report.
const (
	DateHeader ReportHeader = 1 + iota
	DescriptionHeader
	CategoryHeader
	BankCategoryHeader
	AmountHeader
)

var months = [5]string{ //nolint:gochecknoglobals
	"Date",
	"Description",
	"Category",
	"Bank category",
	"Amount",
}
