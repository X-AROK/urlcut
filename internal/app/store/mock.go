package store

type MockStore struct{}

func NewMockStore() MockStore {
	return MockStore{}
}

func (s MockStore) Get(id string) (string, bool) {
	if id == "test" {
		return "https://practicum.yandex.ru", true
	}

	return "", false
}

func (s MockStore) Add(_ string) string {
	return "test"
}
