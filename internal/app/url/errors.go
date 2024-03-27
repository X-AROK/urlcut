package url

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("url not found")
)

type AlreadyExistsError struct {
	ID string
}

func (err *AlreadyExistsError) Error() string {
	return fmt.Sprintf("original url already exists with id '%s'", err.ID)
}

func NewAlreadyExistsError(id string) *AlreadyExistsError {
	return &AlreadyExistsError{ID: id}
}
