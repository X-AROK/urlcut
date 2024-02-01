package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/X-AROK/urlcut/internal/app/store"
	"github.com/X-AROK/urlcut/internal/util"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		goToId(w, r)
	} else if r.Method == http.MethodPost {
		createShort(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}
}

func createShort(w http.ResponseWriter, r *http.Request) {
	buff, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error in data", http.StatusBadRequest)
		return
	}
	url := string(buff)

	id := util.GenerateId(8)
	store.Set(id, url)

	data := []byte("http://" + r.Host + "/" + id)
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func goToId(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	url, ok := store.Get(id)
	if !ok {
		http.Error(w, "Id not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
