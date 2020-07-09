package usecases

import (
	"testing"
	"time"

	"github.com/jinzhu/now"

	. "github.com/onsi/gomega"
)

const errNotEqual = "not equal"

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
