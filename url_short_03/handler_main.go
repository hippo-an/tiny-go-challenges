package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/hippo-an/tiny-go-challenges/url_short_03/handlers"
)

const (
	jsonString = `
[
  {
    "path": "/dev-hippo-an",
    "url": "https://github.com/dev-hippo-an"
  }
]
`
)

func main() {
	isLoadYamlFile := flag.Bool("load-yaml-file", false, "load the YAML file from a file rather than from a json format string")
	flag.Parse()
	fmt.Println("this is isLoadYamlFile: ", *isLoadYamlFile)

	defaultHandler := defaultHandler()

	var pathBytes []byte
	// Build the Handler using the default map handler as the fallback
	var handler http.Handler
	if *isLoadYamlFile {
		yamlPath, err := os.ReadFile("url_short_03/paths.yaml")
		if err != nil {
			fmt.Println("error while reading paths yaml file:", err)
			return
		}
		pathBytes = yamlPath
		yamlHandler, err := handlers.YAMLHandler(pathBytes, defaultHandler)
		if err != nil {
			return
		}
		handler = yamlHandler
	} else {
		pathBytes = []byte(jsonString)
		jsonHandler, err := handlers.JSONHandler(pathBytes, defaultHandler)
		if err != nil {
			return
		}
		handler = jsonHandler
	}
	fmt.Println("Starting the server on :8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		return
	}
}

func defaultHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "Welcome to the new world..!!")
	})

	// Build the MapHandler using the mux as the fallback
	defaultPathToUrls := map[string]string{
		"/naver": "https://naver.com",
		"/yaml":  "https://godoc.org/gopkg.in/yaml.v2",
	}

	return handlers.MapHandler(defaultPathToUrls, mux)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
