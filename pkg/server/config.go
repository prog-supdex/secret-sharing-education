package server

type Config struct {
	ServerPort              int
	RequestsLimit           int
	Within                  int
	IpBucketLifeTimeSeconds int
}

func NewConfig() *Config {
	return &Config{}
}
