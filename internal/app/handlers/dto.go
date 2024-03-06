package handlers

import "github.com/X-AROK/urlcut/internal/app/url"

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	URL string `json:"result"`
}

func NewShortenResponse(url *url.URL, baseURL string) ShortenResponse {
	return ShortenResponse{URL: baseURL + "/" + url.ShortURL}
}
