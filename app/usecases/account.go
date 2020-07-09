package usecases

import "fmt"

// AccountRepo - represents AccountRepo repository
type AccountRepo interface {
	Set(key, account string) error
	Get(key string) (string, error)
}

// NewAccount - constructor for Account use case
func NewAccount(repo AccountRepo) *Account {
	return &Account{repo: repo}
}

// Account - represents Account use case
type Account struct {
	repo AccountRepo
}

// Get - returns account
func (a Account) Get(chatID int64) (string, error) {
	key := fmt.Sprintf("account_%v", chatID)
	return a.repo.Get(key)
}

// Set - save account
func (a Account) Set(chatID int64, account string) error {
	key := fmt.Sprintf("account_%v", chatID)
	return a.repo.Set(key, account)
}
