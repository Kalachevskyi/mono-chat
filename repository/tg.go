package repository

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func NewChat() *Chat {
	return &Chat{}
}

type Chat struct{}

func (c *Chat) GetFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {

		return nil, errors.Errorf("can't get file by url: %s", url)
	}

	return resp.Body, nil
}
