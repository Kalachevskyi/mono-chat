package telegram

import (
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/Kalachevskyi/mono-chat/app/model"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"
)

// TransactionUC - represents a use-case interface for processing business logic of "MonoBank" transactions API.
type TransactionUC interface {
	GetTransactions(token, account string, userID uuid.UUID, from time.Time, to time.Time) (io.Reader, error)
	ParseDate(period string) (from time.Time, to time.Time, err error)
	Locale() *time.Location
}

// NewTransaction - builds "NewTransaction" internal handler.
func NewTransaction(t TokenUC, tr TransactionUC, a AccountUC, cu ChatUserUC, b *BotWrapper) *Transaction {
	return &Transaction{
		tokenUC:       t,
		transactionUC: tr,
		accountUC:     a,
		chatUserUC:    cu,
		BotWrapper:    b,
	}
}

// Transaction - represents an internal handler for processing "MonoBank" transactions API.
type Transaction struct {
	tokenUC       TokenUC
	transactionUC TransactionUC
	accountUC     AccountUC
	chatUserUC    ChatUserUC
	*BotWrapper
}

// Handle - process the "MonoBank" transactions API, send the result to the user.
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
	userID, err := t.chatUserUC.GetChatUserID(chatID)
	if err != nil {
		t.sendDefaultErr(chatID, err)

		return
	}

	token, err := t.tokenUC.Get(userID)
	if err != nil {
		t.sendDefaultErr(chatID, err)

		return
	}

	account, err := t.accountUC.Get(userID)
	if err == model.ErrNil {
		t.sendMSG(tg.NewMessage(chatID, "Please set account."))

		return
	}

	if err != nil {
		t.sendDefaultErr(chatID, err)

		return
	}

	fileResp, err := t.transactionUC.GetTransactions(token, account, userID, from, to)
	if err != nil {
		t.sendDefaultErr(chatID, err)

		return
	}
	reader := tg.FileReader{
		Name:   fmt.Sprintf("%s-%s%s", from.Format(dateTimePattern), to.Format(dateTimePattern), ".csv"),
		Reader: fileResp,
		Size:   -1,
	}

	msg := tg.NewDocumentUpload(chatID, reader)
	t.sendMSG(msg)
}
