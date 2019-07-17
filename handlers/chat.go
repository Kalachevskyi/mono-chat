package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
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

type ChatBuilder struct {
	Updates  tg.UpdatesChannel
	TokeUC   TokenUC
	Handlers map[HandlerKey]Handler
}

func (c *ChatBuilder) Build() *Chat {
	return &Chat{
		updates:  c.Updates,
		tokeUC:   c.TokeUC,
		handlers: c.Handlers,
	}
}

type Chat struct {
	updates  tg.UpdatesChannel
	tokeUC   TokenUC
	handlers map[HandlerKey]Handler
	BotWrapper
}

func (c *Chat) Handle() {
	for u := range c.updates {
		chatID := u.Message.Chat.ID

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
			msg := "can't find handler %v"
			c.sendDefaultErr(chatID, errors.Errorf(msg, FileReportHandler))
			continue
		}

		if u.Message.Command() != "" {
			switch u.Message.Command() {
			case getCommand, todayCommand, currentMonthCommand:
				if h, ok := c.handlers[Transactions]; ok {
					h.Handle(u)
					continue
				}
			case tokenCommand:
				if err := c.tokeUC.Set(u.Message.Chat.ID, u.Message.CommandArguments()); err != nil {
					c.sendDefaultErr(u.Message.Chat.ID, err)
					return
				}
				c.sendMSG(tg.NewMessage(u.Message.Chat.ID, "successfully set token"))
				return
			}
		}

		c.sendMSG(tg.NewMessage(u.Message.Chat.ID, defaultErrMSG))
	}
}
