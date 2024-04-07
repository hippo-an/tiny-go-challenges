package forms

import (
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

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		f.has(field)
	}
}

func (f *Form) has(field string) {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
	}
}
