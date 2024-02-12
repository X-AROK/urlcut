package store

import "github.com/X-AROK/urlcut/internal/util"

type MapStore struct {
	values map[string]string
}

func NewMapStore() MapStore {
	return MapStore{values: make(map[string]string)}
}

func (s MapStore) Get(id string) (string, error) {
	v, ok := s.values[id]
	if !ok {
		return v, ErrorNotFound
	}
	return v, nil
}

func (s MapStore) Add(v string) string {
	id := util.GenerateID(8)
	s.values[id] = v
	return id
}
