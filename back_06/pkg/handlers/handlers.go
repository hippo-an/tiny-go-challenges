package handlers

import (
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/config"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/models"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/pkg/render"
	"net/http"
)

var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handler
func NewHandlers(r *Repository) {
	Repo = r
}

func (re *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "home.page.tmpl", &models.TemplateData{})
}

func (re *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, "about.page.tmpl", &models.TemplateData{})
}
