package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/X-AROK/urlcut/internal/app/compress"
	"github.com/X-AROK/urlcut/internal/app/logger"
	"github.com/X-AROK/urlcut/internal/app/url"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MainRouter(ctx context.Context, s url.Repository, baseURL string) chi.Router {
	m := url.NewManager(s)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer, logger.ResponseLogger, logger.RequestLogger, compress.GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Post("/", createShort(ctx, m, baseURL))
		r.Get("/{id}", goToID(ctx, m))
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", shorten(ctx, m, baseURL))
			r.Post("/shorten/batch", batch(ctx, m, baseURL))
		})
	})

	return r
}

func createShort(ctx context.Context, m *url.Manager, baseURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		buff, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u := url.NewURL(string(buff))
		id, err := m.AddURL(ctx, u)
		var existsErr *url.AlreadyExistsError
		if err != nil && !errors.As(err, &existsErr) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		status := http.StatusCreated
		if existsErr != nil {
			id = existsErr.ID
			status = http.StatusConflict
		}

		data := []byte(baseURL + "/" + id)
		w.WriteHeader(status)
		w.Write(data)
	}
}

func goToID(ctx context.Context, m *url.Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		id := chi.URLParam(r, "id")
		url, err := m.GetURL(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, url.OriginalURL, http.StatusTemporaryRedirect)
	}
}

func shorten(ctx context.Context, m *url.Manager, baseURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u := url.NewURL(req.URL)
		_, err := m.AddURL(ctx, u)
		var existsErr *url.AlreadyExistsError
		if err != nil && !errors.As(err, &existsErr) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		status := http.StatusCreated
		if existsErr != nil {
			status = http.StatusConflict
		}

		res := NewShortenResponse(u, baseURL)
		resp, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(resp)
	}
}

func batch(ctx context.Context, m *url.Manager, baseURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		req := BatchRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		urls := req.ToURLs()
		err := m.AddURLsBatch(ctx, urls)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res := NewBatchResponse(urls, baseURL)
		resp, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	}
}
