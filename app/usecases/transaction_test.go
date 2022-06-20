package usecases_test

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"testing"
	"time"

	"github.com/Kalachevskyi/mono-chat/app/model"
	uc "github.com/Kalachevskyi/mono-chat/app/usecases"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

func TestNewTransaction(t *testing.T) {
	RegisterTestingT(t)
	want := &uc.Transaction{}
	got := uc.NewTransaction(nil, nil, nil, uc.Date{})
	Ω(got).To(Equal(want), fmt.Sprintf(errDefaultMsg, got))
}

func TestTransaction_GetTransactions(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	date, _ := uc.NewDateOld(nil)

	type args struct {
		token   string
		account string
		chatID  int64
		from    time.Time
		to      time.Time
	}
	type fields struct {
		apiRepo     func(args) uc.MonoRepo
		mappingRepo func(args) uc.MappingRepo
		log         func() uc.Logger
		Date        uc.Date
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    func() io.Reader
		wantErr bool
	}{
		{
			name: `test-case1: error from the repo GetTransactions`,
			fields: fields{
				apiRepo: func(a args) uc.MonoRepo {
					repo := NewMockMonoRepo(mockCtrl)
					err := errors.New("some error")
					repo.EXPECT().GetTransactions(a.token, a.account, a.from, a.to).Return(nil, err).Times(1)
					return repo
				},
				mappingRepo: func(a args) uc.MappingRepo { return nil },
				log:         func() uc.Logger { return nil },
			},
			args: args{
				token:   "some_token",
				account: "some_account",
			},
			wantErr: true,
		},
		{
			name: `test-case1: success execution`,
			fields: fields{
				apiRepo: func(a args) uc.MonoRepo {
					transactions := []model.Transaction{
						{
							ID:          "ZuHWzqkKGVo=",
							Mcc:         7997,
							Amount:      -95000,
							Time:        1554466347,
							Description: "Покупка щастя",
						},
					}
					repo := NewMockMonoRepo(mockCtrl)
					repo.EXPECT().GetTransactions(a.token, a.account, a.from, a.to).Return(transactions, nil).Times(1)
					return repo
				},
				mappingRepo: func(a args) uc.MappingRepo {
					key := fmt.Sprintf("%s%s", strconv.Itoa(int(a.chatID)), "_mapping")
					catMap := map[string]model.CategoryMapping{"7997": {}}
					repo := NewMockMappingRepo(mockCtrl)
					repo.EXPECT().Get(key).Return(catMap, nil).Times(1)
					return repo
				},
				log: func() uc.Logger { return nil },
			},
			args: args{
				token:   "some_token",
				account: "some_account",
			},
			want: func() io.Reader {
				records := [][]string{
					{
						uc.DateHeader.Str(),
						uc.DescriptionHeader.Str(),
						uc.CategoryHeader.Str(),
						uc.BankCategoryHeader.Str(),
						uc.AmountHeader.Str(),
					},
					{
						"05.04.2019 15:12:27",
						"Покупка щастя",
						"",
						"7997",
						"-950.00",
					},
				}
				buf := &bytes.Buffer{}
				wr := csv.NewWriter(buf)
				if err := wr.WriteAll(records); err != nil {
					Ω(err).To(BeNil(), errNotEqual)
				}
				return buf
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tr := uc.NewTransaction(tt.fields.apiRepo(tt.args), tt.fields.mappingRepo(tt.args), tt.fields.log(), *date)
		got, err := tr.GetTransactions(tt.args.token, tt.args.account, tt.args.chatID, tt.args.from, tt.args.to)
		if tt.wantErr {
			Ω(err).NotTo(BeNil(), errNotEqual)
			continue
		} else {
			Ω(err).To(BeNil(), errNotEqual)
		}
		Ω(got).To(Equal(tt.want()), fmt.Sprintf(errDefaultMsg, got))
	}
}
