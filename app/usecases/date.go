package usecases

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/now"

	"github.com/pkg/errors"
)

// TimeLocation - application location
const timeLocation = "Europe/Kiev"

//timeDurationDay - time duration for days
const timeDurationDay = 24 * time.Hour

//Date patterns
const (
	yearMonthPattern = "01.2006"
	datePattern      = "02.01.2006"
	dateTimePattern  = "02.01.2006T15.04"
)

// Default regexp patterns for Date
const (
	ddPattern        = `^\d{1,2}-\d{1,2}` //range of date for the current month/year
	ddmmyyyyPattern  = `\d{2}\.\d{2}\.\d{4}-\d{2}\.\d{2}\.\d{4}`
	ddmmyyyytPattern = `\d{2}\.\d{2}\.\d{4}T\d{2}\.\d{2}-\d{2}\.\d{2}\.\d{4}T\d{2}\.\d{2}`
)

var (
	ddRegexp        = regexp.MustCompile(ddPattern)        //nolint:gochecknoglobals
	ddmmyyyyRegexp  = regexp.MustCompile(ddmmyyyyPattern)  //nolint:gochecknoglobals
	ddmmyyyytRegexp = regexp.MustCompile(ddmmyyyytPattern) //nolint:gochecknoglobals
)

const errParseTime = "can't parse time"

// NewDate - Date type constructor
// "loc" - can be empty, default parameter "Europe/Kiev"
func NewDate(loc *time.Location) (*Date, error) {
	if loc != nil {
		return &Date{loc: loc}, nil
	}

	loc, err := time.LoadLocation(timeLocation)
	if err != nil {
		return nil, errors.Wrap(err, "can't set time location")
	}

	return &Date{loc: loc}, nil
}

// Date - precompiled regex for dates, time location
type Date struct {
	loc *time.Location
}

func (d Date) getFilter(name string) (*filter, error) {
	name = strings.TrimSuffix(name, csvSuffix)
	datesSetLen := 2

	dates := strings.Split(name, "-")
	if len(dates) != datesSetLen { //2 - must have two date by - separator
		return nil, errors.New("can't split dates by -")
	}

	if ddRegexp.MatchString(name) {
		return d.parseDate(d.prepareDays(dates[0], dates[1]))
	}

	if ddmmyyyyRegexp.MatchString(name) {
		return d.parseDate(dates[0], dates[1])
	}

	if ddmmyyyytRegexp.MatchString(name) {
		return d.parseDateTime(dates[0], dates[1])
	}

	return nil, errors.New("can't find date pattern")
}

// prepareDays - prepare "from, to" dates adding current month/year
func (d Date) prepareDays(fromStr, toStr string) (from, to string) {
	yearMonth := time.Now().Format(yearMonthPattern)
	prefix := "0"
	minLen := 1

	if len(fromStr) == minLen {
		fromStr = prefix + fromStr
	}

	if len(toStr) == minLen {
		toStr = prefix + toStr
	}

	from = fmt.Sprintf("%s.%s", fromStr, yearMonth)
	to = fmt.Sprintf("%s.%s", toStr, yearMonth)

	return from, to
}

//parseDate - parse range of time according layout, exclude time
func (d Date) parseDate(fromStr, toStr string) (*filter, error) {
	from, err := time.ParseInLocation(datePattern, fromStr, d.loc)
	if err != nil {
		return nil, errors.Wrap(err, errParseTime)
	}

	to, err := time.ParseInLocation(datePattern, toStr, d.loc)
	if err != nil {
		return nil, errors.Wrap(err, errParseTime)
	}

	filter := filter{
		from:     now.New(from).BeginningOfDay(),
		to:       now.New(to).EndOfDay(),
		truncate: timeDurationDay,
	}

	return &filter, nil
}

//parseDateTime - parse range of time according layout
func (d Date) parseDateTime(fromStr, toStr string) (*filter, error) {
	from, err := time.ParseInLocation(dateTimePattern, fromStr, d.loc)
	if err != nil {
		return nil, errors.Wrap(err, errParseTime)
	}

	to, err := time.ParseInLocation(dateTimePattern, toStr, d.loc)
	if err != nil {
		return nil, errors.Wrap(err, errParseTime)
	}

	filter := filter{
		from:     from,
		to:       to,
		truncate: time.Minute,
	}

	return &filter, nil
}
