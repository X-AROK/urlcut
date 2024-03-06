package store

import "github.com/X-AROK/urlcut/internal/app/url"

type MockStore struct{}

func NewMockStore() *MockStore {
	return &MockStore{}
}

func (s MockStore) Get(id string) (*url.URL, error) {
	if id == "test" {
		return &url.URL{ShortURL: "test", OriginalURL: "https://practicum.yandex.ru"}, nil
	}

	return &url.URL{}, url.ErrorNotFound
}

func (s MockStore) Add(u *url.URL) (string, error) {
	u.ShortURL = "test"
	return "test", nil
}
