package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hippo-an/tiny-go-challenges/adventure_04/internal/handlers"
	"github.com/hippo-an/tiny-go-challenges/adventure_04/internal/models"
)

func main() {
	port := flag.Int64("port", 3000, "the port to start the Adventure web application on")
	fileName := flag.String("file", "gopher.json", "the JSON file with the choose your own adventure story")

	flag.Parse()

	file, err := os.Open(fmt.Sprintf("./adventure_04/%s", *fileName))

	if err != nil {
		log.Fatalf("error while reading file %s; %s\n", *fileName, err)
	}

	story, err := models.JsonStory(file)

	if err != nil {
		log.Fatal("error while decode file to struct")
	}

	handler := handlers.NewHandler(&story)
	log.Println("server is starting on port: ", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}
