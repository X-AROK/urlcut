package url

type Repository interface {
	Get(string) (URL, error)
	Add(URL) string
}
