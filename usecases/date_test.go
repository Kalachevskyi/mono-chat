package usecases

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/Kalachevskyi/mono-chat/config"
	. "github.com/onsi/gomega"
)

const errNotEqual = "not equal"

func TestDate_Init(t *testing.T) {
	RegisterTestingT(t)
	type args struct {
		// regexp patterns
		dateShort string
		date      string
		dateTime  string
		// time location
		location string
	}

	type fields struct {
		dateShortRegexp func(a args) *regexp.Regexp
		dateRegexp      func(a args) *regexp.Regexp
		dateTimeRegexp  func(a args) *regexp.Regexp
		loc             func(a args) *time.Location
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test-case1: compile date error",
			args: args{
				date: "fail(",
			},
			wantErr: true,
		},
		{
			name: "test-case2: compile dateShort error",
			args: args{
				dateShort: "fail(",
			},
			wantErr: true,
		},
		{
			name: "test-case3: compile dateTime error",
			args: args{
				dateTime: "fail(",
			},
			wantErr: true,
		},
		{
			name: "test-case4: compile location error",
			args: args{
				location: "fail(",
			},
			wantErr: true,
		},
		{
			name: "test-case5: success",
			args: args{
				dateShort: "success",
				date:      "success",
				dateTime:  "success",
				location:  "Europe/Kiev",
			},
			fields: fields{
				dateShortRegexp: func(a args) *regexp.Regexp {
					result, err := regexp.Compile(a.dateShort)
					Ω(err).To(BeNil(), errNotEqual)
					return result
				},
				dateRegexp: func(a args) *regexp.Regexp {
					result, err := regexp.Compile(a.date)
					Ω(err).To(BeNil(), errNotEqual)
					return result
				},
				dateTimeRegexp: func(a args) *regexp.Regexp {
					result, err := regexp.Compile(a.dateTime)
					Ω(err).To(BeNil(), errNotEqual)
					return result
				},
				loc: func(a args) *time.Location {
					result, err := time.LoadLocation(a.location)
					Ω(err).To(BeNil(), errNotEqual)
					return result
				},
			},
		},
	}
	for _, tt := range tests {
		d := &Date{}
		if err := d.Init(tt.args.dateShort, tt.args.date, tt.args.dateTime, tt.args.location); tt.wantErr {
			Ω(err).NotTo(BeNil(), errNotEqual)
			continue
		} else {
			Ω(err).To(BeNil(), errNotEqual)
		}
		Ω(d.dateShortRegexp).To(Equal(tt.fields.dateShortRegexp(tt.args)), errNotEqual)
		Ω(d.dateRegexp).To(Equal(tt.fields.dateRegexp(tt.args)), errNotEqual)
		Ω(d.dateTimeRegexp).To(Equal(tt.fields.dateTimeRegexp(tt.args)), errNotEqual)
		Ω(d.loc).To(Equal(tt.fields.loc(tt.args)), errNotEqual)
	}
}

func TestDate_getFilter(t *testing.T) {
	RegisterTestingT(t)
	dateShort := config.DateShortRegexpPattern
	date := config.DateRegexpPattern
	dateTime := config.DateTimeRegexpPattern
	location := config.TimeLocation

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    func(*time.Location) *filter
		wantErr bool
	}{
		{
			name:    `test-case1: error with "can't split dates by -""`,
			args:    args{name: "not valid"},
			wantErr: true,
		},
		{
			name: `test-case2: success with short date`,
			args: args{name: "01-05"},
			want: func(loc *time.Location) *filter {
				yearMonth := time.Now().Format(yearMonthPattern)
				fromStr := fmt.Sprintf("%s.%s%s", "01", yearMonth, timeFromPattern)
				toStr := fmt.Sprintf("%s.%s%s", "05", yearMonth, timeToPattern)
				from, err := time.ParseInLocation(dateTimePattern, fromStr, loc)
				Ω(err).To(BeNil(), errNotEqual)
				to, err := time.ParseInLocation(dateTimePattern, toStr, loc)
				Ω(err).To(BeNil(), errNotEqual)
				return &filter{
					from:     from,
					to:       to,
					truncate: timeDurationDay,
				}
			},
		},
		{
			name: `test-case3: success with date`,
			args: args{name: "01.08.2019-06.08.2019"},
			want: func(loc *time.Location) *filter {
				from, err := time.ParseInLocation(dateTimePattern, "01.08.2019"+timeFromPattern, loc)
				Ω(err).To(BeNil(), errNotEqual)
				to, err := time.ParseInLocation(dateTimePattern, "05.08.2019"+timeToPattern, loc)
				Ω(err).To(BeNil(), errNotEqual)
				return &filter{
					from:     from,
					to:       to,
					truncate: timeDurationDay,
				}
			},
		},
		{
			name: `test-case4: success with date time`,
			args: args{name: "01.08.2019T15.00-05.08.2019T21.00"},
			want: func(loc *time.Location) *filter {
				from, err := time.ParseInLocation(dateTimePattern, "01.08.2019T15.00", loc)
				Ω(err).To(BeNil(), errNotEqual)
				to, err := time.ParseInLocation(dateTimePattern, "05.08.2019T21.00", loc)
				Ω(err).To(BeNil(), errNotEqual)
				return &filter{
					from:     from,
					to:       to,
					truncate: time.Minute,
				}
			},
		},
		{
			name:    `test-case5: error with can't find date pattern`,
			args:    args{name: "01.08.2019T15:00-05.08.2019T21:00"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		d := &Date{}

		err := d.Init(dateShort, date, dateTime, location)
		Ω(err).To(BeNil(), errNotEqual)

		filter, err := d.getFilter(tt.args.name)
		if tt.wantErr {
			Ω(err).NotTo(BeNil(), errNotEqual)
			continue
		} else {
			Ω(err).To(BeNil(), errNotEqual)
		}

		Ω(filter).To(Equal(tt.want(d.loc)), errNotEqual)
	}
}

func TestDate_parseTime(t *testing.T) {
	RegisterTestingT(t)
	dateShort := config.DateShortRegexpPattern
	date := config.DateRegexpPattern
	dateTime := config.DateTimeRegexpPattern
	location := config.TimeLocation

	type args struct {
		fromStr string
		toStr   string
		tr      time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    func(*time.Location) *filter
		wantErr bool
	}{
		{
			name:    `test-case1: error with parsing "from" parameter`,
			args:    args{fromStr: "fail"},
			wantErr: true,
		},
		{
			name:    `test-case1: error with parsing "to" parameter`,
			args:    args{fromStr: "02.01.2006T15.04", toStr: "fail"},
			wantErr: true,
		},
		{
			name: `test-case3: success`,
			args: args{
				fromStr: "01.08.2019T15.00",
				toStr:   "05.08.2019T21.00",
				tr:      time.Minute,
			},
			want: func(loc *time.Location) *filter {
				from, err := time.ParseInLocation(dateTimePattern, "01.08.2019T15.00", loc)
				Ω(err).To(BeNil(), errNotEqual)
				to, err := time.ParseInLocation(dateTimePattern, "05.08.2019T21.00", loc)
				Ω(err).To(BeNil(), errNotEqual)
				return &filter{
					from:     from,
					to:       to,
					truncate: time.Minute,
				}
			},
		},
	}
	for _, tt := range tests {
		d := &Date{}
		err := d.Init(dateShort, date, dateTime, location)
		Ω(err).To(BeNil(), errNotEqual)

		got, err := d.parseTime(tt.args.fromStr, tt.args.toStr, tt.args.tr)
		if tt.wantErr {
			Ω(err).NotTo(BeNil(), errNotEqual)
			continue
		} else {
			Ω(err).To(BeNil(), errNotEqual)
		}

		Ω(got).To(Equal(tt.want(d.loc)), errNotEqual)

	}
}
