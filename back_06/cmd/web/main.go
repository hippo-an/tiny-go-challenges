package main

import (
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/handlers"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/render"
	"log"
	"net/http"
)

const portNumber = ":8008"

func main() {
	var app config.AppConfig
	tc, err := render.GenerateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)

	http.HandleFunc("/", repo.Home)
	http.HandleFunc("/about", repo.About)
	log.Fatal(http.ListenAndServe(portNumber, nil))
}
