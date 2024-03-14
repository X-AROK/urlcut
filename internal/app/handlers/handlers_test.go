package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/X-AROK/urlcut/internal/app/store"
	"github.com/X-AROK/urlcut/internal/app/url"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type send struct {
	method    string
	path      string
	data      string
	urlParams map[string]string
}
type want struct {
	code     int
	response string
	headers  map[string]string
}
type test struct {
	name string
	send send
	want want
}

func runTests(t *testing.T, handler http.HandlerFunc, tests []test) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer([]byte(tt.send.data))
			req := httptest.NewRequest(tt.send.method, tt.send.path, b)

			ctx := chi.NewRouteContext()
			for k, v := range tt.send.urlParams {
				ctx.URLParams.Add(k, v)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			w := httptest.NewRecorder()
			handler(w, req)
			res := w.Result()

			defer res.Body.Close()
			respBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			get := string(respBody)

			assert.Equal(t, tt.want.code, res.StatusCode)

			contentType := res.Header.Get("Content-Type")
			if strings.Contains(contentType, "application/json") {
				if strings.HasPrefix(get, "[") {
					var exp, act any
					require.NoError(t, json.Unmarshal([]byte(get), &act))
					require.NoError(t, json.Unmarshal([]byte(tt.want.response), &exp))
					assert.ElementsMatch(t, exp, act)
				} else {
					assert.JSONEq(t, tt.want.response, get)
				}
			} else {
				assert.Equal(t, tt.want.response, get)
			}

			for header, value := range tt.want.headers {
				h := res.Header.Get(header)
				assert.Equal(t, value, h)
			}
		})
	}
}

func TestCreateShortHandler(t *testing.T) {
	m := url.NewManager(store.NewMockStore())
	baseURL := "http://localhost:8000"
	handler := createShort(context.Background(), m, baseURL)

	tests := []test{
		{
			name: "create url",
			send: send{
				method:    http.MethodPost,
				path:      "/",
				data:      "https://practicum.yandex.ru",
				urlParams: map[string]string{},
			},
			want: want{
				code:     http.StatusCreated,
				response: baseURL + "/test",
				headers:  map[string]string{},
			},
		},
	}

	runTests(t, handler, tests)
}

func TestGoToIDHandler(t *testing.T) {
	m := url.NewManager(store.NewMockStore())
	handler := goToID(context.Background(), m)

	tests := []test{
		{
			name: "redirect by id",
			send: send{
				method:    http.MethodGet,
				path:      "/{id}",
				data:      "",
				urlParams: map[string]string{"id": "test"},
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
			name: "id not found",
			send: send{
				method:    http.MethodGet,
				path:      "/{id}",
				data:      "",
				urlParams: map[string]string{"id": "test2"},
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "id not found\n",
				headers:  map[string]string{},
			},
		},
	}

	runTests(t, handler, tests)
}

func TestShortenHandler(t *testing.T) {
	m := url.NewManager(store.NewMockStore())
	baseURL := "http://localhost:8000"
	handler := shorten(context.Background(), m, baseURL)

	tests := []test{
		{
			name: "create url",
			send: send{
				method:    http.MethodPost,
				path:      "/",
				data:      "{\"url\": \"https://practicum.yandex.ru\"}",
				urlParams: map[string]string{},
			},
			want: want{
				code:     http.StatusCreated,
				response: fmt.Sprintf("{\"result\":\"%s/test\"}", baseURL),
				headers:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "incorrect json #1",
			send: send{
				method:    http.MethodPost,
				path:      "/",
				data:      "{\"url\" \"https://practicum.yandex.ru\"}",
				urlParams: map[string]string{},
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "invalid character '\"' after object key\n",
				headers:  map[string]string{},
			},
		},
		{
			name: "incorrect json #2",
			send: send{
				method:    http.MethodPost,
				path:      "/",
				data:      "{",
				urlParams: map[string]string{},
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "unexpected EOF\n",
				headers:  map[string]string{},
			},
		},
		{
			name: "incorrect json #3",
			send: send{
				method:    http.MethodPost,
				path:      "/",
				data:      "{\"url\": [1, 2, 3]}",
				urlParams: map[string]string{},
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "json: cannot unmarshal array into Go struct field ShortenRequest.url of type string\n",
				headers:  map[string]string{},
			},
		},
	}

	runTests(t, handler, tests)
}

func TestBatchHandler(t *testing.T) {
	m := url.NewManager(store.NewMockStore())
	baseURL := "http://localhost:8000"
	handler := batch(context.Background(), m, baseURL)

	tests := []test{
		{
			name: "create url",
			send: send{
				method:    http.MethodPost,
				path:      "/",
				data:      "[{\"correlation_id\": \"1\", \"original_url\": \"\"}]",
				urlParams: map[string]string{},
			},
			want: want{
				code:     http.StatusCreated,
				response: fmt.Sprintf("[{\"correlation_id\": \"1\", \"short_url\": \"%s/test1\"}]", baseURL),
				headers:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "create urls",
			send: send{
				method:    http.MethodPost,
				path:      "/",
				data:      "[{\"correlation_id\": \"1\", \"original_url\": \"https://practicum.yandex.ru\"}, {\"correlation_id\": \"2\", \"original_url\": \"https://vk.com\"}]",
				urlParams: map[string]string{},
			},
			want: want{
				code:     http.StatusCreated,
				response: fmt.Sprintf("[{\"correlation_id\": \"1\", \"short_url\": \"%s/test1\"}, {\"correlation_id\": \"2\", \"short_url\": \"%[1]s/test2\"}]", baseURL),
				headers:  map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "empty",
			send: send{
				method:    http.MethodPost,
				path:      "/",
				data:      "[]",
				urlParams: map[string]string{},
			},
			want: want{
				code:     http.StatusCreated,
				response: "[]",
				headers:  map[string]string{"Content-Type": "application/json"},
			},
		},
	}

	runTests(t, handler, tests)
}
