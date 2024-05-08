package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/driver"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/handlers"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/helpers"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/render"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/repository"
)

const portNumber = ":8008"

var dbHost *string
var dbName *string
var dbUser *string
var dbPass *string
var dbPort *string
var dbSSL *string

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
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", *dbUser, *dbPass, *dbHost, *dbPort, *dbName, *dbSSL)
	db, err := driver.ConnectSQL("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)

	listenForMail()

	//from := "sehyeong@here.com"
	//auth := smtp.PlainAuth("", from, "", "localhost")
	//smtp.SendMail("localhost:1025", auth, from, []string{"hello@there.com"}, []byte("Hello world"))

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

	// read flags

	inProduction := flag.Bool("production", false, "Application is in production")
	useCache := flag.Bool("cache", false, "Use template cache")
	dbHost = flag.String("dbhost", "localhost", "Database host")
	dbName = flag.String("dbname", "booking", "Database name")
	dbUser = flag.String("dbuser", "root", "Database user")
	dbPass = flag.String("dbpass", "secret", "Database password")
	dbPort = flag.String("dbport", "15432", "Database port")
	dbSSL = flag.String("dbssl", "disable", "Database ssl settings(disable, prefer, require)")

	flag.Parse()

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// production mode settings =============================
	app.InProduction = *inProduction

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
	app.UseCache = *useCache

	// session settings =====================================
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return nil
}
