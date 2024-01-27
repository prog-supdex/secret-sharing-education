package cli

import (
	"flag"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets/handlers"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/server"
	"net/http"
	"os"
)

var showVersion bool

func init() {
	flag.BoolVar(&showVersion, "v", false, "Show the project version")
}

func Run() error {
	cfg, err, stopProgram := NewCliConfig()
	if err != nil {
		return err
	}

	if stopProgram {
		os.Exit(0)
	}

	fileStore, err := filestore.New(cfg.Filestore.FullFilePath())
	if err != nil {
		return err
	}

	srv, err := server.New(cfg.Server)
	if err != nil {
		return err
	}

	secretManager := secrets.NewSecretManager(fileStore)
	secretHandler := handlers.NewSecretHandler(secretManager)

	routes := secretHandler.Routes()
	// Add HealthHandler to map
	routes[cfg.HealthPath] = http.HandlerFunc(server.HealthHandler)

	srv.Mount(routes)

	srv.Run()

	return nil
}
