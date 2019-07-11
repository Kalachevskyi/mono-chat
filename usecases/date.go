package usecases

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//Date patterns
const (
	yearMonthPattern      = "01.2006"
	dateTimePattern       = "02.01.2006T15.04"
	dateTimeReportPattern = "02.01.2006 15:04:05"
	timeFromPattern       = "T00.00"
	timeToPattern         = "T23.59"
)

//Regexp patterns
const (
	dateShortRegexpPattern = `^\d{2}-\d{2}` //range of date for the current month
	dateRegexpPattern      = `\d{2}\.\d{2}\.\d{4}-\d{2}\.\d{2}\.\d{4}`
	dateTimeRegexpPattern  = `\d{2}\.\d{2}\.\d{4}T\d{2}\.\d{2}-\d{2}\.\d{2}\.\d{4}T\d{2}\.\d{2}`
)

type Date struct {
	dateShortRegexp *regexp.Regexp
	dateRegexp      *regexp.Regexp
	dateTimeRegexp  *regexp.Regexp
	loc             *time.Location
}

func (d *Date) Init() error {
	r, err := regexp.Compile(dateShortRegexpPattern)
	if err != nil {
		return errors.WithStack(err)
	}
	d.dateShortRegexp = r

	r, err = regexp.Compile(dateRegexpPattern)
	if err != nil {
		return errors.WithStack(err)
	}
	d.dateRegexp = r

	r, err = regexp.Compile(dateTimeRegexpPattern)
	if err != nil {
		return errors.WithStack(err)
	}
	d.dateTimeRegexp = r

	loc, err := time.LoadLocation(timeLocation)
	if err != nil {
		return errors.Wrap(err, "can't set time location")
	}
	d.loc = loc

	return nil
}

func (d *Date) getFilter(name string) (*filter, error) {
	name = strings.TrimSuffix(name, csvSuffix)
	dates := strings.Split(name, "-")
	if len(dates) != 2 { //2 - must have two date by - separator
		return nil, errors.New("can't split dates by -")
	}

	if d.dateShortRegexp.MatchString(name) {
		yearMonth := time.Now().Format(yearMonthPattern)
		from := fmt.Sprintf("%s.%s%s", dates[0], yearMonth, timeFromPattern)
		to := fmt.Sprintf("%s.%s%s", dates[1], yearMonth, timeToPattern)
		return d.parseTime(from, to, timeDurationDay)
	}

	if d.dateRegexp.MatchString(name) {
		return d.parseTime(dates[0]+timeFromPattern, dates[1]+timeToPattern, timeDurationDay)
	}

	if d.dateTimeRegexp.MatchString(name) {
		return d.parseTime(dates[0], dates[1], time.Minute)
	}

	return nil, errors.New("can't find date pattern")
}

//parseTime - parse range of time according layout
func (d *Date) parseTime(fromStr, toStr string, tr time.Duration) (*filter, error) {
	errMsq := "can't parse time"
	from, err := time.ParseInLocation(dateTimePattern, fromStr, d.loc)
	if err != nil {
		return nil, errors.Wrap(err, errMsq)
	}

	to, err := time.ParseInLocation(dateTimePattern, toStr, d.loc)
	if err != nil {
		return nil, errors.Wrap(err, errMsq)
	}

	filter := filter{
		start:    from,
		end:      to,
		truncate: tr,
	}

	return &filter, nil
}
