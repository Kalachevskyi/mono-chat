package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

// AccountUC - represents a use-case interface for processing business logic "Account" use case.
type AccountUC interface {
	Get(userID uuid.UUID) (string, error)
	Set(userID uuid.UUID, account string) error
}

// NewAccount - builds "NewAccount" internal handler.
func NewAccount(accountUC AccountUC, chatUserUC ChatUserUC, botWrapper *BotWrapper) *Account {
	return &Account{
		accountUC:  accountUC,
		chatUserUC: chatUserUC,
		BotWrapper: botWrapper,
	}
}

// Account - represents an internal handler for processing "Account".
type Account struct {
	accountUC  AccountUC
	chatUserUC ChatUserUC
	*BotWrapper
}

// Handle - process the "Token", send the result to the user.
func (a *Account) Handle(u tg.Update) { // nolint:dupl
	userID, err := a.chatUserUC.GetChatUserID(u.Message.Chat.ID)
	if err != nil {
		a.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	if err := a.accountUC.Set(userID, u.Message.CommandArguments()); err != nil {
		a.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	a.sendMSG(tg.NewMessage(u.Message.Chat.ID, "successfully set account"))
}
