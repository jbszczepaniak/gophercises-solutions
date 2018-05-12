package handler

import (
	"net/http/httptest"
	"strings"
	"testing"
)

var singleStory = `
{
  "intro": {
    "title": "The Little Blue Gopher",
    "story": [
      "Once upon a time, long long ago, there was a little blue gopher. Our little blue friend wanted to go on an adventure, but he wasn't sure where to go. Will you go on an adventure with him?",
      "One of his friends once recommended going to New York to make friends at this mysterious thing called \"GothamGo\". It is supposed to be a big event with free swag and if there is one thing gophers love it is free trinkets. Unfortunately, the gopher once heard a campfire story about some bad fellas named the Sticky Bandits who also live in New York. In the stories these guys would rob toy stores and terrorize young boys, and it sounded pretty scary.",
      "On the other hand, he has always heard great things about Denver. Great ski slopes, a bad hockey team with cheap tickets, and he even heard they have a conference exclusively for gophers like himself. Maybe Denver would be a safer place to visit."
    ],
    "options": [
      {
        "text": "That story about the Sticky Bandits isn't real, it is from Home Alone 2! Let's head to New York.",
        "arc": "new-york"
      },
      {
        "text": "Gee, those bandits sound pretty real to me. Let's play it safe and try our luck in Denver.",
        "arc": "denver"
      }
    ]
	}
}
`

func TestHandlerReturns200(t *testing.T) {
	// request := httptest.NewRequest("GET", "/", nil)
	// response := httptest.NewRecorder()
	// handler := &cyoaHandler{}
	// handler.ServeHTTP(response, request)
	// want := 200
	// got := response.Code
	// if want != got {
	// 	t.Errorf("want %d, got %d status code", got, want)
	// }
}

func TestNewHandlerUnmarshallStories(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		_, err := NewCyoaHandler(singleStory)
		if err != nil {
			t.Errorf("error not expected")
		}
	})
	t.Run("invalid JSON", func(t *testing.T) {
		_, err := NewCyoaHandler("invalid JSON")
		if err == nil {
			t.Errorf("expected error did not occur")
		}
	})
}

func TestHandlerReturnsLinksFromJSON(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	handler, _ := NewCyoaHandler(singleStory)

	handler.ServeHTTP(response, request)

	if !strings.Contains(response.Body.String(), "<a href='/new-york'>") {
		t.Errorf("body does not contain link to story")
	}
	if !strings.Contains(response.Body.String(), "<a href='/new-york'>") {
		t.Errorf("body does not contain link to story")
	}
}
