package store

type MockStore struct{}

func NewMockStore() MockStore {
	return MockStore{}
}

func (s MockStore) Get(id string) (string, error) {
	if id == "test" {
		return "https://practicum.yandex.ru", nil
	}

	return "", ErrorNotFound
}

func (s MockStore) Add(_ string) string {
	return "test"
}
