package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/config"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/server"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/version"
	"github.com/sethvargo/go-envconfig"
	"os"
)

var showVersion bool

func init() {
	flag.BoolVar(&showVersion, "v", false, "Show the project version")
}

func Run() error {
	flag.Parse()
	if showVersion {
		fmt.Printf("Version: %s", version.Version())
		os.Exit(0)
	}

	cfg, err := initConfig()
	if err != nil {
		return err
	}

	if cfg.DataFilePath == "" {
		return errors.New("ENV DATA_FILE_PATH is blank")
	}

	srv, err := server.New(*cfg)
	if err != nil {
		return err
	}

	srv.Run()

	return nil
}

func initConfig() (*config.Config, error) {
	ctx := context.Background()
	cfg := config.New()
	err := envconfig.Process(ctx, cfg)

	return cfg, err
}
