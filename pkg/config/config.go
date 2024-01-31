package config

import (
	"github.com/prog-supdex/mini-project/milestone-code/pkg/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/logger"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/server"
)

type Config struct {
	Filestore filestore.Config
	Server    server.Config
	Logger    logger.Config
}

func New() *Config {
	config := Config{
		Server:    *server.NewConfig(),
		Filestore: *filestore.NewConfig(),
		Logger:    *logger.NewConfig(),
	}

	return &config
}
