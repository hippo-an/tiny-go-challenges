package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hippo-an/tiny-go-challenges/adventure_04/internal/models"
)

func init() {
	templateBytes, err := os.ReadFile("./adventure_04/web/template.gohtml")
	if err != nil {
		log.Fatal("error while read index html file; ", err)
	}
	tmpl = template.Must(template.New("").Parse(string(templateBytes)))
}

type AdventureHandler struct {
	story    *models.Story
	tmpl     *template.Template
	pathFunc func(r *http.Request) string
}

func (h AdventureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFunc(r)
	story := *h.story
	if chapter, ok := story[path]; ok {
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

type HandlerOption func(h *AdventureHandler)

var tmpl *template.Template

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:]
}

func NewHandler(s *models.Story, opts ...HandlerOption) AdventureHandler {
	handler := AdventureHandler{
		story:    s,
		tmpl:     tmpl,
		pathFunc: defaultPathFn,
	}

	for _, opt := range opts {
		opt(&handler)
	}

	return handler
}

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *AdventureHandler) {
		h.tmpl = t
	}
}

func WithPathFunc(fn func(r *http.Request) string) HandlerOption {
	return func(h *AdventureHandler) {
		h.pathFunc = fn
	}
}
