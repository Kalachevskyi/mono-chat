package usecases_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/Kalachevskyi/mono-chat/app/model"
	uc "github.com/Kalachevskyi/mono-chat/app/usecases"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

const errNotEqual = "not equal"

func TestMapping_Parse(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)

	type fields struct {
		mappingRepo func() uc.MappingRepo
	}

	type args struct {
		chatID int64
		r      func() io.Reader
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: `test-case1: error with csv reader`,
			fields: fields{
				mappingRepo: func() uc.MappingRepo {
					return NewMockMappingRepo(mockCtrl)
				}},
			args: args{
				chatID: 0,
				r: func() io.Reader {
					data, err := ioutil.ReadFile("./testdata/mapping_file_error.json")
					Ω(err).To(BeNil(), errNotEqual)
					return bytes.NewReader(data)
				},
			},
			wantErr: true,
		},
		{
			name: `test-case2: error, not equal number of columns`,
			fields: fields{
				mappingRepo: func() uc.MappingRepo {
					return NewMockMappingRepo(mockCtrl)
				}},
			args: args{
				chatID: 0,
				r: func() io.Reader {
					data, err := ioutil.ReadFile("./testdata/mapping_len_line_error.csv")
					Ω(err).To(BeNil(), errNotEqual)
					return bytes.NewReader(data)
				},
			},
			wantErr: true,
		},
		{
			name: `test-case3: success`,
			fields: fields{
				mappingRepo: func() uc.MappingRepo {
					mapping := map[string]model.CategoryMapping{
						"4111": {
							Mono: "4111",
							App:  "Transport",
						},
						"7230": {
							Mono: "7230",
							App:  "Hair care",
						},
					}
					repo := NewMockMappingRepo(mockCtrl)
					repo.EXPECT().Set("0_mapping", mapping)
					return repo
				},
			},
			args: args{
				chatID: 0,
				r: func() io.Reader {
					data, err := ioutil.ReadFile("./testdata/mapping.csv")
					Ω(err).To(BeNil(), errNotEqual)
					return bytes.NewReader(data)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		mappingRepo := tt.fields.mappingRepo()
		m := uc.NewMapping(mappingRepo, nil)
		err := m.Parse(tt.args.chatID, tt.args.r())
		Ω(err != nil).To(Equal(tt.wantErr), errNotEqual)
	}
}

func TestMapping_Validate(t *testing.T) {
	RegisterTestingT(t)
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    `test-case1: error with validation csv format of the file`,
			args:    args{name: "file"},
			wantErr: true,
		},
		{
			name:    `test-case2: success`,
			args:    args{name: "file.csv"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		m := &uc.Mapping{}

		err := m.Validate(tt.args.name)
		if tt.wantErr {
			Ω(err).NotTo(BeNil(), errNotEqual)
			continue
		} else {
			Ω(err).To(BeNil(), errNotEqual)
		}
	}
}

func TestNewMapping(t *testing.T) {
	RegisterTestingT(t)
	want := &uc.Mapping{}
	got := uc.NewMapping(nil, nil)
	Ω(got).To(Equal(want), fmt.Sprintf(errDefaultMsg, got))
}
