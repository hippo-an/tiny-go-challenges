package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/config"
	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/template"
)

// Generator handles project generation
type Generator struct {
	config   *config.ProjectConfig
	template *template.Engine
}

// New creates a new Generator instance
func New(cfg *config.ProjectConfig) *Generator {
	return &Generator{
		config:   cfg,
		template: template.NewEngine(),
	}
}

// Generate creates the project structure and files
func (g *Generator) Generate() error {
	fmt.Printf("Creating project: %s\n", g.config.Name)

	// Check if output directory exists
	if _, err := os.Stat(g.config.OutputDir); !os.IsNotExist(err) {
		return config.ErrOutputDirExists
	}

	// Create directory structure
	if err := g.createDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Generate files from templates
	if err := g.generateFiles(); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	fmt.Printf("✓ Created directory structure\n")
	fmt.Printf("✓ Generated configuration files\n")

	return nil
}

func (g *Generator) createDirectories() error {
	dirs := []string{
		"cmd/server",
		"internal/domain",
		"internal/application",
		"internal/infrastructure/database",
		"internal/infrastructure/http",
		"internal/interfaces/http",
		"pkg",
		"web/templates/layouts",
		"web/templates/pages",
		"web/templates/components",
		"web/static/css",
		"web/static/js",
		"web/tailwind",
		"migrations",
		"sqlc/queries",
	}

	// Add gRPC directories if enabled
	if g.config.IncludeGRPC {
		dirs = append(dirs, "internal/interfaces/grpc", "proto")
	}

	// Add AI directories if enabled
	if g.config.IncludeAI {
		dirs = append(dirs, "internal/infrastructure/llm", "internal/infrastructure/prompt")
	}

	for _, dir := range dirs {
		path := filepath.Join(g.config.OutputDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func (g *Generator) generateFiles() error {
	// Generate base files
	files := []struct {
		template string
		output   string
	}{
		{"base/go.mod.tmpl", "go.mod"},
		{"base/Makefile.tmpl", "Makefile"},
		{"base/.gitignore.tmpl", ".gitignore"},
		{"base/.air.toml.tmpl", ".air.toml"},
		{"base/README.md.tmpl", "README.md"},
		{"base/cmd/server/main.go.tmpl", "cmd/server/main.go"},
	}

	for _, f := range files {
		if err := g.generateFile(f.template, f.output); err != nil {
			return fmt.Errorf("failed to generate %s: %w", f.output, err)
		}
	}

	// Generate framework-specific files
	if err := g.generateFrameworkFiles(); err != nil {
		return err
	}

	// Generate database-specific files
	if err := g.generateDatabaseFiles(); err != nil {
		return err
	}

	// Generate frontend files
	if err := g.generateFrontendFiles(); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateFile(templateName, outputPath string) error {
	content, err := g.template.Render(templateName, g.config)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(g.config.OutputDir, outputPath)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, []byte(content), 0644)
}

func (g *Generator) generateFrameworkFiles() error {
	frameworkDir := fmt.Sprintf("http/%s", g.config.Framework)
	files := []struct {
		template string
		output   string
	}{
		{frameworkDir + "/server.go.tmpl", "internal/infrastructure/http/server.go"},
		{frameworkDir + "/routes.go.tmpl", "internal/interfaces/http/routes.go"},
		{frameworkDir + "/handler.go.tmpl", "internal/interfaces/http/handler.go"},
	}

	for _, f := range files {
		if err := g.generateFile(f.template, f.output); err != nil {
			// Skip if template doesn't exist (graceful fallback)
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func (g *Generator) generateDatabaseFiles() error {
	if g.config.Database == config.DatabaseNone {
		return nil
	}

	dbDir := fmt.Sprintf("database/%s", g.config.Database)
	files := []struct {
		template string
		output   string
	}{
		{dbDir + "/sqlc.yaml.tmpl", "sqlc.yaml"},
		{dbDir + "/db.go.tmpl", "internal/infrastructure/database/db.go"},
		{dbDir + "/schema.sql.tmpl", "migrations/001_init.sql"},
		{dbDir + "/queries.sql.tmpl", "sqlc/queries/queries.sql"},
	}

	for _, f := range files {
		if err := g.generateFile(f.template, f.output); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
	}

	return nil
}

func (g *Generator) generateFrontendFiles() error {
	files := []struct {
		template string
		output   string
	}{
		{"frontend/package.json.tmpl", "package.json"},
		{"frontend/tailwind.config.js.tmpl", "web/tailwind/tailwind.config.js"},
		{"frontend/input.css.tmpl", "web/tailwind/input.css"},
		{"frontend/base.templ.tmpl", "web/templates/layouts/base.templ"},
		{"frontend/index.templ.tmpl", "web/templates/pages/index.templ"},
	}

	for _, f := range files {
		if err := g.generateFile(f.template, f.output); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
	}

	return nil
}
