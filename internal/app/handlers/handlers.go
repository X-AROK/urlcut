package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/X-AROK/urlcut/internal/app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MainRouter(s store.Repository, baseURL string) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Route("/", func(r chi.Router) {
		r.Post("/", createShort(s, baseURL))
		r.Get("/{id}", goToID(s))
	})

	return r
}

func createShort(s store.Repository, baseURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		buff, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error in data", http.StatusBadRequest)
			return
		}
		url := string(buff)

		id := s.Add(url)

		data := []byte(baseURL + "/" + id)
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

func goToID(s store.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/")
		url, ok := s.Get(id)
		if !ok {
			http.Error(w, "Id not found", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
