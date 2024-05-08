package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/forms"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/helpers"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/render"
	"github.com/gorilla/mux"
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

	res, ok := conf.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		conf.App.Session.Put(r.Context(), "error", "reservation is not correctly set")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	room, err := conf.Repo.GetRoomById(res.RoomID)

	if err != nil {
		conf.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	conf.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func ReservationSummary(w http.ResponseWriter, r *http.Request) {

	reservation, ok := conf.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		conf.App.ErrorLog.Println("cannot get item from session")
		conf.App.Session.Put(r.Context(), "error", "reservation is not correctly set")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	conf.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		conf.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	dateLayout := "2006-01-02"

	startDate, err := time.Parse(dateLayout, sd)
	if err != nil {
		conf.App.Session.Put(r.Context(), "error", "start date format is incorrect")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(dateLayout, ed)
	if err != nil {
		conf.App.Session.Put(r.Context(), "error", "end date format is incorrect")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		conf.App.Session.Put(r.Context(), "error", "room id is incorrect")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
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

	reservationId, err := conf.Repo.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		conf.App.Session.Put(r.Context(), "error", "can't insert reservation into database")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        roomId,
		ReservationID: reservationId,
		RestrictionID: 1,
	}

	err = conf.Repo.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		conf.App.Session.Put(r.Context(), "error", "can't insert room restriction")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	htmlMsg := fmt.Sprintf(`
		<h1>Reservation Confirmation</h1>
		Dear %s,
		The reservation confirmed from %s to %s
	`, reservation.FirstName, sd, ed)

	msg := models.MailData{
		To:      reservation.Email,
		From:    "me@here.com",
		Subject: "Reservation Confirmation.",
		Content: htmlMsg,
	}

	conf.App.MailChan <- msg

	conf.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start-date")
	ed := r.Form.Get("end-date")
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

	rooms, err := conf.Repo.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
	}

	if len(rooms) == 0 {
		conf.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	conf.App.Session.Put(r.Context(), "reservation", reservation)
	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start-modal")
	ed := r.Form.Get("end-modal")

	format := "2006-01-02"

	startDate, err := time.Parse(format, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(format, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	availability, err := conf.Repo.SearchAvailabilityByDateByRoomId(startDate, endDate, roomId)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	response := struct {
		Ok        bool   `json:"ok"`
		Msg       string `json:"msg"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		RoomId    int    `json:"room_id"`
	}{
		Ok:        availability,
		Msg:       "",
		StartDate: sd,
		EndDate:   ed,
		RoomId:    roomId,
	}

	prettyJson, err := json.MarshalIndent(response, "", "  ")

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(prettyJson)
	return
}

func ChooseRoom(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	roomId, ok := params["roomId"]
	if !ok {
		helpers.ServerError(w, errors.New("id is not correct"))
		return
	}

	res, ok := conf.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		conf.App.ErrorLog.Println("cannot get reservation from session")
		conf.App.Session.Put(r.Context(), "error", "reservation is not correctly set")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	intRoomId, err := strconv.Atoi(roomId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.RoomID = intRoomId

	conf.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
func BookRoom(w http.ResponseWriter, r *http.Request) {
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	format := "2006-01-02"

	startDate, err := time.Parse(format, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(format, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var res models.Reservation
	res.RoomID = roomId
	res.StartDate = startDate
	res.EndDate = endDate

	conf.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		conf.App.Session.Put(r.Context(), "error", "can't parse the form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required(email, password)

	if !form.Valid() {
		return
	}

	id, err := conf.Repo.Authenticate(email, password)

	if err != nil {
		helpers.ServerError(w, err)
		conf.App.Session.Put(r.Context(), "error", "can't parse the form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	conf.App.Session.Put(r.Context(), "user_id", id)
	conf.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func PostLogout(w http.ResponseWriter, r *http.Request) {
	_ = conf.App.Session.Destroy(r.Context())
	_ = conf.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
