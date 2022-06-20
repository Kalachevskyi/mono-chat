package mono

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/Kalachevskyi/mono-chat/app/model"
)

const (
	domainMono   = "https://api.monobank.ua"
	tokenMonoKey = "X-Token"
)

// NewMono - builds Mono repository.
func NewMono(log Logger) *Mono {
	return &Mono{log: log}
}

// Mono - represents the Mono repository for getting transaction from MonoBank telegram.
type Mono struct {
	log Logger
}

// GetTransactions - return Transactions from MonoBank.
func (m *Mono) GetTransactions(token, account string, from, to time.Time) ([]model.Transaction, error) {
	url := fmt.Sprintf("%s/personal/statement/%s/%d/%d", domainMono, account, from.Unix(), to.Unix())

	req, err := http.NewRequest(http.MethodGet, url, nil) // nolint:noctx
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set(tokenMonoKey, token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer closeBody(resp.Body, m.log)

	transactions := make([]model.Transaction, 0)
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, errors.WithStack(err)
	}

	return transactions, nil
}

// GetClientInfo - returns information about accounts (card, currency).
func (m Mono) GetClientInfo(token string) (c model.ClientInfo, err error) {
	url := fmt.Sprintf("%s/personal/client-info", domainMono)

	req, err := http.NewRequest(http.MethodGet, url, nil) // nolint:noctx
	if err != nil {
		return c, errors.WithStack(err)
	}

	req.Header.Set(tokenMonoKey, token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c, errors.WithStack(err)
	}

	defer closeBody(resp.Body, m.log)

	clientInfo := model.ClientInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&clientInfo); err != nil {
		return c, errors.WithStack(err)
	}

	return clientInfo, nil
}
