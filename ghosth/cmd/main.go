package main

import (
	"log"
	"os"

	"github.com/hippo-an/tiny-go-challenges/ghosth/internal/server"
	"github.com/hippo-an/tiny-go-challenges/ghosth/internal/store"
)

func main() {
	logger := log.New(os.Stdout, "[Ghosth] ", log.LstdFlags)

	port := 9000

	logger.Print("Creating guests store..")
	guestDb := store.NewGuestStore(logger)
	guestDb.AddGuest(store.Guest{Name: "Sehyeong", Email: "sehyeong@zzzzzz.ii"})

	srv, err := server.NewServer(logger, port, guestDb)
	if err != nil {
		logger.Fatalf("Error when creating server: %s", err)
		os.Exit(1)
	}
	if err := srv.StreamStart(); err != nil {
		logger.Fatalf("Error when starting server: %s", err)
		os.Exit(1)
	}
}
