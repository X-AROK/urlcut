package main

import (
	"net/http"

	"github.com/X-AROK/urlcut/internal/app/config"
	"github.com/X-AROK/urlcut/internal/app/handlers"
	"github.com/X-AROK/urlcut/internal/app/logger"
	"github.com/X-AROK/urlcut/internal/app/store"
	"go.uber.org/zap"
)

func main() {
	c := config.NewConfigFromFlags()
	logger.Initialize(c.LoggerLevel)

	logger.Log.Info(
		"Starting server",
		zap.String("addr", c.Addr),
	)
	if err := run(c.Addr, c.BaseURL); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func run(addr, baseURL string) error {
	s := store.NewMapStore()
	r := handlers.MainRouter(&s, baseURL)

	return http.ListenAndServe(addr, r)
}
