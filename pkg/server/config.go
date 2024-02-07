package server

type Config struct {
	ServerPort int
	Request    ConfigRequest
}

type ConfigRequest struct {
	RequestsLimit int
	Within        int
}

func NewConfig() *Config {
	return &Config{}
}
