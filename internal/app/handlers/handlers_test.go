package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/X-AROK/urlcut/internal/app/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer([]byte(body)))
	require.NoError(t, err)
	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestMainHandler(t *testing.T) {
	s := store.NewMockStore()
	baseURL := "http://localhost:8000"
	ts := httptest.NewServer(MainRouter(s, baseURL))

	type send struct {
		method string
		path   string
		data   string
	}
	type want struct {
		code     int
		response string
		headers  map[string]string
	}
	tests := []struct {
		name string
		send send
		want want
	}{
		{
			name: "create url",
			send: send{
				method: http.MethodPost,
				path:   "/",
				data:   "https://practicum.yandex.ru",
			},
			want: want{
				code:     http.StatusCreated,
				response: baseURL + "/test",
				headers:  map[string]string{},
			},
		},
		{
			name: "redirect by id",
			send: send{
				method: http.MethodGet,
				path:   "/test",
				data:   "",
			},
			want: want{
				code:     http.StatusTemporaryRedirect,
				response: "<a href=\"https://practicum.yandex.ru\">Temporary Redirect</a>.\n\n",
				headers: map[string]string{
					"Location": "https://practicum.yandex.ru",
				},
			},
		},
		{
			name: "not allowed method",
			send: send{
				method: http.MethodPut,
				path:   "/",
				data:   "",
			},
			want: want{
				code:     http.StatusMethodNotAllowed,
				response: "",
				headers:  map[string]string{},
			},
		},
		{
			name: "id not found",
			send: send{
				method: http.MethodGet,
				path:   "/test2",
				data:   "",
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "Id not found\n",
				headers:  map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, get := testRequest(t, ts, tt.send.method, tt.send.path, tt.send.data)

			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, get)

			for header, value := range tt.want.headers {
				h := res.Header.Get(header)
				assert.Equal(t, value, h)
			}
		})
	}
}
