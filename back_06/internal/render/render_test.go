package render

import (
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()

	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "hello")
	result := AddDefaultData(&td, r)

	if result.Flash != "hello" {
		t.Error("flash value of hello not found in session")
	}
}

func TestTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := GenerateTemplateCache()

	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter
	if err := Template(&ww, r, "home.page.tmpl", &models.TemplateData{}); err != nil {
		t.Error(err)
	}
	if err := Template(&ww, r, "non-existent.page.tmpl", &models.TemplateData{}); err == nil {
		t.Error(err)
	}

}

func TestNewTemplates(t *testing.T) {
	NewTemplate(app)
}

func TestGenerateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := GenerateTemplateCache()

	if err != nil {
		t.Error(err)
	}

}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}
