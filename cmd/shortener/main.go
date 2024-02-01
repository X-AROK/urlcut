package main

import (
	"net/http"

	"github.com/X-AROK/urlcut/internal/app/handlers"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.MainHandler)

	return http.ListenAndServe(":8080", mux)
}
