package handlers

import (
	"fmt"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/forms"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/helpers"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/render"
	"net/http"
	"strconv"
	"time"
)

var conf *config.Config

// NewHandlers sets the repository for the handler
func NewHandlers(c *config.Config) {
	conf = c
}

func Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func Reservation(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["reservation"] = models.Reservation{}

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func ReservationSummary(w http.ResponseWriter, r *http.Request) {

	reservation, ok := conf.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		conf.App.ErrorLog.Println("cannot get item from session")
		conf.App.Session.Put(r.Context(), "error", "reservation is not correctly set")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	conf.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("start_date")

	dateLayout := "2006-01-02"

	startDate, err := time.Parse(dateLayout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(dateLayout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomId,
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

	err = conf.Repo.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	conf.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func PostAvailability(w http.ResponseWriter, r *http.Request) {
	startDate := r.Form.Get("start-date")
	endDate := r.Form.Get("end-date")

	fmt.Fprintln(w, startDate, endDate)
}

func AvailabilityJson(w http.ResponseWriter, r *http.Request) {

}

func Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}
