package handlers

import (
	"io"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MappingUC interface {
	Validate(name string) error
	Parse(chatID int64, r io.Reader) error
	GetFile(url string) (io.ReadCloser, error)
}

func NewMapping(uc MappingUC, botWrapper *BotWrapper) *Mapping {
	return &Mapping{uc: uc, BotWrapper: botWrapper}
}

type Mapping struct {
	uc MappingUC
	*BotWrapper
}

func (m *Mapping) Handle(u tg.Update) {
	if err := m.uc.Validate(u.Message.Document.FileName); err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileTG, err := m.bot.GetFile(tg.FileConfig{FileID: u.Message.Document.FileID})
	if err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	file, err := m.uc.GetFile(fileTG.Link(m.bot.Token))
	if err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	if err := m.uc.Parse(u.Message.Chat.ID, file); err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}
	m.sendMSG(tg.NewMessage(u.Message.Chat.ID, "mapping successfully loaded"))
}
