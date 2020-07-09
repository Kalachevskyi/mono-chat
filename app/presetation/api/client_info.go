package api

import (
	"fmt"

	"github.com/Kalachevskyi/mono-chat/app/model"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

const accuracy = 100

// ClientInfoUC - represents Client Info use case
type ClientInfoUC interface {
	GetClientInfo(token string) (model.ClientInfo, error)
}

// NewClientInfo - represents ClientInfo constructor
func NewClientInfo(tokenUC TokenUC, clientInfoUC ClientInfoUC, botWrapper *BotWrapper) *ClientInfo {
	return &ClientInfo{
		tokenUC:      tokenUC,
		clientInfoUC: clientInfoUC,
		BotWrapper:   botWrapper,
	}
}

// ClientInfo - represents ClientInfo handler struct
type ClientInfo struct {
	tokenUC      TokenUC
	clientInfoUC ClientInfoUC
	*BotWrapper
}

// Handle  - represents ClientInfo handler
func (c *ClientInfo) Handle(u tg.Update) {
	chatID := u.Message.Chat.ID

	token, err := c.tokenUC.Get(chatID)
	if err != nil {
		c.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	clientInfo, err := c.clientInfoUC.GetClientInfo(token)
	if err != nil {
		c.sendDefaultErr(u.Message.Chat.ID, err)
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
