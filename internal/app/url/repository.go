package url

import "context"

type Repository interface {
	Get(context.Context, string) (*URL, error)
	Add(context.Context, *URL) (string, error)
}
