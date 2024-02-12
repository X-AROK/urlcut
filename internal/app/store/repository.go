package store

import "errors"

var (
	ErrorNotFound = errors.New("id not found")
)

type Repository interface {
	Get(string) (string, error)
	Add(string) string
}
