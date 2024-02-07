package cli

import (
	"encoding/json"
	"flag"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/logger"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/secrets/handlers"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/server"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/version"
	"log/slog"
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

	logger.InitLogger(cfg.Logger, os.Stdout)

	fileStore, err := filestore.New(cfg.Filestore.FullFilePath())
	if err != nil {
		slog.Error("FileStore initialize error: " + err.Error())
		return err
	}

	rateLimiter := server.NewRateLimit(cfg.Server.Request)

	srv, err := server.New(cfg.Server, rateLimiter)
	if err != nil {
		slog.Error("Server initialize error: " + err.Error())
		return err
	}

	secretManager := secrets.NewSecretManager(fileStore)
	secretHandler := handlers.NewSecretHandler(secretManager)

	routes := secretHandler.Routes()

	routesKeys := make([]string, 0, len(routes))
	for k := range routes {
		routesKeys = append(routesKeys, k)
	}

	urlsJson, _ := json.Marshal(routesKeys)

	slog.Info("Starting application",
		"version", version.Version(),
		"serverPort", cfg.Server.ServerPort,
		"endpoints", string(urlsJson))

	srv.Mount(routes)
	srv.Run()

	return nil
}
