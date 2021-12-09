package app

import (
	"github.com/gorilla/mux"
	"github.com/mahjadan/login/cmd/handle"
	"net/http"
	"time"
)

type Server struct {
	http.Server
}

func New() Server {
	return Server{
		http.Server{
			Addr:         "127.0.0.1:8080",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) SetupRoutes(router *mux.Router, handler handle.Handler) {
	router.HandleFunc("/login", handler.Login).Methods(http.MethodPost)
	router.HandleFunc("/register", handler.Register).Methods(http.MethodPost)
	s.Handler = router
}
