package telegram

import (
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/Kalachevskyi/mono-chat/app/model"
)

const accuracy = 100

// ClientInfoUC - represents Client Info use case.
type ClientInfoUC interface {
	GetClientInfo(token string) (model.ClientInfo, error)
}

// NewClientInfo - represents ClientInfo constructor.
func NewClientInfo(tokenUC TokenUC, clientInfoUC ClientInfoUC, chatUserUC ChatUserUC, botWrapper *BotWrapper) *ClientInfo {
	return &ClientInfo{
		tokenUC:      tokenUC,
		clientInfoUC: clientInfoUC,
		chatUserUC:   chatUserUC,
		BotWrapper:   botWrapper,
	}
}

// ClientInfo - represents ClientInfo handler struct.
type ClientInfo struct {
	tokenUC      TokenUC
	clientInfoUC ClientInfoUC
	chatUserUC   ChatUserUC
	*BotWrapper
}

// Handle  - represents ClientInfo handler.
func (c *ClientInfo) Handle(u tg.Update) {
	userID, err := c.chatUserUC.GetChatUserID(u.Message.Chat.ID)
	if err != nil {
		c.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	chatID := u.Message.Chat.ID
	token, err := c.tokenUC.Get(userID)
	if err != nil {
		c.sendDefaultErr(chatID, err)

		return
	}

	clientInfo, err := c.clientInfoUC.GetClientInfo(token)
	if err != nil {
		c.sendDefaultErr(chatID, err)

		return
	}

	var resp string
	for _, val := range clientInfo.Accounts {
		resp = fmt.Sprintf("%sid: %s\n", resp, val.ID)
		resp = fmt.Sprintf("%s    currency code: %v\n", resp, val.CurrencyCode)
		resp = fmt.Sprintf("%s    balance: %v\n", resp, val.Balance/accuracy)
		resp = fmt.Sprintf("%s    type: %v\n\n", resp, val.Type)
	}

	msg := tg.NewMessage(chatID, resp)
	c.sendMSG(msg)
}
