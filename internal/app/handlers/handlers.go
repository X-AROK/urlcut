package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/X-AROK/urlcut/internal/app/store"
)

func MainHandler(s store.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			goToID(w, r, s)
		} else if r.Method == http.MethodPost {
			createShort(w, r, s)
		} else {
			http.Error(w, "Method not allowed", http.StatusBadRequest)
		}
	}
}

func createShort(w http.ResponseWriter, r *http.Request, s store.Repository) {
	buff, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error in data", http.StatusBadRequest)
		return
	}
	url := string(buff)

	id := s.Add(url)

	data := []byte("http://" + r.Host + "/" + id)
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func goToID(w http.ResponseWriter, r *http.Request, s store.Repository) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	url, ok := s.Get(id)
	if !ok {
		http.Error(w, "Id not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
