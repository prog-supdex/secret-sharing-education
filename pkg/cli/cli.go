package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/config"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets/handlers"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/server"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/version"
	"github.com/sethvargo/go-envconfig"
	"net/http"
	"os"
)

var showVersion bool

func init() {
	flag.BoolVar(&showVersion, "v", false, "Show the project version")
}

func Run() error {
	flag.Parse()
	if showVersion {
		fmt.Print(version.Version())
		os.Exit(0)
	}

	cfg, err := initConfig()
	if err != nil {
		return err
	}

	if cfg.Filestore.DataFilePath == "" {
		return errors.New("ENV DATA_FILE_PATH is blank")
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

func initConfig() (*config.Config, error) {
	ctx := context.Background()
	cfg := config.New()
	err := envconfig.Process(ctx, &cfg)

	return &cfg, err
}
