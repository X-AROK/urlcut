package store

import "github.com/X-AROK/urlcut/internal/app/url"

type MockStore struct{}

func NewMockStore() *MockStore {
	return &MockStore{}
}

func (s MockStore) Get(id string) (url.URL, error) {
	if id == "test" {
		return url.NewURL("https://practicum.yandex.ru"), nil
	}

	return url.URL{}, url.ErrorNotFound
}

func (s MockStore) Add(_ url.URL) (string, error) {
	return "test", nil
}
