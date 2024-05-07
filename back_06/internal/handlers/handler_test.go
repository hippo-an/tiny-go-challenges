package handlers

import (
	"context"
	"fmt"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{
		"home", "/", "GET", http.StatusOK,
	},
	{
		"about", "/about", "GET", http.StatusOK,
	},
	{
		"gq", "/generals-quarters", "GET", http.StatusOK,
	},
	{
		"ms", "/majors-suite", "GET", http.StatusOK,
	},
	{
		"ms", "/search-availability", "GET", http.StatusOK,
	},
	{
		"ms", "/contact", "GET", http.StatusOK,
	},
	{
		"ms", "/make-reservation", "GET", http.StatusOK,
	},
	{
		"ms", "/reservation-summary", "GET", http.StatusOK,
	},
	//{
	//	"post-search-avail", "/search-availability", "POST", []postData{
	//		{key: "start-date", value: "2024-01-01"},
	//		{key: "end-date", value: "2024-01-31"},
	//	}, http.StatusOK,
	//},
	//{
	//	"post-search-avail-json", "/search-availability-json", "POST", []postData{}, http.StatusOK,
	//},
	//{
	//	"make-reservation", "/make-reservation", "POST", []postData{
	//		{key: "first_name", value: "Sehyeong"},
	//		{key: "last_name", value: "An"},
	//		{key: "email", value: "sehyeong@aaaa.com"},
	//		{key: "phone", value: "010-xxxx-xxxx"},
	//	}, http.StatusOK,
	//},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {

		res, err := ts.Client().Get(ts.URL + e.url)

		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, res.StatusCode)
		}

	}
}

func TestPostReservation(t *testing.T) {

	reqBody := "start_date=2025-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2025-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=sehyeong")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=an")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sehyeong@aaa.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=12345678")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := loadSessionContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted: %d", rr.Code, http.StatusSeeOther)
	}

	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = loadSessionContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = PostReservation

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2025-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=sehyeong")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=an")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sehyeong@aaa.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=12345678")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = loadSessionContext(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = PostReservation

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := loadSessionContext(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted: %d", rr.Code, http.StatusOK)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = loadSessionContext(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler = Reservation

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted: %d", rr.Code, http.StatusOK)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = loadSessionContext(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 1000
	session.Put(ctx, "reservation", reservation)
	handler = Reservation

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted: %d", rr.Code, http.StatusOK)
	}

}

func loadSessionContext(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
