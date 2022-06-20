package usecases

import (
	"fmt"

	"github.com/google/uuid"
)

// UserRepo - represents AccountRepo repository.
type UserRepo interface {
	Set(key string) error
	CheckUser(key string) (bool, error)
}

// NewUser - constructor for User use case.
func NewUser(repo UserRepo) *User {
	return &User{repo: repo}
}

// User - represents User use case.
type User struct {
	repo UserRepo
}

// CheckUser - returns true if user exists.
func (a User) CheckUser(userID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("%s_%v", userKey, userID)

	return a.repo.CheckUser(key)
}

// Set - save account key.
func (a User) Set(userID uuid.UUID) error {
	key := fmt.Sprintf("%s_%v", userKey, userID)

	return a.repo.Set(key)
}
