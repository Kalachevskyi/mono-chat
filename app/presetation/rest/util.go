package rest

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const dateTimePattern = "02.01.2006T15.04"

func sendUserUnauthorizedError(w http.ResponseWriter, log Logger) {
	http.Error(w, "user unauthorized", http.StatusUnauthorized)
	log.Error("user unauthorized")
}

func sendWrongUUIDError(w http.ResponseWriter, log Logger, user string) {
	http.Error(w, "can't parse user to UUID", http.StatusBadRequest)
	log.Errorf("can't parse user to UUID: %s", user)
}

func sendServerError(w http.ResponseWriter, log Logger, msg string) {
	http.Error(w, "internal server error", http.StatusInternalServerError)
	log.Error(msg)
}

func sendCantFindUserError(w http.ResponseWriter, log Logger, user uuid.UUID) {
	http.Error(w, "can't find user", http.StatusUnauthorized)
	log.Errorf("can't find user: %v", user)
}

func sendBadRequestError(w http.ResponseWriter, log Logger, msg string) {
	http.Error(w, fmt.Sprintf("bad request: %s", msg), http.StatusBadRequest)
	log.Errorf("bad request: %s", msg)
}
