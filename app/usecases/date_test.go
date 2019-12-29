package usecases

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

const errNotEqual = "not equal"

func TestDate_getFilter(t *testing.T) {
	RegisterTestingT(t)

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
			args: args{name: "01.08.2019-05.08.2019"},
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
		d, _ := NewDate(nil)

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
		d, _ := NewDate(nil)
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
