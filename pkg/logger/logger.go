package logger

import (
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

func InitLogger(config Config, output io.Writer) {
	var programLevel = new(slog.LevelVar)

	switch strings.ToUpper(config.LogLevel) {
	case "INFO":
		{
			programLevel.Set(slog.LevelInfo)
		}
	case "DEBUG":
		{
			programLevel.Set(slog.LevelDebug)
		}
	case "WARN":
		{
			programLevel.Set(slog.LevelWarn)
		}
	case "ERROR":
		{
			programLevel.Set(slog.LevelError)
		}
	}

	handler := tint.NewHandler(output, &tint.Options{
		Level:      programLevel,
		TimeFormat: time.Kitchen,
		NoColor:    config.DisableColor || !isatty.IsTerminal(os.Stdout.Fd()),
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
