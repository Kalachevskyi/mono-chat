package handlers

type TokenUC interface {
	Set(chatID int64, token string) error
	Get(chatID int64) (string, error)
}
