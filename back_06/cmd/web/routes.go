package main

import (
	"net/http"

	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/handlers"
	"github.com/gorilla/mux"
)

func route(app *config.AppConfig) http.Handler {

	r := mux.NewRouter()

	r.Use(LoggingMiddleware)
	r.Use(CSRFMiddleware)
	r.Use(SessionLoad)

	_ = r.HandleFunc("/", handlers.Home).Methods("GET")
	_ = r.HandleFunc("/about", handlers.About).Methods("GET")
	_ = r.HandleFunc("/generals-quarters", handlers.Generals).Methods("GET")
	_ = r.HandleFunc("/majors-suite", handlers.Majors).Methods("GET")
	_ = r.HandleFunc("/contact", handlers.Contact).Methods("GET")

	_ = r.HandleFunc("/search-availability", handlers.SearchAvailability).Methods("GET")
	_ = r.HandleFunc("/search-availability", handlers.PostSearchAvailability).Methods("POST")
	_ = r.HandleFunc("/search-availability-json", handlers.AvailabilityJson).Methods("POST")

	_ = r.HandleFunc("/choose-room/{roomId}", handlers.ChooseRoom).Methods("GET")
	_ = r.HandleFunc("/book-room", handlers.BookRoom).Methods("GET")

	_ = r.HandleFunc("/make-reservation", handlers.Reservation).Methods("GET")
	_ = r.HandleFunc("/reservation-summary", handlers.ReservationSummary).Methods("GET")

	_ = r.HandleFunc("/make-reservation", handlers.PostReservation).
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded")

	fileServer := http.FileServer(http.Dir("./static/"))
	_ = r.PathPrefix("/").Handler(http.StripPrefix("/static", fileServer))

	admRouter := r.PathPrefix("/admin").Subrouter()

	admRouter.Use(Auth)

	_ = admRouter.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}).Methods("GET")
	return r
}
