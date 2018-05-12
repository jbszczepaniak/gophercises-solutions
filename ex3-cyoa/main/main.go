package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jedruniu/gophercises-solutions/ex3-cyoa"
)

func main() {
	file, err := os.Open("gopher.json")
	if err != nil {
		panic(err)
	}
	stories, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	cyoa, err := handler.NewCyoaHandler(string(stories), "page_template.html")
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", cyoa)
}
