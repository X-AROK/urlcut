package config

type Config struct {
	Addr    string
	BaseURL string
}

func New() Config {
	return Config{}
}
