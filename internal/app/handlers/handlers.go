package handlers

import (
	"net/http"

	"github.com/X-AROK/urlcut/internal/app/managers"
	"github.com/X-AROK/urlcut/internal/app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MainRouter(s store.Repository, baseURL string) chi.Router {
	m := managers.NewURLSManager(s)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Route("/", func(r chi.Router) {
		r.Post("/", createShort(m, baseURL))
		r.Get("/{id}", goToID(m))
	})

	return r
}

func createShort(m managers.URLSManager, baseURL string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		id, err := m.AddURL(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := []byte(baseURL + "/" + id)
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

func goToID(m managers.URLSManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := m.GetURL(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
