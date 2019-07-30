package handlers

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/now"
)

type ApiUC interface {
	GetTransactions(token string, chatID int64, from time.Time, to time.Time) (io.Reader, error)
	ParseDate(period string) (from time.Time, to time.Time, err error)
	Locale() *time.Location
}

func NewTransaction(tokenUC TokenUC, apiUC ApiUC, botWrapper *BotWrapper) *Transaction {
	return &Transaction{tokenUC: tokenUC, apiUC: apiUC, BotWrapper: botWrapper}
}

type Transaction struct {
	tokenUC TokenUC
	apiUC   ApiUC
	*BotWrapper
}

func (t *Transaction) Handle(u tg.Update) {
	var (
		from, to time.Time
		timeNow  = now.New(time.Now().In(t.apiUC.Locale()))
	)
	switch u.Message.Command() {
	case getCommand:
		fromTime, toTime, err := t.apiUC.ParseDate(u.Message.CommandArguments())
		if err != nil {
			t.sendDefaultErr(u.Message.Chat.ID, err)
			return
		}
		from, to = fromTime, toTime
	case todayCommand:
		from, to = timeNow.BeginningOfDay(), timeNow.EndOfDay()
	case currentMonthCommand:
		from, to = timeNow.BeginningOfMonth(), timeNow.EndOfMonth()
	default:
		t.sendDefaultErr(u.Message.Chat.ID, errors.New("can't detect command"))
		return
	}

	chatID := u.Message.Chat.ID
	token, err := t.tokenUC.Get(chatID)
	if err != nil {
		t.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileResp, err := t.apiUC.GetTransactions(token, chatID, from, to)
	if err != nil {
		t.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}
	reader := tg.FileReader{
		Name:   fmt.Sprintf("%s-%s%s", from.Format(dateTimePattern), to.Format(dateTimePattern), ".csv"),
		Reader: fileResp,
		Size:   -1,
	}
	msg := tg.NewDocumentUpload(u.Message.Chat.ID, reader)
	t.sendMSG(msg)
}
