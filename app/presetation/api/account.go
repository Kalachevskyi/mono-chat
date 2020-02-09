// Package api is the business logic layer of the application.
package api

import tg "github.com/go-telegram-bot-api/telegram-bot-api"

// NewAccount - builds "NewAccount" internal handler
func NewAccount(accountUC AccountUC, botWrapper *BotWrapper) *Account {
	return &Account{
		accountUC:  accountUC,
		BotWrapper: botWrapper,
	}
}

// Account - represents an internal handler for processing "Account"
type Account struct {
	accountUC AccountUC
	*BotWrapper
}

// Handle - process the "Token", send the result to the user
func (a *Account) Handle(u tg.Update) {
	if err := a.accountUC.Set(u.Message.Chat.ID, u.Message.CommandArguments()); err != nil {
		a.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	a.sendMSG(tg.NewMessage(u.Message.Chat.ID, "successfully set account"))
}
