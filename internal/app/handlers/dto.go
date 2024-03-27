package handlers

import "github.com/X-AROK/urlcut/internal/app/url"

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	URL string `json:"result"`
}

type BatchRequestElement struct {
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type BatchRequest []BatchRequestElement

type BatchResponseElement struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type BatchResponse []BatchResponseElement

func NewShortenResponse(url *url.URL, baseURL string) ShortenResponse {
	return ShortenResponse{URL: baseURL + "/" + url.ShortURL}
}

func NewBatchResponse(urls *url.URLsBatch, baseURL string) BatchResponse {
	resp := make([]BatchResponseElement, len(*urls))
	i := 0
	for k, v := range *urls {
		resp[i] = BatchResponseElement{ID: k, ShortURL: baseURL + "/" + v.ShortURL}
		i++
	}

	return resp
}

func (br BatchRequest) ToURLs() *url.URLsBatch {
	urls := &url.URLsBatch{}
	for _, v := range br {
		(*urls)[v.ID] = url.NewURL(v.OriginalURL)
	}

	return urls
}
