package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/repository"
	"html/template"
	"log"
)

// Config is the repository type
type Config struct {
	App  *AppConfig
	Repo repository.Repository
}

// NewConfig creates a new repository
func NewConfig(a *AppConfig, r repository.Repository) *Config {
	return &Config{
		App:  a,
		Repo: r,
	}
}

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}
