package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("templates/*.tmpl")),
	}
}

type Count struct {
	Count int
}

type Contact struct {
	Name  string
	Email string
}

func newContact(name, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
	}
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

func (d Data) hasEmail(email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func newData() Data {
	return Data{
		Contacts: Contacts{
			newContact("John", "john@example.com"),
			newContact("Pale", "pale@example.com"),
			newContact("Trace", "trace@example.com"),
			newContact("Jane", "jane@example.com"),
			newContact("Log", "log@example.com"),
			newContact("Jump", "jump@example.com"),
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	count := &Count{}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "count-index", count)
	})

	e.POST("/count", func(c echo.Context) error {
		count.Count++
		return c.Render(http.StatusOK, "count", count)
	})

	page := newPage()

	e.GET("/contacts", func(c echo.Context) error {
		return c.Render(http.StatusOK, "contact-index", page)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.hasEmail(email) {
			fd := newFormData()
			fd.Values["name"] = name
			fd.Values["email"] = email
			fd.Errors["email"] = "Email already exists"
			page.Form = fd
			return c.Render(http.StatusUnprocessableEntity, "form", page.Form)
		}

		contact := newContact(name, email)
		page.Data.Contacts = append(
			page.Data.Contacts,
			contact,
		)

		c.Render(http.StatusOK, "form", newFormData())
		return c.Render(http.StatusOK, "oob-contact", contact)
	})

	log.Fatal(e.Start(":8080"))
}
