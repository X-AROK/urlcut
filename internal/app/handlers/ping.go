package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

func PingDB(ctx context.Context, db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		if err := db.PingContext(ctxTimeout); err != nil {
			http.Error(w, "ошибка соединения с бд", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
