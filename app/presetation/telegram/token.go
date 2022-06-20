package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

// TokenUC - represents a usecase interface for processing "Token" business logic.
type TokenUC interface {
	Set(userID uuid.UUID, token string) error
	Get(userID uuid.UUID) (string, error)
}

// NewToken - builds "NewToken" internal handler.
func NewToken(tokenUC TokenUC, chatUserUC ChatUserUC, botWrapper *BotWrapper) *Token {
	return &Token{
		tokenUC:    tokenUC,
		chatUserUC: chatUserUC,
		BotWrapper: botWrapper,
	}
}

// Token - represents an internal handler for processing "Token".
type Token struct {
	tokenUC    TokenUC
	chatUserUC ChatUserUC
	*BotWrapper
}

// Handle - process the "Token", send the result to the user.
func (t *Token) Handle(u tg.Update) { // nolint:dupl
	userID, err := t.chatUserUC.GetChatUserID(u.Message.Chat.ID)
	if err != nil {
		t.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	if err := t.tokenUC.Set(userID, u.Message.CommandArguments()); err != nil {
		t.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}
	t.sendMSG(tg.NewMessage(u.Message.Chat.ID, "successfully set token"))
}
