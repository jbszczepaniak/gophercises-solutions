package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
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

func NewCyoaHandler(stories, tmplPath string) (*cyoaHandler, error) {
	var unmarshalled Stories
	err := json.Unmarshal([]byte(stories), &unmarshalled)
	if err != nil {
		return nil, err
	}
	templ, err := template.ParseFiles(tmplPath)
	if err != nil {
		return nil, err
	}
	return &cyoaHandler{stories: unmarshalled, templ: templ}, err
}

type cyoaHandler struct {
	stories Stories
	templ   *template.Template
}

func (h *cyoaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	if name == "" {
		name = "intro"
	}

	story, ok := h.stories[name]
	if ok {
		h.templ.Execute(w, story)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Chapter not found"))
}
