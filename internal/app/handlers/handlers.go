package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/X-AROK/urlcut/internal/app/compress"
	"github.com/X-AROK/urlcut/internal/app/logger"
	"github.com/X-AROK/urlcut/internal/app/url"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MainRouter(s url.Repository, baseURL string) chi.Router {
	m := url.NewManager(s)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer, logger.ResponseLogger, logger.RequestLogger, compress.GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Post("/", createShort(m, baseURL))
		r.Get("/{id}", goToID(m))
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", shorten(m, baseURL))
		})
	})

	return r
}

func createShort(m *url.Manager, baseURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		buff, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u := url.NewURL(string(buff))
		id, err := m.AddURL(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := []byte(baseURL + "/" + id)
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

func goToID(m *url.Manager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		id := chi.URLParam(r, "id")
		url, err := m.GetURL(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, url.Addr, http.StatusTemporaryRedirect)
	}
}

func shorten(m *url.Manager, baseURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var u url.URL
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := m.AddURL(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res := url.NewResult(baseURL + "/" + id)
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
