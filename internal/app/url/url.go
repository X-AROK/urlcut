package url

import "errors"

var (
	ErrorNotFound = errors.New("id not found")
)

type URL struct {
	Addr string
}

func NewURL(addr string) URL {
	return URL{Addr: addr}
}
