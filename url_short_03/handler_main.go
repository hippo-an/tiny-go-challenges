package main

import (
	"fmt"
	"github.com/dev-hippo-an/tiny-go-challenges/url_short_03/handlers"
	"net/http"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/naver": "https://naver.com",
		"/yaml":  "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := handlers.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /dev-hippo-an
  url: https://github.com/dev-hippo-an
- path: /ejejc
  url: https://github.com/ejejc
`
	yamlHandler, err := handlers.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "Welcome to the new world..!!")
	})
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
