package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/config"
	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/executor"
	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/template"
)

// Generator handles project generation
type Generator struct {
	config   *config.ProjectConfig
	template *template.Engine
	executor *executor.Executor
	modifier *executor.Modifier
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fmt.Printf("Creating project: %s\n\n", g.config.Name)

	// Step 1: Check required tools
	fmt.Println("Checking required tools...")
	tmpExecutor := executor.New(".")
	results, err := tmpExecutor.CheckTools(ctx)
	if err != nil {
		executor.PrintToolStatus(results)
		return err
	}
	executor.PrintToolStatus(results)

	// Step 2: Check if output directory exists
	if _, err := os.Stat(g.config.OutputDir); !os.IsNotExist(err) {
		return config.ErrOutputDirExists
	}

	// Step 3: Create directory structure
	fmt.Println("Creating directory structure...")
	if err := g.createDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Initialize executor and modifier with project directory
	g.executor = executor.New(g.config.OutputDir)
	g.modifier = executor.NewModifier(g.config.OutputDir)

	// Step 4: Execute CLI commands
	fmt.Println("Initializing project...")
	if err := g.executeCLICommands(ctx); err != nil {
		g.cleanup()
		return fmt.Errorf("failed to execute CLI commands: %w", err)
	}

	// Step 5: Modify CLI-generated files
	fmt.Println("Configuring project files...")
	if err := g.modifyCLIGeneratedFiles(); err != nil {
		return fmt.Errorf("failed to modify generated files: %w", err)
	}

	// Step 6: Generate template-based files
	fmt.Println("Generating source files...")
	if err := g.generateTemplateFiles(); err != nil {
		return fmt.Errorf("failed to generate template files: %w", err)
	}

	// Step 7: Install Go dependencies and tidy
	fmt.Println("Installing dependencies...")
	if err := g.installGoDependencies(ctx); err != nil {
		return fmt.Errorf("failed to install Go dependencies: %w", err)
	}

	fmt.Println()
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

func (g *Generator) executeCLICommands(ctx context.Context) error {
	// 1. go mod init
	fmt.Printf("  go mod init %s\n", g.config.ModulePath)
	if err := g.executor.GoModInit(ctx, g.config.ModulePath); err != nil {
		return err
	}

	// 2. npm init
	fmt.Println("  npm init")
	if err := g.executor.NpmInit(ctx); err != nil {
		return err
	}

	// 3. npm install tailwindcss packages
	packages := []string{"tailwindcss", "@tailwindcss/forms", "@tailwindcss/typography"}
	fmt.Printf("  npm install %s\n", packages)
	if err := g.executor.NpmInstall(ctx, packages, true); err != nil {
		return err
	}

	// 4. air init
	fmt.Println("  air init")
	if err := g.executor.AirInit(ctx); err != nil {
		return err
	}

	// 5. tailwindcss init
	fmt.Println("  npx tailwindcss init")
	if err := g.executor.TailwindInit(ctx); err != nil {
		return err
	}

	return nil
}

func (g *Generator) modifyCLIGeneratedFiles() error {
	// Modify air.toml
	if err := g.modifier.ModifyAirConfig(); err != nil {
		fmt.Printf("  Warning: could not modify .air.toml: %v\n", err)
	}

	// Modify and move tailwind.config.js
	if err := g.modifier.ModifyTailwindConfig(); err != nil {
		fmt.Printf("  Warning: could not modify tailwind.config.js: %v\n", err)
	}
	if err := g.modifier.MoveTailwindConfig(); err != nil {
		fmt.Printf("  Warning: could not move tailwind.config.js: %v\n", err)
	}

	// Modify package.json
	if err := g.modifier.ModifyPackageJSON(); err != nil {
		fmt.Printf("  Warning: could not modify package.json: %v\n", err)
	}

	return nil
}

func (g *Generator) generateTemplateFiles() error {
	// Generate remaining template-based files
	files := []struct {
		template string
		output   string
	}{
		{"base/Makefile.tmpl", "Makefile"},
		{"base/.gitignore.tmpl", ".gitignore"},
		{"base/README.md.tmpl", "README.md"},
		{"base/cmd/server/main.go.tmpl", "cmd/server/main.go"},
	}

	for _, f := range files {
		if err := g.generateFile(f.template, f.output); err != nil {
			return fmt.Errorf("failed to generate %s: %w", f.output, err)
		}
	}

	// Generate framework files (Gin)
	if err := g.generateFrameworkFiles(); err != nil {
		return err
	}

	// Generate database-specific files
	if err := g.generateDatabaseFiles(); err != nil {
		return err
	}

	// Generate frontend template files
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
	// Gin is the only supported framework
	files := []struct {
		template string
		output   string
	}{
		{"http/gin/server.go.tmpl", "internal/infrastructure/http/server.go"},
		{"http/gin/routes.go.tmpl", "internal/interfaces/http/routes.go"},
		{"http/gin/handler.go.tmpl", "internal/interfaces/http/handler.go"},
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
	// Only generate template-based frontend files
	// package.json and tailwind.config.js are now CLI-generated
	files := []struct {
		template string
		output   string
	}{
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

func (g *Generator) installGoDependencies(ctx context.Context) error {
	packages := []string{
		"github.com/gin-gonic/gin",
		"github.com/a-h/templ",
	}

	// Add database driver
	switch g.config.Database {
	case config.DatabasePostgres:
		packages = append(packages, "github.com/jackc/pgx/v5")
	case config.DatabaseMySQL:
		packages = append(packages, "github.com/go-sql-driver/mysql")
	case config.DatabaseSQLite:
		packages = append(packages, "modernc.org/sqlite")
	}

	fmt.Printf("  go get %v\n", packages)
	if err := g.executor.GoGet(ctx, packages); err != nil {
		return err
	}

	fmt.Println("  go mod tidy")
	return g.executor.GoModTidy(ctx)
}

func (g *Generator) cleanup() {
	// Remove partially created project on error
	if g.config.OutputDir != "" {
		os.RemoveAll(g.config.OutputDir)
	}
}
