package main

import (
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func route(app *config.AppConfig) http.Handler {

	r := mux.NewRouter()

	r.Use(LoggingMiddleware)
	r.Use(CSRFMiddleware)
	r.Use(SessionLoad)

	r.HandleFunc("/", handlers.Repo.Home).Methods("GET")
	r.HandleFunc("/about", handlers.Repo.About).Methods("GET")

	return r
}
