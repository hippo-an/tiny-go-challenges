package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/render"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/repository"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplate = "./../../templates"

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})

	app.InProduction = false

	tc, err := GenerateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = !app.InProduction

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.InfoLog = infoLog
	app.ErrorLog = errorLog
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// db settings ==========================================
	repo := repository.NewTestRepository()

	cfg := config.NewConfig(&app, repo)
	NewHandlers(cfg)
	render.NewRenderer(&app)

	r := mux.NewRouter()

	r.Use(LoggingMiddleware)
	r.Use(SessionLoad)

	_ = r.HandleFunc("/", Home).Methods("GET")
	_ = r.HandleFunc("/about", About).Methods("GET")
	_ = r.HandleFunc("/generals-quarters", Generals).Methods("GET")
	_ = r.HandleFunc("/majors-suite", Majors).Methods("GET")

	_ = r.HandleFunc("/search-availability", SearchAvailability).Methods("GET")
	_ = r.HandleFunc("/search-availability", PostSearchAvailability).Methods("POST")
	_ = r.HandleFunc("/search-availability-json", AvailabilityJson).Methods("POST")

	_ = r.HandleFunc("/contact", Contact).Methods("GET")
	_ = r.HandleFunc("/make-reservation", Reservation).Methods("GET")
	_ = r.HandleFunc("/reservation-summary", ReservationSummary).Methods("GET")
	_ = r.HandleFunc("/make-reservation", PostReservation).
		Methods("POST").
		Headers("Content-Type", "application/x-www-form-urlencoded")

	fileServer := http.FileServer(http.Dir("./static/"))
	_ = r.PathPrefix("/").Handler(http.StripPrefix("/static", fileServer))

	return r
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" || strings.HasPrefix(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		log.Printf("[%s] %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func CSRFMiddleware(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)

}

func GenerateTestTemplateCache() (map[string]*template.Template, error) {
	fullTemplateCache := map[string]*template.Template{}
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplate))
	if err != nil {
		return fullTemplateCache, err
	}

	layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
	if err != nil {
		return fullTemplateCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		createdTemplate, err := template.New(name).ParseFiles(page)
		if err != nil {
			return fullTemplateCache, err
		}

		if len(layouts) > 0 {
			associatedTemplateWithLayout, err := createdTemplate.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
			if err != nil {
				return fullTemplateCache, err
			}
			createdTemplate = associatedTemplateWithLayout
		}

		fullTemplateCache[name] = createdTemplate
	}

	return fullTemplateCache, nil
}
