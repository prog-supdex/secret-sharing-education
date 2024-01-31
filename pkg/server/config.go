package server

type Config struct {
	ServerPort int
}

func NewConfig() *Config {
	return &Config{}
}
