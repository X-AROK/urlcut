package store

import (
	"sync"

	"github.com/X-AROK/urlcut/internal/app/url"
	"github.com/X-AROK/urlcut/internal/util"
)

type MapStore struct {
	mx     sync.Mutex
	values map[string]url.URL
}

func NewMapStore() *MapStore {
	return &MapStore{values: make(map[string]url.URL)}
}

func (s *MapStore) Get(id string) (url.URL, error) {
	s.mx.Lock()
	v, ok := s.values[id]
	s.mx.Unlock()

	if !ok {
		return v, url.ErrorNotFound
	}
	return v, nil
}

func (s *MapStore) Add(v url.URL) (string, error) {
	id := util.GenerateID(8)
	s.mx.Lock()
	s.values[id] = v
	s.mx.Unlock()
	return id, nil
}
