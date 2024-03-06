package url

type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewURL(originalURL string) *URL {
	return &URL{ShortURL: "", OriginalURL: originalURL}
}
