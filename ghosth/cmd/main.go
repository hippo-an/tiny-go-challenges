package main

import (
	"github.com/hippo-an/tiny-go-challenges/ghosth/internal/server"
	"github.com/hippo-an/tiny-go-challenges/ghosth/internal/store"
	"log"
	"os"
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
	if err := srv.Start(); err != nil {
		logger.Fatalf("Error when starting server: %s", err)
		os.Exit(1)
	}
}
