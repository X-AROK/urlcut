package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		uri := r.RequestURI

		start := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(start)

		Log.Info(
			"Request",
			zap.String("method", method),
			zap.String("uri", uri),
			zap.String("duration", duration.String()),
		)
	})
}

func ResponseLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseData := &responseData{}
		wp := newLoggingResponseWriter(w, responseData)
		h.ServeHTTP(&wp, r)

		Log.Info(
			"Response",
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	})
}
