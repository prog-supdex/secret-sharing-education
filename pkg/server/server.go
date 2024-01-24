package server

import (
	"github.com/gorilla/mux"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/config"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets/handlers"
	"log"
	"net/http"
)

type Server struct {
	port          string
	secretManager secrets.Manager
}

func New(c config.Config) (*Server, error) {
	fileStore, err := filestore.New(c.FullFilePath())

	if err != nil {
		return nil, err
	}

	secretManager := secrets.NewSecretManager(fileStore)

	return &Server{
		port:          c.ServerPort,
		secretManager: secretManager,
	}, nil
}

func (s Server) Run() {
	r := mux.NewRouter()

	handlers.NewSecretHandler(s.secretManager).RegisterHandler(r)

	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	err := http.ListenAndServe(s.port, r)

	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
