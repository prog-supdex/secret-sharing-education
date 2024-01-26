package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	port string
	mux  *mux.Router
}

type RoutesMapping map[string]http.Handler

func New(c Config) (*Server, error) {
	return &Server{
		port: c.ServerPort,
		mux:  mux.NewRouter(),
	}, nil
}

func (s Server) Run() {
	err := http.ListenAndServe(s.port, s.mux)

	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (s Server) Mount(routes RoutesMapping) {
	for path, handler := range routes {
		s.mux.Handle(path, handler)
	}
}
