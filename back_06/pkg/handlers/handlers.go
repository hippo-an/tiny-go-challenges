package handlers

import (
	"fmt"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/models"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/render"
	"log"
	"net/http"
)

var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handler
func NewHandlers(r *Repository) {
	Repo = r
}

func (re *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	re.App.Session.Put(r.Context(), "remote_ip", remoteIp)
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (re *Repository) About(w http.ResponseWriter, r *http.Request) {
	remoteIp := re.App.Session.GetString(r.Context(), "remote_ip")
	log.Println(remoteIp)
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (re *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}

func (re *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (re *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (re *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func (re *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	startDate := r.Form.Get("start-date")
	endDate := r.Form.Get("end-date")

	fmt.Fprintln(w, startDate, endDate)
}

func (re *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}
