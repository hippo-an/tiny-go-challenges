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
			"title":       strings.Title,
			"lower":       strings.ToLower,
			"upper":       strings.ToUpper,
			"replace":     strings.ReplaceAll,
			"contains":    strings.Contains,
			"ToSnakeCase": toSnakeCase,
			"ToCamelCase": toCamelCase,
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
			"databaseImport": databaseImport,
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

// databaseImport returns the import path for the selected database driver
func databaseImport(database config.Database) string {
	imports := map[config.Database]string{
		config.DatabasePostgres: "github.com/jackc/pgx/v5",
		config.DatabaseSQLite:   "modernc.org/sqlite",
	}
	return imports[database]
}

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	// Replace hyphens with underscores
	s = strings.ReplaceAll(s, "-", "_")

	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// toCamelCase converts a string to CamelCase
func toCamelCase(s string) string {
	// Replace hyphens and underscores with spaces for splitting
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}
