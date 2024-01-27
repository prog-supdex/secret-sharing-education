package server

type Config struct {
	ServerPort string `env:"LISTEN_PORT, default=:8080"`
}

func NewConfig() *Config {
	return &Config{}
}
