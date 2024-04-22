package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	params             []postData
	expectedStatusCode int
}{
	{
		"home", "/", "GET", []postData{}, http.StatusOK,
	},
	{
		"about", "/about", "GET", []postData{}, http.StatusOK,
	},
	{
		"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK,
	},
	{
		"ms", "/majors-suite", "GET", []postData{}, http.StatusOK,
	},
	{
		"ms", "/search-availability", "GET", []postData{}, http.StatusOK,
	},
	{
		"ms", "/contact", "GET", []postData{}, http.StatusOK,
	},
	{
		"ms", "/make-reservation", "GET", []postData{}, http.StatusOK,
	},
	{
		"ms", "/reservation-summary", "GET", []postData{}, http.StatusOK,
	},
	{
		"post-search-avail", "/search-availability", "POST", []postData{
			{key: "start-date", value: "2024-01-01"},
			{key: "end-date", value: "2024-01-31"},
		}, http.StatusOK,
	},
	{
		"post-search-avail-json", "/search-availability-json", "POST", []postData{}, http.StatusOK,
	},
	{
		"make-reservation", "/make-reservation", "POST", []postData{
			{key: "first_name", value: "Sehyeong"},
			{key: "last_name", value: "An"},
			{key: "email", value: "sehyeong@aaaa.com"},
			{key: "phone", value: "010-xxxx-xxxx"},
		}, http.StatusOK,
	},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			res, err := ts.Client().Get(ts.URL + e.url)

			if err != nil {
				t.Error(err)
			}

			if res.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, res.StatusCode)
			}
		} else {
			values := url.Values{}

			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			res, err := ts.Client().PostForm(ts.URL+e.url, values)

			if err != nil {
				t.Error(err)
			}

			if res.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, res.StatusCode)
			}
		}
	}
}
