package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/hello", nil)

	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("invalid when should have been valid")
	}
}

func TestFrom_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/hello", nil)

	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("should not valid when required fields are missing")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "a")
	postData.Add("c", "a")
	r = httptest.NewRequest("POST", "/hello", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("should valid when required fields are missing")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/hello", nil)

	form := New(r.PostForm)
	form.has("a")

	if form.Valid() {
		t.Error("should not valid when required fields are missing")

	}
}
