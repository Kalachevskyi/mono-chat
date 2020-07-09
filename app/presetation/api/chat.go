package api

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Error messages
const defaultErrMSG = "Sorry, I can't process this message, view the logs or contact the owner of the service."

const dateTimePattern = "02.01.2006T15.04"

// Command messages
const (
	getCommand          = "get"
	todayCommand        = "today"
	currentMonthCommand = "month"
	tokenCommand        = "token"
	accountCommand      = "account"
	infoCommand         = "info"
)

// Logger - represents the application's logger interface
type Logger interface {
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}

// Handler - represents internal handler interface
type Handler interface {
	Handle(u tg.Update)
}

// NewChat - builds main chat handler
func NewChat(updates tg.UpdatesChannel, handlers map[HandlerKey]Handler, botWrapper *BotWrapper) *Chat {
	return &Chat{
		updates:    updates,
		handlers:   handlers,
		BotWrapper: botWrapper,
	}
}

// Chat - main chat handler
type Chat struct {
	updates  tg.UpdatesChannel
	handlers map[HandlerKey]Handler
	*BotWrapper
}

// Handle - routes between internal handlers depending on the type of message
func (c *Chat) Handle() {
	for u := range c.updates {
		if u.Message == nil { // ignore any non-Message Updates
			continue
		}

		if u.Message.Document != nil {
			switch u.Message.Document.FileName {
			case "mapping.csv":
				if h, ok := c.handlers[MappingHandler]; ok {
					h.Handle(u)
					continue
				}
			default:
				if h, ok := c.handlers[FileReportHandler]; ok {
					h.Handle(u)
					continue
				}
			}

			continue
		}

		if u.Message.Command() != "" {
			switch u.Message.Command() {
			case getCommand, todayCommand, currentMonthCommand:
				if h, ok := c.handlers[TransactionsHandler]; ok {
					h.Handle(u)
					continue
				}
			case tokenCommand:
				if h, ok := c.handlers[TokenHandler]; ok {
					h.Handle(u)
					continue
				}
				return
			case accountCommand:
				if h, ok := c.handlers[AccountHandler]; ok {
					h.Handle(u)

					continue
				}
				return
			case infoCommand:
				if h, ok := c.handlers[ClientInfoHandler]; ok {
					h.Handle(u)
					continue
				}
				return
			}
		}

		c.sendMSG(tg.NewMessage(u.Message.Chat.ID, defaultErrMSG))
	}
}
