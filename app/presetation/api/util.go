package api

import (
	"fmt"
	"io"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

// HandlerKey - type for naming handlers
type HandlerKey int

// Handlers keys
const (
	FileReportHandler HandlerKey = 1 + iota
	MappingHandler
	TransactionsHandler
	TokenHandler
	ClientInfoHandler
	AccountHandler
)

// ErrStack - add stack to error, work with "github.com/pkg/errors" package
func ErrStack(err error) error {
	if err == nil {
		return err
	}

	return fmt.Errorf("%+v", err)
}

func close(closer io.Closer, log Logger) {
	if err := closer.Close(); err != nil {
		log.Errorf("handlers.Chat.close: can't close body: err=%s", ErrStack(err))
	}
}

// NewBotWrapper - builds "NewBotWrapper"
func NewBotWrapper(bot *tg.BotAPI, log Logger) *BotWrapper {
	return &BotWrapper{bot: bot, log: log}
}

// BotWrapper - represents the Telegram chatbot wrapper, register an error if it occurred
type BotWrapper struct {
	bot *tg.BotAPI
	log Logger
}

func (c *BotWrapper) sendDefaultErr(chatID int64, err error) {
	c.log.Error(ErrStack(err))
	c.sendMSG(tg.NewMessage(chatID, defaultErrMSG))
}

func (c *BotWrapper) sendMSG(msg tg.Chattable) {
	if _, err := c.bot.Send(msg); err != nil {
		c.log.Errorf("can't send err message: err=%+v", errors.WithStack(err))
	}
}
