package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/driver"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/handlers"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/render"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/repository"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8008"

var (
	app     config.AppConfig
	session *scs.SessionManager
)

// main is main application function
func main() {
	err := setupAppConfig()

	if err != nil {
		log.Fatal(err)
	}

	// db settings ==========================================
	db, err := driver.ConnectSQL("postgres", "postgresql://root:secret@localhost:15432/booking?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	repo := repository.NewPostgresRepository(db)

	cfg := config.NewConfig(&app, repo)
	handlers.NewHandlers(cfg)

	svr := http.Server{
		Handler:      route(&app),
		Addr:         portNumber,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
	}

	log.Println("server starting on port", portNumber)
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func setupAppConfig() error {
	// encode GOB
	gob.Register(models.User{})
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// production mode settings =============================
	app.InProduction = false

	// logger settings ======================================
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	// go template settings =================================
	tc, err := render.GenerateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = app.InProduction

	// session settings =====================================
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	render.NewRenderer(&app)

	return nil
}
