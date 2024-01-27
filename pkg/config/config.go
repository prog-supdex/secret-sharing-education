package config

import (
	"github.com/prog-supdex/mini-project/milestone-code/pkg/filestore"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/server"
)

type Config struct {
	Filestore  filestore.Config
	Server     server.Config
	HealthPath string
}

func New() Config {
	config := Config{
		Server:     *server.NewConfig(),
		Filestore:  *filestore.NewConfig(),
		HealthPath: "/healthcheck",
	}

	return config
}
