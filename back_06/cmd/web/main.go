package main

import (
	"context"
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/handlers"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/render"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const portNumber = ":8008"

var (
	app     config.AppConfig
	session *scs.SessionManager
)

// main is main application function
func main() {

	err := run()

	if err != nil {
		log.Fatal(err)
	}
	svr := http.Server{
		Handler:      route(&app),
		Addr:         portNumber,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
	}

	go func() {
		log.Println("server starting on port", portNumber)
		if err := svr.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	svr.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down....")
	os.Exit(0)
}

func run() error {
	gob.Register(models.Reservation{})

	app.InProduction = false

	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	tc, err := render.GenerateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = app.InProduction

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)

	return nil
}
