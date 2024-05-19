package main

import (
	"github.com/dev-hippo-an/tiny-go-challenges/htmxepl/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()

	router.Get("/foo", handlers.Make(handlers.HandleFoo))

	listenAddr := os.Getenv("LISTEN_ADDR")
	slog.Info("HTTP server started", "listenAddr", listenAddr)
	err := http.ListenAndServe(listenAddr, router)
	log.Fatal(err)
}
