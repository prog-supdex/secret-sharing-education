package filestore

import (
	"log/slog"
	"os"
	"path"
)

type Config struct {
	DataFilePath string
	rootPath     string
}

func NewConfig() *Config {
	currentPath, err := os.Getwd()
	if err != nil {
		slog.Error("Getting a rooted path name was failed: " + err.Error())
		os.Exit(1)
	}

	return &Config{rootPath: currentPath}
}

func (c Config) FullFilePath() string {
	return path.Join(c.rootPath, c.DataFilePath)
}
