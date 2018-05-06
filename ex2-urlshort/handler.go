package urlshort

import (
	"encoding/json"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) (http.HandlerFunc, error) {
	return Handler(pathsToUrls, fallback)
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	mappings := []urlPath{}
	err := yaml.Unmarshal(yml, &mappings)
	if err != nil {
		return nil, err
	}
	return Handler(mappings, fallback)
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	mappings := []urlPath{}
	err := json.Unmarshal(jsn, &mappings)
	if err != nil {
		return nil, err
	}
	return Handler(mappings, fallback)
}

func Handler(mappings interface{}, fallback http.Handler) (http.HandlerFunc, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, ok := findURL(r.URL.String(), mappings)
		if ok {
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}), nil
}

func findURL(path string, pathMappings interface{}) (url string, found bool) {
	switch pathMappings := pathMappings.(type) {
	case []urlPath:
		for _, m := range pathMappings {
			if m.Path == path {
				url = m.URL
				found = true
				return
			}
		}
	case map[string]string:
		url, found = pathMappings[path]
	}
	return
}

type urlPath struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}
