package handlers

import (
	"fmt"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/forms"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/helpers"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/render"
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

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (re *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (re *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (re *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["reservation"] = models.Reservation{}

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (re *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	reservation, ok := re.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		re.App.ErrorLog.Println("cannot get item from session")
		re.App.Session.Put(r.Context(), "error", "reservation is not correctly set")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	re.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (re *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	re.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
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

func (re *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {

}

func (re *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}
