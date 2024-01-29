package logger

import (
	"log/slog"
	"strings"
)

func InitLogger(level string) {
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

	logger := slog.New(NewHandler(&slog.HandlerOptions{
		Level: programLevel,
	}))

	slog.SetDefault(logger)
}
