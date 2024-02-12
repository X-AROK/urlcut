package managers

import (
	"errors"
	"io"
	"net/http"

	"github.com/X-AROK/urlcut/internal/app/store"
	"github.com/go-chi/chi/v5"
)

type URLSManager struct {
	s store.Repository
}

func NewURLSManager(s store.Repository) URLSManager {
	return URLSManager{s: s}
}

func (m URLSManager) AddURL(r *http.Request) (string, error) {
	buff, err := io.ReadAll(r.Body)
	if err != nil {
		return "", errors.New("error in data")
	}
	url := string(buff)

	return m.s.Add(url), nil
}

func (m URLSManager) GetURL(r *http.Request) (string, error) {
	id := chi.URLParam(r, "id")
	url, err := m.s.Get(id)
	return url, err
}
