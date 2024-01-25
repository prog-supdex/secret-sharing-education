package config

import (
	"log"
	"os"
	"path"
)

type Config struct {
	ServerPort   string `env:"LISTEN_PORT, default=:8080"`
	DataFilePath string `env:"DATA_FILE_PATH, required"`
	RootPath     string
}

func New() *Config {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{RootPath: currentPath}
}

func (c Config) FullFilePath() string {
	return path.Join(c.RootPath, c.DataFilePath)
}
