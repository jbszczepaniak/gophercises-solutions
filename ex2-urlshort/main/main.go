package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	urlshort "github.com/jedruniu/gophercises-solutions/ex2-urlshort"
)

var yamlFile string

func init() {
	flag.StringVar(&yamlFile, "yaml", "", "file with short path -> url mappings")
	flag.Parse()
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := getYaml()
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	err = http.ListenAndServe(":8080", yamlHandler)
	if err != nil {
		panic(err)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

var fallbackYaml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

func getYaml() []byte {
	if yamlFile == "" {
		return []byte(fallbackYaml)
	}
	file, err := os.Open(yamlFile)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return content
}
