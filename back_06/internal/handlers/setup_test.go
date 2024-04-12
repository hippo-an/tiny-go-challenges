package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/render"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
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

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplate(&app)

	r := mux.NewRouter()

	r.Use(LoggingMiddleware)
	r.Use(SessionLoad)

	_ = r.HandleFunc("/", Repo.Home).Methods("GET")
	_ = r.HandleFunc("/about", Repo.About).Methods("GET")
	_ = r.HandleFunc("/generals-quarters", Repo.Generals).Methods("GET")
	_ = r.HandleFunc("/majors-suite", Repo.Majors).Methods("GET")

	_ = r.HandleFunc("/search-availability", Repo.Availability).Methods("GET")
	_ = r.HandleFunc("/search-availability", Repo.PostAvailability).Methods("POST")
	_ = r.HandleFunc("/search-availability-json", Repo.AvailabilityJson).Methods("POST")

	_ = r.HandleFunc("/contact", Repo.Contact).Methods("GET")
	_ = r.HandleFunc("/make-reservation", Repo.Reservation).Methods("GET")
	_ = r.HandleFunc("/reservation-summary", Repo.ReservationSummary).Methods("GET")
	_ = r.HandleFunc("/make-reservation", Repo.PostReservation).
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
