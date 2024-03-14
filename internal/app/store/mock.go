package store

import (
	"context"
	"fmt"
	"github.com/X-AROK/urlcut/internal/app/url"
)

type MockStore struct{}

func NewMockStore() *MockStore {
	return &MockStore{}
}

func (s MockStore) Get(ctx context.Context, id string) (*url.URL, error) {
	if id == "test" {
		return &url.URL{ShortURL: "test", OriginalURL: "https://practicum.yandex.ru"}, nil
	}

	return &url.URL{}, url.ErrorNotFound
}

func (s MockStore) Add(ctx context.Context, u *url.URL) (string, error) {
	u.ShortURL = "test"
	return "test", nil
}

func (s MockStore) AddBatch(ctx context.Context, urls *url.URLsBatch) error {
	for k, v := range *urls {
		id := fmt.Sprintf("test%s", k)
		v.ShortURL = id
	}

	return nil
}
