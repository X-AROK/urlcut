package url

type Manager struct {
	s Repository
}

func NewManager(s Repository) *Manager {
	return &Manager{s: s}
}

func (m *Manager) AddURL(u URL) (string, error) {
	return m.s.Add(u), nil
}

func (m *Manager) GetURL(id string) (URL, error) {
	url, err := m.s.Get(id)
	return url, err
}
