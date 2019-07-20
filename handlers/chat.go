package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Error messages
const defaultErrMSG = "Sorry, I can't process this message, view the logs or contact the owner of the service."

const dateTimePattern = "02.01.2006T15.04"

//Command messages
const (
	getCommand          = "get"
	todayCommand        = "today"
	currentMonthCommand = "month"
	tokenCommand        = "token"
)

type Logger interface {
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}

type Handler interface {
	Handle(u tg.Update)
}

func NewChat(updates tg.UpdatesChannel, handlers map[HandlerKey]Handler, botWrapper *BotWrapper) *Chat {
	return &Chat{
		updates:    updates,
		handlers:   handlers,
		BotWrapper: botWrapper,
	}
}

type Chat struct {
	updates  tg.UpdatesChannel
	handlers map[HandlerKey]Handler
	*BotWrapper
}

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
			}
		}

		c.sendMSG(tg.NewMessage(u.Message.Chat.ID, defaultErrMSG))
	}
}
