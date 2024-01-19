package handlers

import (
	"github.com/dev-hippo-an/tiny-go-challenges/adventure_04/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type AdventureHandler struct {
	story *models.Story
}

var pathHtml string

func init() {
	templateBytes, err := os.ReadFile("./adventure_04/web/template.gohtml")
	if err != nil {
		log.Fatal("error while read index html file; ", err)
	}

	pathHtml = string(templateBytes)
}

func (h AdventureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	tmpl := template.Must(template.New(path).Parse(pathHtml))
	story := *h.story

	if chapter, ok := story[path[1:]]; ok {
		err := tmpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	} else {
		http.Error(w, "Chapter not found.", http.StatusNotFound)
	}

}

func NewHandler(s *models.Story) AdventureHandler {
	return AdventureHandler{
		story: s,
	}
}
