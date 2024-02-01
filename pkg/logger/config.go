package logger

type Config struct {
	LogLevel     string
	DisableColor bool
}

func NewConfig() *Config {
	return &Config{}
}
