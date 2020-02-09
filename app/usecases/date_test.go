package usecases

import (
	"testing"
	"time"

	"github.com/jinzhu/now"

	. "github.com/onsi/gomega"
)

const errNotEqual = "not equal"

//func TestDate_getFilter(t *testing.T) {
//	RegisterTestingT(t)
//	d, err := NewDate(nil)
//	Ω(err).To(BeNil(), errNotEqual)
//
//	type args struct {
//		name string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *filter
//		wantErr bool
//	}{
//		{
//			name:    `test-case1: error with "can't split dates by -""`,
//			args:    args{name: "not valid"},
//			wantErr: true,
//		},
//		{
//			name: `test-case2: success with short date 01-05`,
//			args: args{name: "01-05"},
//			want: &filter{
//				from:     time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, d.loc),
//				to:       time.Date(time.Now().Year(), time.Now().Month(), 5, 23, 59, 0, 0, d.loc),
//				truncate: timeDurationDay,
//			},
//		},
//		{
//			name: `test-case3: success with short date 1-05`,
//			args: args{name: "1-05"},
//			want: &filter{
//				from:     time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, d.loc),
//				to:       time.Date(time.Now().Year(), time.Now().Month(), 5, 23, 59, 0, 0, d.loc),
//				truncate: timeDurationDay,
//			},
//		},
//		{
//			name: `test-case4: success with short date 01-5`,
//			args: args{name: "01-5"},
//			want: &filter{
//				from:     time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, d.loc),
//				to:       time.Date(time.Now().Year(), time.Now().Month(), 5, 23, 59, 0, 0, d.loc),
//				truncate: timeDurationDay,
//			},
//		},
//		{
//			name: `test-case5: success with short date 1-5`,
//			args: args{name: "1-5"},
//			want: &filter{
//				from:     time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, d.loc),
//				to:       time.Date(time.Now().Year(), time.Now().Month(), 5, 23, 59, 0, 0, d.loc),
//				truncate: timeDurationDay,
//			},
//		},
//		{
//			name: `test-case6: success with date`,
//			args: args{name: "01.08.2019-05.08.2019"},
//			want: &filter{
//				from:     time.Date(2019, 8, 1, 0, 0, 0, 0, d.loc),
//				to:       time.Date(2019, 8, 5, 23, 59, 0, 0, d.loc),
//				truncate: timeDurationDay,
//			},
//		},
//		{
//			name: `test-case7: success with date time`,
//			args: args{name: "01.08.2019T15.00-05.08.2019T21.00"},
//			want: &filter{
//				from:     time.Date(2019, 8, 1, 15, 0, 0, 0, d.loc),
//				to:       time.Date(2019, 8, 5, 21, 0, 0, 0, d.loc),
//				truncate: time.Minute,
//			},
//		},
//		{
//			name:    `test-case8: error with can't find date pattern`,
//			args:    args{name: "01.08.2019T15:00-05.08.2019T21:00"},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		filter, err := d.getFilter(tt.args.name)
//		Ω(err != nil).To(Equal(tt.wantErr), errNotEqual)
//		Ω(filter).To(Equal(tt.want), errNotEqual)
//	}
//}

func TestDate_parseTime(t *testing.T) {
	RegisterTestingT(t)

	type args struct {
		fromStr string
		toStr   string
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
			want:    func(location *time.Location) *filter { return nil },
			wantErr: true,
		},
		{
			name:    `test-case1: error with parsing "to" parameter`,
			args:    args{fromStr: "02.01.2006T15.04", toStr: "fail"},
			want:    func(location *time.Location) *filter { return nil },
			wantErr: true,
		},
		{
			name: `test-case3: success`,
			args: args{
				fromStr: "01.08.2019",
				toStr:   "05.08.2019",
			},
			want: func(loc *time.Location) *filter {
				from, _ := time.ParseInLocation(datePattern, "01.08.2019", loc)
				from = now.New(from).BeginningOfDay()
				to, _ := time.ParseInLocation(datePattern, "05.08.2019", loc)
				to = now.New(to).EndOfDay()
				return &filter{
					from:     from,
					to:       to,
					truncate: timeDurationDay,
				}
			},
		},
	}
	for _, tt := range tests {
		d, _ := NewDate(nil)
		got, err := d.parseDate(tt.args.fromStr, tt.args.toStr)
		Ω(tt.wantErr).To(Equal(err != nil), errNotEqual)
		Ω(got).To(Equal(tt.want(d.loc)), errNotEqual)
	}
}
