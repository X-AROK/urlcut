package handlers

import (
	"io"
	"net/http"

	"github.com/X-AROK/urlcut/internal/app/logger"
	"github.com/X-AROK/urlcut/internal/app/url"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MainRouter(s url.Repository, baseURL string) chi.Router {
	m := url.NewManager(s)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer, logger.ResponseLogger, logger.RequestLogger)
	r.Route("/", func(r chi.Router) {
		r.Post("/", createShort(m, baseURL))
		r.Get("/{id}", goToID(m))
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
