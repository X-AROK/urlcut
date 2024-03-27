package url

type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLsBatch map[string]*URL

func NewURL(originalURL string) *URL {
	return &URL{ShortURL: "", OriginalURL: originalURL}
}
