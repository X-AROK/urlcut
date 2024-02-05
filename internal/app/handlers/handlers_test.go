package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/X-AROK/urlcut/internal/app/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandler(t *testing.T) {
	handler := MainHandler(store.NewMockStore())

	type send struct {
		method string
		path   string
		body   string
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
				body:   "https://practicum.yandex.ru",
			},
			want: want{
				code:     http.StatusCreated,
				response: "http://example.com/test",
				headers:  map[string]string{},
			},
		},
		{
			name: "redirect by id",
			send: send{
				method: http.MethodGet,
				path:   "/test",
				body:   "",
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
				body:   "",
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "Method not allowed\n",
				headers:  map[string]string{},
			},
		},
		{
			name: "id not found",
			send: send{
				method: http.MethodGet,
				path:   "/test2",
				body:   "",
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
			b := strings.NewReader(tt.send.body)
			req := httptest.NewRequest(tt.send.method, tt.send.path, b)
			w := httptest.NewRecorder()

			handler(w, req)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(body))

			for header, value := range tt.want.headers {
				h := w.Header().Get(header)
				assert.Equal(t, value, h)
			}
		})
	}
}
