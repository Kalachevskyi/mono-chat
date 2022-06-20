package telegram

import (
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// NewTelegram - builds Telegram repository.
func NewTelegram() *Telegram {
	return &Telegram{}
}

// Telegram - represents the Telegram repository for communication with Telegram telegram.
type Telegram struct{}

// GetFile -  get the file from Telegram REST API, makes HTTP call to telegram API.
func (c *Telegram) GetFile(u *url.URL) (io.ReadCloser, error) {
	resp, err := http.Get(u.String()) //nolint
	if err != nil {
		return nil, errors.Errorf("can't get file by url: %s", u.String())
	}

	return resp.Body, nil
}
