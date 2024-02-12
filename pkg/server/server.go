package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	port    string
	mux     *mux.Router
	limiter RateLimiter
}

type RoutesMapping map[string]http.HandlerFunc

func New(c Config, r RateLimiter) (*Server, error) {
	return &Server{
		port:    fmt.Sprintf(":%d", c.ServerPort),
		mux:     mux.NewRouter(),
		limiter: r,
	}, nil
}

func (s Server) Run() {
	err := http.ListenAndServe(s.port, s.mux)

	s.mux.Handle("/healthcheck", s.limiter.IpRateLimiter(HealthHandler))

	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (s Server) Mount(routes RoutesMapping) {
	for path, handler := range routes {
		s.mux.Handle(path, s.limiter.IpRateLimiter(handler))
	}
}
