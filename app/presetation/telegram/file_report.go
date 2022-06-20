package telegram

import (
	"io"
	"net/url"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

// CsvUC - represents a usecase interface for processing business logic of a CSV report.
type CsvUC interface {
	Validate(name string) error
	GetFile(u *url.URL) (io.ReadCloser, error)
	Parse(userID uuid.UUID, fileName string, r io.Reader) (io.Reader, error)
}

// NewFileReport - builds "FileReport" internal handler.
func NewFileReport(csvUC CsvUC, chatUserUC ChatUserUC, botWrapper *BotWrapper) *FileReport {
	return &FileReport{
		csvUC:      csvUC,
		chatUserUC: chatUserUC,
		BotWrapper: botWrapper,
	}
}

// FileReport - represents an internal handler for processing a CSV report.
type FileReport struct {
	csvUC      CsvUC
	chatUserUC ChatUserUC
	*BotWrapper
}

// Handle - process the CSV MonoBank report, send processed result to the user.
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

	fileURL, err := url.Parse(fileTG.Link(f.bot.Token))
	if err != nil {
		f.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	file, err := f.csvUC.GetFile(fileURL)
	if err != nil {
		f.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	userID, err := f.chatUserUC.GetChatUserID(u.Message.Chat.ID)
	if err != nil {
		f.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	fileResp, err := f.csvUC.Parse(userID, u.Message.Document.FileName, file)
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
