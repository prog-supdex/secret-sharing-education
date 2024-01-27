package filestore

import (
	"log"
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
		log.Fatal(err)
	}

	return &Config{rootPath: currentPath}
}

func (c Config) FullFilePath() string {
	return path.Join(c.rootPath, c.DataFilePath)
}
