package config

import (
	"github.com/prog-supdex/mini-project/milestone-code/pkg/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/server"
)

type Config struct {
	Filestore filestore.Config
	Server    server.Config
	LogLevel  string
	LogKind   string
}

func New() *Config {
	config := Config{
		Server:    *server.NewConfig(),
		Filestore: *filestore.NewConfig(),
		LogLevel:  "INFO",
		LogKind:   "Text",
	}

	return &config
}
