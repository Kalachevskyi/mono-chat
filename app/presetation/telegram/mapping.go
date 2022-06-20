package telegram

import (
	"io"
	"net/url"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

// MappingUC - represents a usecase interface for processing category mapping business logic.
type MappingUC interface {
	Validate(name string) error
	Parse(userID uuid.UUID, r io.Reader) error
	GetFile(u *url.URL) (io.ReadCloser, error)
}

// NewMapping - builds "NewMapping" internal handler.
func NewMapping(mappingUC MappingUC, chatUserUC ChatUserUC, botWrapper *BotWrapper) *Mapping {
	return &Mapping{
		mappingUC:  mappingUC,
		chatUserUC: chatUserUC,
		BotWrapper: botWrapper,
	}
}

// Mapping - represents an internal handler for processing category mapping.
type Mapping struct {
	mappingUC  MappingUC
	chatUserUC ChatUserUC
	*BotWrapper
}

// Handle - process category mapping, send the result to the user.
func (m *Mapping) Handle(u tg.Update) {
	if err := m.mappingUC.Validate(u.Message.Document.FileName); err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	fileTG, err := m.bot.GetFile(tg.FileConfig{FileID: u.Message.Document.FileID})
	if err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	fileURL, err := url.Parse(fileTG.Link(m.bot.Token))
	if err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	file, err := m.mappingUC.GetFile(fileURL)
	if err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	userID, err := m.chatUserUC.GetChatUserID(u.Message.Chat.ID)
	if err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}

	if err := m.mappingUC.Parse(userID, file); err != nil {
		m.sendDefaultErr(u.Message.Chat.ID, err)

		return
	}
	m.sendMSG(tg.NewMessage(u.Message.Chat.ID, "mapping successfully loaded"))
}
