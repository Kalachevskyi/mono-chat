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

// Date - precompiled regex for dates, time location
type Date struct {
	dateShortRegexp *regexp.Regexp
	dateRegexp      *regexp.Regexp
	dateTimeRegexp  *regexp.Regexp
	loc             *time.Location
}

// Init - compile regex for date, load location
func (d *Date) Init(dateShort, date, dateTime, location string) error {
	r, err := regexp.Compile(dateShort)
	if err != nil {
		return errors.WithStack(err)
	}
	d.dateShortRegexp = r

	r, err = regexp.Compile(date)
	if err != nil {
		return errors.WithStack(err)
	}
	d.dateRegexp = r

	r, err = regexp.Compile(dateTime)
	if err != nil {
		return errors.WithStack(err)
	}
	d.dateTimeRegexp = r

	loc, err := time.LoadLocation(location)
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
		from:     from,
		to:       to,
		truncate: tr,
	}

	return &filter, nil
}
