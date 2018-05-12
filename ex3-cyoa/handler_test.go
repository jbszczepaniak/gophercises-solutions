package handler

import (
	"net/http/httptest"
	"strings"
	"testing"
)

var introStory = `
{
  "intro": {
    "title": "The Little Blue Gopher",
    "story": [
      "Once upon a time, long...",
      "One of his friends once ...",
      "On the other hand, he..."
    ],
    "options": [
      {
        "text": "That story about.",
        "arc": "new-york"
      },
      {
        "text": "Gee, those bandits.",
        "arc": "denver"
      }
    ]
	}
}
`
var introStoryAsHTML = `
<h1>The Little Blue Gopher</h1>

<p>Once upon a time, long...</p>
<p>One of his friends once ...</p>
<p>On the other hand, he...</p>

</br>

<ul>
<li><a href='/new-york'>That story about.</li>
<li><a href='/denver'>Gee, those bandits.</li>
</ul>
`
var templateURL = "main/page_template.html"

func TestNewHandlerUnmarshallStories(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		_, err := NewCyoaHandler(introStory, templateURL)
		if err != nil {
			t.Errorf("error not expected")
		}
	})
	t.Run("invalid JSON", func(t *testing.T) {
		_, err := NewCyoaHandler("invalid JSON", templateURL)
		if err == nil {
			t.Errorf("expected error did not occur")
		}
	})
}

func TestBaseUrlReturnsIntroStory(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	handler, _ := NewCyoaHandler(introStory, templateURL)

	handler.ServeHTTP(response, request)

	if stripSpaces(response.Body.String()) != stripSpaces(introStoryAsHTML) {
		t.Errorf("expected HTML differs from expected")
	}
}

func TestKnownStoryIsReturned(t *testing.T) {
	request := httptest.NewRequest("GET", "/intro", nil)
	response := httptest.NewRecorder()
	handler, _ := NewCyoaHandler(introStory, templateURL)

	handler.ServeHTTP(response, request)

	if stripSpaces(response.Body.String()) != stripSpaces(introStoryAsHTML) {
		t.Errorf("expected HTML differs from expected")
	}
}

func TestUnknownStoryReturns404(t *testing.T) {
	request := httptest.NewRequest("GET", "/unknown", nil)
	response := httptest.NewRecorder()
	handler, _ := NewCyoaHandler(introStory, templateURL)

	handler.ServeHTTP(response, request)

	got := response.Code
	want := 404

	if got != want {
		t.Errorf("want status code %d, got %d", want, got)
	}
}

func stripSpaces(text string) string {
	return strings.Replace(strings.Replace(text, " ", "", -1), "\n", "", -1)
}
