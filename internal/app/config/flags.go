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

	var level string
	flag.StringVar(&level, "l", "info", "Уроввень логироввания")

	var fileStoragePath string
	flag.StringVar(&fileStoragePath, "f", "/tmp/short-url-db.json", "Файл хранилища")

	var dsn string
	flag.StringVar(&dsn, "d", "", "Строка подключения БД")

	flag.Parse()

	builder := NewConfigBuilder()
	builder.WithAddr(addr).WithBaseURL(baseURL).WithLoggerLevel(level).WithFileStorage(fileStoragePath).WithDSN(dsn)

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		builder.WithAddr(envAddr)
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		builder.WithBaseURL(envBaseURL)
	}
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		builder.WithLoggerLevel(envLevel)
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		builder.WithFileStorage(envFileStoragePath)
	}
	if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
		builder.WithDSN(envDSN)
	}

	return builder.Build()
}
