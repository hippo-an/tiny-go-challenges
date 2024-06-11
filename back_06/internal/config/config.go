package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/hippo-an/tiny-go-challenges/back_06/internal/models"
	"github.com/hippo-an/tiny-go-challenges/back_06/internal/repository"
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
	MailChan      chan models.MailData
}
