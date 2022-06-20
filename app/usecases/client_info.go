package usecases

import "github.com/Kalachevskyi/mono-chat/app/model"

// ClientInfoRepo - represents ClientInfo repository.
type ClientInfoRepo interface {
	GetClientInfo(token string) (c model.ClientInfo, err error)
}

// NewClientInfo - ClientInfo constructor.
func NewClientInfo(repo ClientInfoRepo) *ClientInfo {
	return &ClientInfo{repo: repo}
}

// ClientInfo - represents ClientInfo use case.
type ClientInfo struct {
	repo ClientInfoRepo
}

// GetClientInfo - returns client info.
func (c ClientInfo) GetClientInfo(token string) (model.ClientInfo, error) {
	return c.repo.GetClientInfo(token)
}
