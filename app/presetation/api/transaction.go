package api

import (
	"fmt"
	"io"
	"time"

	"github.com/Kalachevskyi/mono-chat/app/model"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"
)

// TransactionUC - represents a use-case interface for processing business logic of "MonoBank" transactions API
type TransactionUC interface {
	GetTransactions(token, account string, chatID int64, from time.Time, to time.Time) (io.Reader, error)
	ParseDate(period string) (from time.Time, to time.Time, err error)
	Locale() *time.Location
}

// NewTransaction - builds "NewTransaction" internal handler
func NewTransaction(t TokenUC, a AccountUC, tr TransactionUC, b *BotWrapper) *Transaction {
	return &Transaction{
		tokenUC:       t,
		transactionUC: tr,
		BotWrapper:    b,
		accountUC:     a,
	}
}

// Transaction - represents an internal handler for processing "MonoBank" transactions API
type Transaction struct {
	tokenUC       TokenUC
	transactionUC TransactionUC
	accountUC     AccountUC
	*BotWrapper
}

// Handle - process the "MonoBank" transactions API, send the result to the user
func (t *Transaction) Handle(u tg.Update) {
	var (
		from, to time.Time
		timeNow  = now.New(time.Now().In(t.transactionUC.Locale()))
	)

	switch u.Message.Command() {
	case getCommand:
		var err error
		from, to, err = t.transactionUC.ParseDate(u.Message.CommandArguments())
		if err != nil {
			t.sendDefaultErr(u.Message.Chat.ID, err)
			return
		}
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

	account, err := t.accountUC.Get(chatID)
	if err == model.ErrNil {
		t.sendMSG(tg.NewMessage(chatID, "Please set account."))
		return
	}

	if err != nil {
		t.sendDefaultErr(u.Message.Chat.ID, err)
		return
	}

	fileResp, err := t.transactionUC.GetTransactions(token, account, chatID, from, to)
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
