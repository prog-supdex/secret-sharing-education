package filestore

import (
	"log"
	"os"
	"path"
)

type Config struct {
	DataFilePath string `env:"DATA_FILE_PATH, required"`
	RootPath     string
}

func NewConfig() *Config {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{RootPath: currentPath}
}

func (c Config) FullFilePath() string {
	return path.Join(c.RootPath, c.DataFilePath)
}
