package telegram

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Error messages.
const defaultErrMSG = "Sorry, I can't process this message, view the logs or contact the owner of the service."

const dateTimePattern = "02.01.2006T15.04"

// Command messages.
const (
	getCommand          = "get"
	todayCommand        = "today"
	currentMonthCommand = "month"
	tokenCommand        = "token"
	accountCommand      = "account"
	infoCommand         = "info"
	userCommand         = "user"
)

// Logger - represents the application's logger interface.
type Logger interface {
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}

// Handler - represents internal handler  interface.
type Handler interface {
	Handle(u tg.Update)
}

// NewChat - builds main chat handler.
func NewChat(updates tg.UpdatesChannel, handlers map[HandlerKey]Handler, botWrapper *BotWrapper) *Chat {
	return &Chat{
		updates:    updates,
		handlers:   handlers,
		BotWrapper: botWrapper,
	}
}

// Chat - main chat handler.
type Chat struct {
	updates  tg.UpdatesChannel
	handlers map[HandlerKey]Handler
	*BotWrapper
}

// Handle - routes between internal handlers depending on the type of message.
func (c *Chat) Handle() {
	for u := range c.updates {
		if u.Message == nil { // ignore any non-Message Updates
			continue
		}

		if u.Message.Document != nil {
			switch u.Message.Document.FileName {
			case "mapping.csv":
				c.handle(MappingHandler, u)

				continue
			default:
				c.handle(FileReportHandler, u)

				continue
			}
		}

		if u.Message.Command() != "" {
			switch u.Message.Command() {
			case getCommand, todayCommand, currentMonthCommand:
				c.handle(TransactionsHandler, u)

				continue
			case tokenCommand:
				c.handle(TokenHandler, u)

				continue
			case accountCommand:
				c.handle(AccountHandler, u)

				continue
			case infoCommand:
				c.handle(ClientInfoHandler, u)

				continue
			case userCommand:
				c.handle(ChatUserHandler, u)

				continue
			}
		}

		c.sendMSG(tg.NewMessage(u.Message.Chat.ID, defaultErrMSG))
	}
}

func (c *Chat) handle(key HandlerKey, u tg.Update) {
	if h, ok := c.handlers[key]; ok {
		h.Handle(u)
	}
}
