package rest

import (
	"io"
	"time"

	"github.com/google/uuid"
)

// Logger - represents the application's logger interface.
type Logger interface {
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}

// TransactionUC - represents a use-case interface for processing business logic of "MonoBank" transactions API.
type TransactionUC interface {
	GetTransactions(token, account string, userID uuid.UUID, from time.Time, to time.Time) (io.Reader, error)
	ParseDate(period string) (from time.Time, to time.Time, err error)
	Locale() *time.Location
}

// UserUC represents User use-case interface.
type UserUC interface {
	Set(userID uuid.UUID) error
	CheckUser(userID uuid.UUID) (bool, error)
}

// AccountUC - represents a use-case interface for processing business logic "Account" use case.
type AccountUC interface {
	Get(userID uuid.UUID) (string, error)
	Set(userID uuid.UUID, account string) error
}

// TokenUC - represents a usecase interface for processing "Token" business logic.
type TokenUC interface {
	Set(userID uuid.UUID, token string) error
	Get(userID uuid.UUID) (string, error)
}
