package handlers

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Error messages
const defaultErrMSG = "Sorry, I can't process this message, view the logs or contact the owner of the service."

//Command messages
const (
	getCommand   = "get"
	todayCommand = "today"
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
	Today() (from, to time.Time, err error)
}

func NewChat(bot *tg.BotAPI, uConf tg.UpdateConfig, log Logger, csvUC CsvUC, apiUC ApiUC) *Chat {
	return &Chat{
		bot:          bot,
		updateConfig: uConf,
		log:          log,
		csvUC:        csvUC,
		apiUC:        apiUC,
	}
}

type Chat struct {
	bot          *tg.BotAPI
	updateConfig tg.UpdateConfig
	log          Logger
	csvUC        CsvUC
	apiUC        ApiUC
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
		c.handleGetCommand(u)
		return
	case todayCommand:
		c.HandleTodayCommand(u)
		return
	default:
		c.sendMSG(tg.NewMessage(u.Message.Chat.ID, defaultErrMSG))
	}
}

func (c *Chat) handleGetCommand(u tg.Update) {
	from, to, err := c.apiUC.ParseDate(u.Message.CommandArguments())
	if err != nil {
		c.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileResp, err := c.apiUC.GetTransactions(u.Message.Chat.ID, from, to)
	if err != nil {
		c.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}
	name := fmt.Sprintf("%s.csv", u.Message.CommandArguments())
	reader := tg.FileReader{
		Name:   name,
		Reader: fileResp,
		Size:   -1,
	}
	msg := tg.NewDocumentUpload(u.Message.Chat.ID, reader)
	c.sendMSG(msg)
}

func (c *Chat) HandleTodayCommand(u tg.Update) {
	from, to, err := c.apiUC.Today()
	if err != nil {
		c.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileResp, err := c.apiUC.GetTransactions(u.Message.Chat.ID, from, to)
	if err != nil {
		c.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}
	reader := tg.FileReader{
		Name:   "today.csv",
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
