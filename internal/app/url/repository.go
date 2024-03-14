package url

import "context"

type Repository interface {
	Get(context.Context, string) (*URL, error)
	Add(context.Context, *URL) (string, error)
	AddBatch(ctx context.Context, batch *URLsBatch) error
}
