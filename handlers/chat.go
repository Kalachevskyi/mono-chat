package handlers

import (
	"fmt"
	"io"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/now"
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

type CsvUC interface {
	Validate(name string) error
	GetFile(url string) (io.ReadCloser, error)
	ParseMapping(chatID int64, r io.Reader) error
	ParseReport(chatID int64, fileName string, r io.Reader) (io.Reader, error)
}

type ApiUC interface {
	GetTransactions(chatID int64, from time.Time, to time.Time) (io.Reader, error)
	ParseDate(period string) (from time.Time, to time.Time, err error)
}

type TokenUC interface {
	SaveToken(chatID int64, token string) error
}

type ChatBuilder struct {
	Bot          *tg.BotAPI
	UpdateConfig tg.UpdateConfig
	Log          Logger
	CsvUC        CsvUC
	ApiUC        ApiUC
	TokeUC       TokenUC
}

func (c *ChatBuilder) Build() *Chat {
	return &Chat{
		bot:          c.Bot,
		updateConfig: c.UpdateConfig,
		log:          c.Log,
		csvUC:        c.CsvUC,
		apiUC:        c.ApiUC,
		tokeUC:       c.TokeUC,
	}
}

type Chat struct {
	bot          *tg.BotAPI
	updateConfig tg.UpdateConfig
	log          Logger
	csvUC        CsvUC
	apiUC        ApiUC
	tokeUC       TokenUC
}

func (c *Chat) Handle() {
	updates, err := c.bot.GetUpdatesChan(c.updateConfig)
	if err != nil {
		c.log.Errorf("can't get updates: %v", ErrStack(err))
	}

	for u := range updates {
		if u.Message == nil { // ignore any non-Message Updates
			continue
		}

		if u.Message.Document != nil {
			c.handleFile(u)
			continue
		}

		if u.Message.Command() != "" {
			c.handleCommands(u)
			continue
		}

		c.sendMSG(tg.NewMessage(u.Message.Chat.ID, defaultErrMSG))
	}
}

func (c *Chat) handleCommands(u tg.Update) {
	switch u.Message.Command() {
	case getCommand:
		from, to, err := c.apiUC.ParseDate(u.Message.CommandArguments())
		if err != nil {
			c.sendDefaultErr(u.Message.Chat.ID, err)
			return
		}
		c.handlePeriodCommand(u, from, to)
		return
	case todayCommand:
		c.handlePeriodCommand(u, now.BeginningOfDay(), now.EndOfDay())
		return
	case currentMonthCommand:
		c.handlePeriodCommand(u, now.BeginningOfMonth(), now.EndOfMonth())
		return
	case tokenCommand:

	default:
		c.sendMSG(tg.NewMessage(u.Message.Chat.ID, defaultErrMSG))
	}
}

func (c *Chat) handlePeriodCommand(u tg.Update, from, to time.Time) {

	fileResp, err := c.apiUC.GetTransactions(u.Message.Chat.ID, from, to)
	if err != nil {
		c.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}
	reader := tg.FileReader{
		Name:   fmt.Sprintf("%s-%s%s", from.Format(dateTimePattern), to.Format(dateTimePattern), ".csv"),
		Reader: fileResp,
		Size:   -1,
	}
	msg := tg.NewDocumentUpload(u.Message.Chat.ID, reader)
	c.sendMSG(msg)
}

func (c *Chat) handleFile(u tg.Update) {
	if err := c.csvUC.Validate(u.Message.Document.FileName); err != nil {
		c.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileTG, err := c.bot.GetFile(tg.FileConfig{FileID: u.Message.Document.FileID})
	if err != nil {
		c.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	file, err := c.csvUC.GetFile(fileTG.Link(c.bot.Token))
	if err != nil {
		c.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	switch u.Message.Document.FileName {
	case "mapping.csv":
		if err := c.csvUC.ParseMapping(u.Message.Chat.ID, file); err != nil {
			c.sendDefaultErr(u.Message.Chat.ID, err)
			return
		}
		c.sendMSG(tg.NewMessage(u.Message.Chat.ID, "mapping successfully loaded"))
	default:
		fileResp, err := c.csvUC.ParseReport(u.Message.Chat.ID, u.Message.Document.FileName, file)
		if err != nil {
			c.sendDefaultErr(u.Message.Chat.ID, err)
			return
		}
		name := u.Message.Document.FileName
		reader := tg.FileReader{
			Name:   name,
			Reader: fileResp,
			Size:   -1,
		}
		msg := tg.NewDocumentUpload(u.Message.Chat.ID, reader)
		c.sendMSG(msg)
	}
	c.close(file)
}

func (c *Chat) sendDefaultErr(chatID int64, err error) {
	c.log.Error(ErrStack(err))
	c.sendMSG(tg.NewMessage(chatID, defaultErrMSG))
}

func (c *Chat) sendMSG(msg tg.Chattable) {
	if _, err := c.bot.Send(msg); err != nil {
		c.log.Errorf("can't send err message: err=%+v", errors.WithStack(err))
	}
}

func (c *Chat) close(closer io.Closer) {
	if err := closer.Close(); err != nil {
		c.log.Errorf("handlers.Chat.close: can't close body: err=%s", ErrStack(err))
	}
}
