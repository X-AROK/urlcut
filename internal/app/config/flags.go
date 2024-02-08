package config

import "flag"

var Addr string
var BaseURL string

func ParseFlags() {
	flag.StringVar(&Addr, "a", ":8080", "Адрес запуска HTTP-сервера")
	flag.StringVar(&BaseURL, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL")

	flag.Parse()
}
