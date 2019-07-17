package repository

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func NewTelegram() *Telegram {
	return &Telegram{}
}

type Telegram struct{}

func (c *Telegram) GetFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {

		return nil, errors.Errorf("can't get file by url: %s", url)
	}

	return resp.Body, nil
}
