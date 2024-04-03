package forms

import (
	"net/http"
	"net/url"
)

type Form struct {
	url.Values
	Errors customError
}

func New(data url.Values) *Form {
	return &Form{
		data,
		map[string][]string{},
	}
}

func (f *Form) Has(field string, r *http.Request) bool {
	x := f.Values.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}

	return true
}
