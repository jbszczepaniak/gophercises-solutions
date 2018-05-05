package urlshort

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var pathsToUrls = map[string]string{
	"short1": "http://short1.com",
}

type SpyHandler struct {
	calls int
}

func (spy *SpyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	spy.calls++
}

func TestMapHandler(t *testing.T) {
	cases := []struct {
		name          string
		path          string
		code          int
		location      string
		fallbackCalls int
	}{
		{
			name:          "for known path user is redirected",
			path:          "short1",
			code:          http.StatusFound,
			location:      pathsToUrls["short1"],
			fallbackCalls: 0,
		},
		{
			name:          "for not known path fallback is used",
			path:          "not known",
			code:          http.StatusOK,
			location:      "",
			fallbackCalls: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fallback := &SpyHandler{}
			handler := MapHandler(pathsToUrls, fallback)

			request, _ := http.NewRequest(http.MethodGet, c.path, nil)
			response := httptest.NewRecorder()

			handler(response, request)

			if c.code != response.Code {
				t.Errorf("want %d, got %d status code", c.code, response.Code)
			}
			if c.location != "" && c.location != response.HeaderMap["Location"][0] {
				t.Errorf("want %s, got %s location", c.location, response.HeaderMap["Location"][0])
			}
			if fallback.calls != c.fallbackCalls {
				t.Errorf("want %d, got %d fallback calls", c.fallbackCalls, fallback.calls)
			}
		})
	}
}
