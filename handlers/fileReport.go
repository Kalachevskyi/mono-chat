package handlers

import (
	"io"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CsvUC interface {
	Validate(name string) error
	GetFile(url string) (io.ReadCloser, error)
	Parse(chatID int64, fileName string, r io.Reader) (io.Reader, error)
}

func NewFileReport(csvUC CsvUC, botWrapper *BotWrapper) *FileReport {
	return &FileReport{csvUC: csvUC, BotWrapper: botWrapper}
}

type FileReport struct {
	csvUC CsvUC
	*BotWrapper
}

func (f *FileReport) Handle(u tg.Update) {
	if err := f.csvUC.Validate(u.Message.Document.FileName); err != nil {
		f.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileTG, err := f.bot.GetFile(tg.FileConfig{FileID: u.Message.Document.FileID})
	if err != nil {
		f.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	file, err := f.csvUC.GetFile(fileTG.Link(f.bot.Token))
	if err != nil {
		f.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileResp, err := f.csvUC.Parse(u.Message.Chat.ID, u.Message.Document.FileName, file)
	if err != nil {
		f.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	name := u.Message.Document.FileName
	reader := tg.FileReader{
		Name:   name,
		Reader: fileResp,
		Size:   -1,
	}
	msg := tg.NewDocumentUpload(u.Message.Chat.ID, reader)
	f.sendMSG(msg)

	close(file, f.log)
}
