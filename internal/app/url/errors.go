package url

import "errors"

var (
	ErrorNotFound      = errors.New("url not found")
	ErrorAlreadyExists = errors.New("original url already exists")
)
