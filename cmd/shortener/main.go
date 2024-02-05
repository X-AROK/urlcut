package main

import (
	"net/http"

	"github.com/X-AROK/urlcut/internal/app/handlers"
	"github.com/X-AROK/urlcut/internal/app/store"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.MainHandler(store.NewMapStore()))

	return http.ListenAndServe(":8080", mux)
}
