package render

import (
	"bytes"
	"fmt"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

func NewTemplate(conf *config.AppConfig) {
	app = conf
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func Template(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var templateCache map[string]*template.Template

	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = GenerateTemplateCache()
	}

	t, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("can not find t")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func GenerateTemplateCache() (map[string]*template.Template, error) {
	fullTemplateCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return fullTemplateCache, err
	}

	layouts, err := filepath.Glob("./templates/*.layout.tmpl")
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
			associatedTemplateWithLayout, err := createdTemplate.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return fullTemplateCache, err
			}
			createdTemplate = associatedTemplateWithLayout
		}

		fullTemplateCache[name] = createdTemplate
	}

	return fullTemplateCache, nil
}

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
