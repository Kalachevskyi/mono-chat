package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewService constructor for HTTP service.
func NewService(transactionHandler *Transaction, port int) *Service {
	s := Service{
		transactionHandler: transactionHandler,
		router:             mux.NewRouter(),
		port:               port,
	}

	s.injectRoutes()

	return &s
}

// Service represents http service.
type Service struct {
	port               int
	transactionHandler *Transaction
	router             *mux.Router
}

func (s *Service) injectRoutes() {
	s.router.HandleFunc("/transactions/month", s.transactionHandler.GetCurrentMonth)
	s.router.HandleFunc("/transactions/today", s.transactionHandler.GetCurrentDay)
	s.router.HandleFunc("/transactions/{from}/{to}", s.transactionHandler.Get)
}

// Start HTTP service.
func (s Service) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router)
}
