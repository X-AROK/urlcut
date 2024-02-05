package store

type Repository interface {
	Get(string) (string, bool)
	Add(string) string
}
