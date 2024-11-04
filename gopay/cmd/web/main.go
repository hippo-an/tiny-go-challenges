package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "0.0.1"
const cssVersion = "1"

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	toss struct {
		clientKey   string
		secretKey   string
		securityKey string
	}
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
}

func (a *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", a.config.port),
		Handler:           a.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	a.infoLog.Printf("Starting HTTP server in %s mode on port %d\n", a.config.env, a.config.port)

	return srv.ListenAndServe()
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listne on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|porduction}")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")

	flag.Parse()

	cfg.toss.clientKey = os.Getenv("CLIENT_KEY")
	cfg.toss.secretKey = os.Getenv("SECRET_KEY")
	cfg.toss.securityKey = os.Getenv("SECURITY_KEY")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
	}

	err := app.serve()
	if err != nil {
		app.errorLog.Println("err")
		log.Fatal(err)
	}

}
