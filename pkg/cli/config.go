package cli

import (
	"fmt"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/config"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/version"
	"github.com/urfave/cli/v2"
	"os"
)

func NewCliConfig() (cfg *config.Config, err error, stopProgram bool) {
	cfg = config.New()
	helpOrVersionWereShown := true

	cli.VersionPrinter = func(cCtx *cli.Context) {
		_, _ = fmt.Fprintf(cCtx.App.Writer, "%v\n", cCtx.App.Version)
	}

	app := &cli.App{
		Name:            "secrets-share",
		Version:         version.Version(),
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file_path",
				Usage:       "set the file_path which will store the secrets",
				EnvVars:     []string{"DATA_FILE_PATH"},
				Category:    "Filestore:",
				Destination: &cfg.Filestore.DataFilePath,
				Required:    true,
			},
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "set the file_path which will store the secrets",
				Value:       8080,
				EnvVars:     []string{"LISTEN_PORT"},
				Category:    "Server:",
				Destination: &cfg.Server.ServerPort,
				Required:    true,
			},
		},
		Action: func(nc *cli.Context) error {
			helpOrVersionWereShown = false
			return nil
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		return &config.Config{}, err, false
	}

	if helpOrVersionWereShown {
		return &config.Config{}, nil, true
	}

	return cfg, nil, false
}
