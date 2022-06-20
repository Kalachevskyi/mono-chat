package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

// ChatUserUC - represents a use-case interface for processing business logic "ChatUser" use case.
type ChatUserUC interface {
	SetChatUserID(chatID int64, userID uuid.UUID) error
	GetChatUserID(chatID int64) (uuid.UUID, error)
}

// NewChatUser - builds "ChatUser" internal handler.
func NewChatUser(userUC ChatUserUC, botWrapper *BotWrapper) *ChatUser {
	return &ChatUser{
		userUC:     userUC,
		BotWrapper: botWrapper,
	}
}

// ChatUser - represents an internal handler for processing "ChatUser".
type ChatUser struct {
	userUC ChatUserUC
	*BotWrapper
}

// Handle - process the "ChatUser ID", send the result to the user.
func (a *ChatUser) Handle(u tg.Update) {
	userID, err := uuid.Parse(u.Message.CommandArguments())
	if err != nil {
		a.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	if err := a.userUC.SetChatUserID(u.Message.Chat.ID, userID); err != nil {
		a.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	a.sendMSG(tg.NewMessage(u.Message.Chat.ID, "successfully set ChatUser ID"))
}
