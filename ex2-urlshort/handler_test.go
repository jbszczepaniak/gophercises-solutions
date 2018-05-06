package urlshort

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCase struct {
	name          string
	path          string
	code          int
	location      string
	fallbackCalls int
}

var cases = []testCase{
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

func TestMapHandler(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fallback := &SpyHandler{}
			handler, _ := MapHandler(pathsToUrls, fallback)
			testHandler(t, handler, c, fallback)
		})
	}
}

func TestYAMLHandler(t *testing.T) {
	t.Run("returns error if yaml unmarshalling did not succed", func(t *testing.T) {
		_, err := YAMLHandler([]byte("garbage yaml"), nil)
		if err == nil {
			t.Errorf("expected err to be returned but it didn't")
		}
	})

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fallback := &SpyHandler{}
			handler, _ := YAMLHandler([]byte(yamlMapping), fallback)
			testHandler(t, handler, c, fallback)
		})
	}
}

func TestJSONHandler(t *testing.T) {
	t.Run("returns error if yaml unmarshalling did not succed", func(t *testing.T) {
		_, err := JSONHandler([]byte("garbage json"), nil)
		if err == nil {
			t.Errorf("expected err to be returned but it didn't")
		}
	})

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fallback := &SpyHandler{}
			handler, _ := JSONHandler([]byte(jsonMapping), fallback)
			testHandler(t, handler, c, fallback)
		})
	}
}

var pathsToUrls = map[string]string{
	"short1": "http://short1.com",
}

var yamlMapping = `
- path: short1
  url: http://short1.com
`
var jsonMapping = `
[
	{
		"path": "short1",
		"url": "http://short1.com"
	}
]
`

type SpyHandler struct {
	calls int
}

func (spy *SpyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	spy.calls++
}

func assertCode(t *testing.T, want int, response *httptest.ResponseRecorder) {
	if want != response.Code {
		t.Errorf("want %d, got %d status code", want, response.Code)
	}
}

func assertLocation(t *testing.T, want string, response *httptest.ResponseRecorder) {
	if want != "" && want != response.HeaderMap["Location"][0] {
		t.Errorf("want %s, got %s location", want, response.HeaderMap["Location"][0])
	}
}

func assertFallbackCount(t *testing.T, want int, got int) {
	if want != got {
		t.Errorf("want %d, got %d fallback calls", want, got)
	}
}

func testHandler(t *testing.T, h http.HandlerFunc, c testCase, fallback *SpyHandler) {
	t.Helper()

	request, _ := http.NewRequest(http.MethodGet, c.path, nil)
	response := httptest.NewRecorder()

	h(response, request)

	assertCode(t, c.code, response)
	assertLocation(t, c.location, response)
	assertFallbackCount(t, c.fallbackCalls, fallback.calls)
}
