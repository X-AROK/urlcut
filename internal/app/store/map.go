package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/X-AROK/urlcut/internal/app/url"
	"github.com/X-AROK/urlcut/internal/util"
)

type MapStore struct {
	mx     sync.Mutex
	values map[string]*url.URL
}

func NewMapStore() *MapStore {
	return &MapStore{values: make(map[string]*url.URL)}
}

func (s *MapStore) Get(ctx context.Context, id string) (*url.URL, error) {
	s.mx.Lock()
	v, ok := s.values[id]
	s.mx.Unlock()

	if !ok {
		return v, url.ErrNotFound
	}
	return v, nil
}

func (s *MapStore) Add(ctx context.Context, v *url.URL) (string, error) {
	id := util.GenerateID(8)
	s.mx.Lock()
	v.ShortURL = id
	s.values[id] = v
	s.mx.Unlock()
	return id, nil
}

func (s *MapStore) AddBatch(ctx context.Context, urls *url.URLsBatch) error {
	for _, u := range *urls {
		_, err := s.Add(ctx, u)
		if err != nil {
			return fmt.Errorf("add to map error: %w", err)
		}
	}

	return nil
}
