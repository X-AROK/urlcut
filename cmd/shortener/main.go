package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/X-AROK/urlcut/internal/app/config"
	"github.com/X-AROK/urlcut/internal/app/handlers"
	"github.com/X-AROK/urlcut/internal/app/logger"
	"github.com/X-AROK/urlcut/internal/app/store"
	"github.com/X-AROK/urlcut/internal/app/url"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	c := config.NewConfigFromFlags()
	logger.Initialize(c.LoggerLevel)

	logger.Log.Info(
		"Starting server",
		zap.String("addr", c.Addr),
	)
	if err := run(c.Addr, c.BaseURL, c.FileStorageFile, c.DSN); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Fatal("Server starting error", zap.Error(err))
	}
}

func run(addr, baseURL, fileStoragePath, dsn string) error {
	var s url.Repository
	var db *sql.DB

	if dsn != "" {
		_db, err := sql.Open("pgx", dsn)
		if err != nil {
			return fmt.Errorf("sql open error: %w", err)
		}
		defer _db.Close()

		dbs, err := store.NewDBStore(_db)
		if err != nil {
			return fmt.Errorf("create db store error: %w", err)
		}

		s = dbs
		db = _db
	} else if fileStoragePath != "" {
		fs, err := store.NewFileStore(fileStoragePath)
		if err != nil {
			return fmt.Errorf("create file store error: %w", err)
		}
		defer fs.Close()
		s = fs
	} else {
		s = store.NewMapStore()
	}

	ctx := context.Background()
	r := handlers.MainRouter(ctx, s, baseURL)

	if db != nil {
		r.Get("/ping", handlers.PingDB(ctx, db))
	}

	return http.ListenAndServe(addr, r)
}
