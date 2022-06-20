package rest

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jinzhu/now"

	"github.com/Kalachevskyi/mono-chat/app/model"
)

const (
	dateFormat          = "2006-01-02"
	authorizationHeader = "Authorization"
)

// HTTP path keys.
const (
	fromKey = "from"
	toKey   = "to"
)

// NewTransaction constructor for Transaction.
func NewTransaction(log Logger, transactionUC TransactionUC, userUC UserUC, accountUC AccountUC, tokenUC TokenUC) *Transaction {
	return &Transaction{
		log:           log,
		transactionUC: transactionUC,
		userUC:        userUC,
		accountUC:     accountUC,
		tokenUC:       tokenUC,
	}
}

// Transaction represents transaction REST handler.
type Transaction struct {
	log           Logger
	transactionUC TransactionUC
	userUC        UserUC
	accountUC     AccountUC
	tokenUC       TokenUC
}

// GetCurrentMonth returns current month transactions.
func (t Transaction) GetCurrentMonth(w http.ResponseWriter, r *http.Request) {
	timeNow := now.New(time.Now().In(t.transactionUC.Locale()))
	from, to := timeNow.BeginningOfMonth(), timeNow.EndOfMonth()
	t.handleTransactions(w, r, from, to)
}

// GetCurrentDay returns current day transactions.
func (t Transaction) GetCurrentDay(w http.ResponseWriter, r *http.Request) {
	timeNow := now.New(time.Now().In(t.transactionUC.Locale()))
	from, to := timeNow.BeginningOfDay(), timeNow.EndOfDay()
	t.handleTransactions(w, r, from, to)
}

// Get - get by date range from, to.
func (t Transaction) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fromRaw, ok := vars[fromKey]
	if !ok {
		sendBadRequestError(w, t.log, "can't get From parameter")

		return
	}

	toRaw, ok := vars[toKey]
	if !ok {
		sendBadRequestError(w, t.log, "can't get To parameter")

		return
	}

	from, err := time.Parse(dateFormat, fromRaw)
	if err != nil {
		sendBadRequestError(w, t.log, "can't parse From parameter")

		return
	}

	to, err := time.Parse(dateFormat, toRaw)
	if err != nil {
		sendBadRequestError(w, t.log, "can't parse To parameter")

		return
	}
	to = now.New(to).EndOfDay()

	t.handleTransactions(w, r, from, to)
}

func (t Transaction) handleTransactions(w http.ResponseWriter, r *http.Request, from, to time.Time) {
	userRaw := r.Header.Get(authorizationHeader)
	if userRaw == "" {
		sendUserUnauthorizedError(w, t.log)

		return
	}

	userID, err := uuid.Parse(userRaw)
	if err != nil {
		sendWrongUUIDError(w, t.log, userRaw)

		return
	}

	ok, err := t.userUC.CheckUser(userID)
	if err != nil {
		sendServerError(w, t.log, err.Error())

		return
	}

	if !ok {
		sendCantFindUserError(w, t.log, userID)

		return
	}

	account, err := t.accountUC.Get(userID)
	if err == model.ErrNil {
		sendCantFindUserError(w, t.log, userID)

		return
	}

	if err != nil {
		sendServerError(w, t.log, err.Error())

		return
	}

	token, err := t.tokenUC.Get(userID)
	if err != nil {
		sendServerError(w, t.log, err.Error())

		return
	}
	fileResp, err := t.transactionUC.GetTransactions(token, account, userID, from, to)
	if err != nil {
		sendServerError(w, t.log, err.Error())

		return
	}

	contentType := fmt.Sprintf("attachment;filename=%s-%s%s", from.Format(dateTimePattern), to.Format(dateTimePattern), ".csv")
	w.Header().Set("Content-Disposition", contentType)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	if _, err = io.Copy(w, fileResp); err != nil {
		t.log.Error(err)
	}
}
