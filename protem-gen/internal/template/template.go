package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	protemgen "github.com/hippo-an/tiny-go-challenges/protem-gen"
	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/config"
)

// Engine handles template rendering
type Engine struct {
	funcMap template.FuncMap
}

// NewEngine creates a new template engine
func NewEngine() *Engine {
	return &Engine{
		funcMap: template.FuncMap{
			"title":    strings.Title,
			"lower":    strings.ToLower,
			"upper":    strings.ToUpper,
			"replace":  strings.ReplaceAll,
			"contains": strings.Contains,
			"hasGRPC": func(cfg *config.ProjectConfig) bool {
				return cfg.IncludeGRPC
			},
			"hasAuth": func(cfg *config.ProjectConfig) bool {
				return cfg.IncludeAuth
			},
			"hasAI": func(cfg *config.ProjectConfig) bool {
				return cfg.IncludeAI
			},
			"hasDB": func(cfg *config.ProjectConfig) bool {
				return cfg.Database != config.DatabaseNone
			},
			"frameworkImport": frameworkImport,
			"databaseImport":  databaseImport,
		},
	}
}

// Render renders a template with the given config
func (e *Engine) Render(templatePath string, cfg *config.ProjectConfig) (string, error) {
	// Read template from embedded filesystem
	content, err := protemgen.TemplatesFS.ReadFile("templates/" + templatePath)
	if err != nil {
		return "", fmt.Errorf("template not found: %s: %w", templatePath, err)
	}

	// Parse template
	tmpl, err := template.New(templatePath).Funcs(e.funcMap).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// frameworkImport returns the import path for the selected framework
func frameworkImport(framework config.Framework) string {
	imports := map[config.Framework]string{
		config.FrameworkGin:   "github.com/gin-gonic/gin",
		config.FrameworkEcho:  "github.com/labstack/echo/v4",
		config.FrameworkChi:   "github.com/go-chi/chi/v5",
		config.FrameworkFiber: "github.com/gofiber/fiber/v2",
	}
	return imports[framework]
}

// databaseImport returns the import path for the selected database driver
func databaseImport(database config.Database) string {
	imports := map[config.Database]string{
		config.DatabasePostgres: "github.com/jackc/pgx/v5",
		config.DatabaseMySQL:    "github.com/go-sql-driver/mysql",
		config.DatabaseSQLite:   "modernc.org/sqlite",
	}
	return imports[database]
}
