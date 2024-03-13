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

func (m *Manager) GetURL(ctx context.Context, id string) (*URL, error) {
	url, err := m.s.Get(ctx, id)
	return url, err
}
