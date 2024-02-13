package config

import (
	"flag"
	"os"
)

func NewConfigFromFlags() Config {
	var addr string
	flag.StringVar(&addr, "a", ":8080", "Адрес запуска HTTP-сервера")

	var baseURL string
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL")

	flag.Parse()

	builder := NewConfigBuilder()
	builder.WithAddr(addr).WithBaseURL(baseURL)

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		builder.WithAddr(envAddr)
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		builder.WithBaseURL(envBaseURL)
	}

	return builder.Build()
}
