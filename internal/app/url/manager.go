package url

import "context"

type Manager struct {
	s Repository
}

func NewManager(s Repository) *Manager {
	return &Manager{s: s}
}

func (m *Manager) AddURL(ctx context.Context, u *URL) (string, error) {
	return m.s.Add(ctx, u)
}

func (m *Manager) AddURLsBatch(ctx context.Context, urls *URLsBatch) error {
	return m.s.AddBatch(ctx, urls)
}

func (m *Manager) GetURL(ctx context.Context, id string) (*URL, error) {
	return m.s.Get(ctx, id)
}
