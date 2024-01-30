package logger

import (
	"io"
	"log/slog"
	"strings"
)

func InitLogger(level string, output io.Writer) {
	var programLevel = new(slog.LevelVar)

	switch strings.ToUpper(level) {
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

	opts := PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: programLevel,
		},
	}
	handler := NewPrettyHandler(output, opts)

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
