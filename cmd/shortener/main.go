package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/X-AROK/urlcut/internal/app/config"
	"github.com/X-AROK/urlcut/internal/app/handlers"
	"github.com/X-AROK/urlcut/internal/app/logger"
	"github.com/X-AROK/urlcut/internal/app/store"
	"github.com/X-AROK/urlcut/internal/app/url"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	c := config.NewConfigFromFlags()
	logger.Initialize(c.LoggerLevel)

	logger.Log.Info(
		"Starting server",
		zap.String("addr", c.Addr),
	)
	if err := run(c.Addr, c.BaseURL, c.FileStorageFile, c.DSN); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatal("Server starting error", zap.Error(err))
	}
}

func run(addr, baseURL, fileStoragePath, dsn string) error {
	var s url.Repository
	var db *sql.DB

	if dsn != "" {
		_db, err := sql.Open("pgx", dsn)
		if err != nil {
			return err
		}
		db = _db
	}
	if fileStoragePath == "" {
		s = store.NewMapStore()
	} else {
		fs, err := store.NewFileStore(fileStoragePath)
		if err != nil {
			return err
		}
		defer fs.Close()
		s = fs
	}

	r := handlers.MainRouter(s, baseURL)

	if db != nil {
		r.Get("/ping", handlers.PingDB(context.Background(), db))
	}

	return http.ListenAndServe(addr, r)
}
