package main

import (
	"net/http"

	"github.com/X-AROK/urlcut/internal/app/config"
	"github.com/X-AROK/urlcut/internal/app/handlers"
	"github.com/X-AROK/urlcut/internal/app/store"
)

func main() {
	config.ParseFlags()

	if err := run(config.Addr, config.BaseURL); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func run(addr, baseURL string) error {
	s := store.NewMapStore()
	r := handlers.MainRouter(s, baseURL)

	return http.ListenAndServe(addr, r)
}
