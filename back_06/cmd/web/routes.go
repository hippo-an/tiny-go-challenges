package main

import (
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func route(app *config.AppConfig) http.Handler {

	r := mux.NewRouter()

	r.Use(LoggingMiddleware)
	r.Use(CSRFMiddleware)
	r.Use(SessionLoad)

	_ = r.HandleFunc("/", handlers.Repo.Home).Methods("GET")
	_ = r.HandleFunc("/about", handlers.Repo.About).Methods("GET")
	_ = r.HandleFunc("/generals-quarters", handlers.Repo.Generals).Methods("GET")
	_ = r.HandleFunc("/majors-suite", handlers.Repo.Majors).Methods("GET")

	_ = r.HandleFunc("/search-availability", handlers.Repo.Availability).Methods("GET")
	_ = r.HandleFunc("/search-availability", handlers.Repo.PostAvailability).Methods("POST")
	_ = r.HandleFunc("/search-availability-json", handlers.Repo.AvailabilityJson).Methods("POST")

	_ = r.HandleFunc("/contact", handlers.Repo.Contact).Methods("GET")
	_ = r.HandleFunc("/make-reservation", handlers.Repo.Reservation).Methods("GET")
	_ = r.HandleFunc("/make-reservation", handlers.Repo.PostReservation).
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded")

	fileServer := http.FileServer(http.Dir("./static/"))
	_ = r.PathPrefix("/").Handler(http.StripPrefix("/static", fileServer))
	return r
}
