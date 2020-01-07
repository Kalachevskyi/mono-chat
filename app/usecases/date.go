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
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// TimeLocation - application location
const timeLocation = "Europe/Kiev"

//timeDurationDay - time duration for days
const timeDurationDay = 24 * time.Hour

//Date patterns
const (
	yearMonthPattern      = "01.2006"
	dateTimePattern       = "02.01.2006T15.04"
	dateTimeReportPattern = "02.01.2006 15:04:05"
	timeFromPattern       = "T00.00"
	timeToPattern         = "T23.59"
)

// Default regexp patterns for Date
const (
	ddRegexp        = `^\d{1,2}-\d{1,2}` //range of date for the current month
	ddmmyyyyRegexp  = `\d{2}\.\d{2}\.\d{4}-\d{2}\.\d{2}\.\d{4}`
	ddmmyyyytRegexp = `\d{2}\.\d{2}\.\d{4}T\d{2}\.\d{2}-\d{2}\.\d{2}\.\d{4}T\d{2}\.\d{2}`
)

var (
	dateShortRegexp = regexp.MustCompile(ddRegexp)        //nolint:gochecknoglobals
	dateRegexp      = regexp.MustCompile(ddmmyyyyRegexp)  //nolint:gochecknoglobals
	dateTimeRegexp  = regexp.MustCompile(ddmmyyyytRegexp) //nolint:gochecknoglobals
)

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
	dates := strings.Split(name, "-")
	if len(dates) != 2 { //2 - must have two date by - separator
		return nil, errors.New("can't split dates by -")
	}

	if dateShortRegexp.MatchString(name) {
		yearMonth := time.Now().Format(yearMonthPattern)
		fromRaw := dates[0]
		if len(fromRaw) == 1 {
			fromRaw = "0" + fromRaw
		}
		toRaw := dates[1]
		if len(toRaw) == 1 {
			toRaw = "0" + toRaw
		}
		from := fmt.Sprintf("%s.%s%s", fromRaw, yearMonth, timeFromPattern)
		to := fmt.Sprintf("%s.%s%s", toRaw, yearMonth, timeToPattern)
		return d.parseTime(from, to, timeDurationDay)
	}

	if dateRegexp.MatchString(name) {
		return d.parseTime(dates[0]+timeFromPattern, dates[1]+timeToPattern, timeDurationDay)
	}

	if dateTimeRegexp.MatchString(name) {
		return d.parseTime(dates[0], dates[1], time.Minute)
	}

	return nil, errors.New("can't find date pattern")
}

//parseTime - parse range of time according layout
func (d Date) parseTime(fromStr, toStr string, tr time.Duration) (*filter, error) {
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
		from:     from,
		to:       to,
		truncate: tr,
	}

	return &filter, nil
}
