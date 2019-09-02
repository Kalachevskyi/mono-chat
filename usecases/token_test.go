package usecases_test

import (
	"errors"
	"fmt"
	"testing"

	uc "github.com/Kalachevskyi/mono-chat/usecases"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
)

const errDefaultMsg = `unexpected value "%v"`

func TestNewToken(t *testing.T) {
	RegisterTestingT(t)
	want := &uc.Token{}
	got := uc.NewToken(nil)
	Ω(got).To(Equal(want), fmt.Sprintf(errDefaultMsg, got))
}

func TestToken_Get(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	type fields struct {
		repo func(chatID int64, want string) uc.TokenRepo
	}
	type args struct {
		chatID int64
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		err     string
	}{
		{
			name: "test-case1: success execution",
			fields: fields{repo: func(chatID int64, want string) uc.TokenRepo {
				repo := NewMockTokenRepo(mockCtrl)
				key := fmt.Sprintf("token_%v", chatID)
				repo.EXPECT().Get(key).Return(want, nil).Times(1)
				return repo
			}},
			args:    args{1},
			want:    "l1lms13d0vc8ks",
			wantErr: false,
		},
		{
			name: "test-case2: repo error",
			fields: fields{repo: func(chatID int64, want string) uc.TokenRepo {
				repo := NewMockTokenRepo(mockCtrl)
				key := fmt.Sprintf("token_%v", chatID)
				err := errors.New("some error")
				repo.EXPECT().Get(key).Return("", err).Times(1)
				return repo
			}},
			args:    args{1},
			want:    "",
			wantErr: true,
			err:     "some error",
		},
	}
	for _, tt := range tests {
		c := uc.NewToken(tt.fields.repo(tt.args.chatID, tt.want))
		got, err := c.Get(tt.args.chatID)
		if tt.wantErr {
			Ω(err).NotTo(BeNil(), fmt.Sprintf(errDefaultMsg, err))
			Ω(err.Error()).To(Equal(tt.err), fmt.Sprintf(errDefaultMsg, err.Error()))
		}
		Ω(got).To(Equal(tt.want), fmt.Sprintf(errDefaultMsg, got))
	}
}

func TestToken_Set(t *testing.T) {
	RegisterTestingT(t)
	mockCtrl := gomock.NewController(t)
	type fields struct {
		repo func(chatID int64, token string) uc.TokenRepo
	}
	type args struct {
		chatID int64
		token  string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     string
	}{
		{
			name: "test-case1: success execution",
			fields: fields{repo: func(chatID int64, token string) uc.TokenRepo {
				repo := NewMockTokenRepo(mockCtrl)
				key := fmt.Sprintf("token_%v", chatID)
				repo.EXPECT().Set(key, token).Return(nil).Times(1)
				return repo
			}},
			args:    args{1, "l1lms13d0vc8ks"},
			wantErr: false,
		},
		{
			name: "test-case2: repo error",
			fields: fields{repo: func(chatID int64, token string) uc.TokenRepo {
				repo := NewMockTokenRepo(mockCtrl)
				key := fmt.Sprintf("token_%v", chatID)
				err := errors.New("some error")
				repo.EXPECT().Set(key, token).Return(err).Times(1)
				return repo
			}},
			args:    args{1, "l1lms13d0vc8ks"},
			wantErr: true,
			err:     "some error",
		},
	}
	for _, tt := range tests {
		tokeRepo := tt.fields.repo(tt.args.chatID, tt.args.token)
		c := uc.NewToken(tokeRepo)
		err := c.Set(tt.args.chatID, tt.args.token)
		if tt.wantErr {
			Ω(err).NotTo(BeNil(), fmt.Sprintf(errDefaultMsg, err))
			Ω(err.Error()).To(Equal(tt.err), fmt.Sprintf(errDefaultMsg, err.Error()))
		}
	}
}
