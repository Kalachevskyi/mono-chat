package usecases

import (
	"fmt"

	"github.com/google/uuid"
)

// ChatUserRepo - represents ChatUser repository interface.
type ChatUserRepo interface {
	Set(key, token string) error
	Get(key string) (string, error)
}

const (
	userKey = "user"
	chatKey = "chat"
)

// NewChatUser constructor for ChatUser.
func NewChatUser(userRepo ChatUserRepo) *ChatUser {
	return &ChatUser{userRepo: userRepo}
}

// ChatUser represents user chat use-case.
type ChatUser struct {
	userRepo ChatUserRepo
}

// SetChatUserID set chat ID with user ID.
func (u ChatUser) SetChatUserID(chatID int64, userID uuid.UUID) error {
	key := fmt.Sprintf("%s_%s_%v", chatKey, userKey, chatID)

	return u.userRepo.Set(key, userID.String())
}

// GetChatUserID get user ID by chat ID.
func (u ChatUser) GetChatUserID(chatID int64) (uuid.UUID, error) {
	key := fmt.Sprintf("%s_%s_%v", chatKey, userKey, chatID)
	val, err := u.userRepo.Get(key)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
