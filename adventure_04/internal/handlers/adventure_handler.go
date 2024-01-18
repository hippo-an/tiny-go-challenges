package handlers

import (
	"github.com/dev-hippo-an/tiny-go-challenges/adventure_04/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
)

var templateMap = make(map[string]string)

func init() {
	templateMap["index"] = "./adventure_04/web/index.gohtml"
}

type AdventureHandler struct {
	story *models.Story
}

func (h AdventureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	indexHtml, err := os.ReadFile(templateMap["index"])
	if err != nil {
		log.Fatal("error while read index html file; ", err)
	}
	tmpl := template.Must(template.New("index").Parse(string(indexHtml)))
	story := *h.story
	err = tmpl.Execute(w, story["intro"])
	if err != nil {
		log.Fatal("error while execute index template parsing")
	}
}

func NewHandler(s *models.Story) AdventureHandler {
	return AdventureHandler{
		story: s,
	}
}
