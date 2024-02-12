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
				Value:       "data.json",
				Category:    "Filestore:",
				Destination: &cfg.Filestore.DataFilePath,
			},
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "set the file_path which will store the secrets",
				Value:       8080,
				EnvVars:     []string{"PORT"},
				Category:    "Server:",
				Destination: &cfg.Server.ServerPort,
			},
			&cli.StringFlag{
				Name:        "log_level",
				Usage:       "set the file_path which will store the secrets",
				EnvVars:     []string{"LOG_LEVEL"},
				Value:       "INFO",
				Category:    "Logger:",
				Destination: &cfg.Logger.LogLevel,
			},
			&cli.BoolFlag{
				Name:        "no_color",
				Usage:       "set the file_path which will store the secrets",
				EnvVars:     []string{"NO_COLOR"},
				Value:       false,
				Category:    "Logger:",
				Destination: &cfg.Logger.DisableColor,
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
