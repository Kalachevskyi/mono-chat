package telegram

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// NewTelegram - builds Telegram repository
func NewTelegram(log Logger) *Telegram {
	return &Telegram{
		log: log,
	}
}

// Telegram - represents the Telegram repository for communication with Telegram api
type Telegram struct {
	log Logger
}

// GetFile -  get the file from Telegram REST API, makes HTTP call to telegram API
func (c *Telegram) GetFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url) //nolint
	if err != nil {
		return nil, errors.Errorf("can't get file by url: %s", url)
	}

	return resp.Body, nil
}
