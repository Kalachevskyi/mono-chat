package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TokenUC interface {
	Set(chatID int64, token string) error
	Get(chatID int64) (string, error)
}

func NewToken(tokenUC TokenUC, botWrapper *BotWrapper) *Token {
	return &Token{
		tokenUC:    tokenUC,
		BotWrapper: botWrapper,
	}
}

type Token struct {
	tokenUC TokenUC
	*BotWrapper
}

func (t *Token) Handle(u tg.Update) {
	if err := t.tokenUC.Set(u.Message.Chat.ID, u.Message.CommandArguments()); err != nil {
		t.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}
	t.sendMSG(tg.NewMessage(u.Message.Chat.ID, "successfully set token"))
}
