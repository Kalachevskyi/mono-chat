package usecases

import (
	"fmt"

	"github.com/google/uuid"
)

const accountKey = "account"

// AccountRepo - represents AccountRepo repository.
type AccountRepo interface {
	Set(key, account string) error
	Get(key string) (string, error)
}

// NewAccount - constructor for Account use case.
func NewAccount(repo AccountRepo) *Account {
	return &Account{repo: repo}
}

// Account - represents Account use case.
type Account struct {
	repo AccountRepo
}

// Get - returns account.
func (a Account) Get(userID uuid.UUID) (string, error) {
	key := fmt.Sprintf("%s_%v", accountKey, userID)

	return a.repo.Get(key)
}

// Set - save account.
func (a Account) Set(userID uuid.UUID, account string) error {
	key := fmt.Sprintf("%s_%v", accountKey, userID)

	return a.repo.Set(key, account)
}
