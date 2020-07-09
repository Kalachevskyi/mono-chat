package api

// AccountUC - represents a use-case interface for processing business logic "Account" use case
type AccountUC interface {
	Get(chatID int64) (string, error)
	Set(chatID int64, account string) error
}
