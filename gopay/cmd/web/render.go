package main

import (
	"embed"
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"strings"
)

type templateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float64
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	IsAuthenticated int
	API             string
	CSSVersion      string
}

var functions = template.FuncMap{}

//go:embed templates
var templateFs embed.FS

func (a *application) addDefaultData(c echo.Context, td *templateData) *templateData {
	return td
}

func (a *application) renderTemplate(c echo.Context, page string, td *templateData, partials ...string) error {
	var t *template.Template
	var err error
	templateToRender := fmt.Sprintf("templates/%s.page.tmpl", page)

	_, templateInMap := a.templateCache[templateToRender]

	if a.config.env == "production" && templateInMap {
		t = a.templateCache[templateToRender]
	} else {
		t, err = a.parseTemplate(partials, page, templateToRender)

		if err != nil {
			a.errorLog.Println(err)
			return err
		}
	}

	if td == nil {
		td = &templateData{}
	}

	td = a.addDefaultData(c, td)

	err = t.Execute(c.Response().Writer, td)

	if err != nil {
		a.errorLog.Println(err)
		return err
	}

	return nil
}

func (a *application) parseTemplate(partials []string, page string, fullQualifiedName string) (*template.Template, error) {
	var t *template.Template
	var err error

	if len(partials) > 0 {
		for i, v := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partial.tmpl", v)
		}
	}

	fullName := fmt.Sprintf("%s.page.tmpl", page)
	if len(partials) > 0 {
		t, err = template.New(fullName).Funcs(functions).ParseFS(
			templateFs,
			"templates/base.layout.tmpl", strings.Join(partials, ","), fullQualifiedName)
	} else {
		t, err = template.New(fullName).Funcs(functions).ParseFS(
			templateFs,
			"templates/base.layout.tmpl", fullQualifiedName)
	}

	if err != nil {
		a.errorLog.Println(err)
		return nil, err
	}

	a.templateCache[fullQualifiedName] = t

	return t, nil
}
