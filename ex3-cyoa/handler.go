package handler

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type Story struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type Stories map[string]Story

func NewCyoaHandler(stories string) (*cyoaHandler, error) {
	var unmarshalled Stories
	err := json.Unmarshal([]byte(stories), &unmarshalled)
	if err != nil {
		return nil, err
	}
	return &cyoaHandler{stories: unmarshalled, curr: "intro"}, err
}

type cyoaHandler struct {
	stories Stories
	curr    string
}

func (h *cyoaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("page_template.html")
	tmplContent, _ := ioutil.ReadAll(file)
	templ, _ := template.New("page").Parse(string(tmplContent))
	templ.Execute(w, h.stories[h.curr])
}
