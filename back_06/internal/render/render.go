package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

var pathToTemplates = "./templates"

func NewTemplate(conf *config.AppConfig) {
	app = conf
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var templateCache map[string]*template.Template

	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = GenerateTemplateCache()
	}

	t, ok := templateCache[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func GenerateTemplateCache() (map[string]*template.Template, error) {
	fullTemplateCache := map[string]*template.Template{}
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return fullTemplateCache, err
	}

	layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
	if err != nil {
		return fullTemplateCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		createdTemplate, err := template.New(name).ParseFiles(page)
		if err != nil {
			return fullTemplateCache, err
		}

		if len(layouts) > 0 {
			associatedTemplateWithLayout, err := createdTemplate.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return fullTemplateCache, err
			}
			createdTemplate = associatedTemplateWithLayout
		}

		fullTemplateCache[name] = createdTemplate
	}

	return fullTemplateCache, nil
}

// archived code =========================================
var cachedTemplate = make(map[string]*template.Template)

func TemplateTestSimple(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	if _, ok := cachedTemplate[t]; !ok {
		log.Println("creating template and adding to cache")
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		log.Println("using cached template")
	}

	tmpl = cachedTemplate[t]
	err = tmpl.Execute(w, nil)

	if err != nil {

	}
}

func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	cachedTemplate[t] = tmpl
	return nil
}
