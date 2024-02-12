package store

import (
	"sync"

	"github.com/X-AROK/urlcut/internal/util"
)

type MapStore struct {
	mx     sync.Mutex
	values map[string]string
}

func NewMapStore() MapStore {
	return MapStore{values: make(map[string]string)}
}

func (s *MapStore) Get(id string) (string, error) {
	s.mx.Lock()
	v, ok := s.values[id]
	s.mx.Unlock()

	if !ok {
		return v, ErrorNotFound
	}
	return v, nil
}

func (s *MapStore) Add(v string) string {
	id := util.GenerateID(8)
	s.mx.Lock()
	s.values[id] = v
	s.mx.Unlock()
	return id
}
