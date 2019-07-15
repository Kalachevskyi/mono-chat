package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/Kalachevskyi/mono-chat/entities"
)

const monoDomain = "https://api.monobank.ua"

type Logger interface {
	Errorf(template string, args ...interface{})
}

func NewMono(log Logger) *Mono {
	return &Mono{log: log}
}

type Mono struct {
	log Logger
}

func (m *Mono) GetTransactions(token string, from, to time.Time) ([]entities.Transaction, error) {
	url := fmt.Sprintf("%s/personal/statement/0/%d/%d", monoDomain, from.Unix(), to.Unix())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set("X-Token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer m.close(resp.Body)

	transactions := make([]entities.Transaction, 0)
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, errors.WithStack(err)
	}

	return transactions, nil
}

func (m *Mono) close(c io.Closer) {
	if err := c.Close(); err != nil {
		m.log.Errorf("%+v", err)
	}
}
