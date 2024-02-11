package config

import (
	"flag"
	"os"
)

func ParseFlags(c *Config) {
	flag.StringVar(&c.Addr, "a", ":8080", "Адрес запуска HTTP-сервера")
	flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL")

	flag.Parse()

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		c.Addr = envAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		c.BaseURL = envBaseURL
	}
}
