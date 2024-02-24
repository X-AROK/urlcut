package url

type URL struct {
	Addr string `json:"url"`
}

type Result struct {
	URL string `json:"result"`
}

func NewURL(addr string) URL {
	return URL{Addr: addr}
}

func NewResult(url string) Result {
	return Result{URL: url}
}
