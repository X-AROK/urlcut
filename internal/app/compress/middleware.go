package compress

import (
	"net/http"
	"strings"
)

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			zw := newGzipWriter(w)
			ow = zw

			defer zw.Close()
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			zr, err := newGzipReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			r.Body = zr
			defer zr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}
