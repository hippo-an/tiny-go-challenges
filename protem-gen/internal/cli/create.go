package cli

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/config"
	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/generator"
	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/prompt"
	"github.com/spf13/cobra"
)

var (
	flagName           string
	flagModule         string
	flagDatabase       string
	flagGRPC           bool
	flagAuth           bool
	flagAI             bool
	flagNonInteractive bool
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Go web application project",
	Long: `Create a new Go web application project with interactive prompts.

This command will guide you through setting up:
  - Project name and Go module path
  - Database (PostgreSQL, SQLite, or none)
  - Optional features (gRPC, Auth boilerplate, AI integration)

The generated project uses Gin as the HTTP framework.

Example:
  protem-gen create
  protem-gen create --name my-app --database postgres`,
	RunE: runCreate,
}

func init() {
	createCmd.Flags().StringVarP(&flagName, "name", "n", "", "Project name")
	createCmd.Flags().StringVarP(&flagModule, "module", "m", "", "Go module path (e.g., github.com/user/project)")
	createCmd.Flags().StringVarP(&flagDatabase, "database", "d", "postgres", "Database: postgres, sqlite, none")
	createCmd.Flags().BoolVar(&flagGRPC, "grpc", false, "Include gRPC support")
	createCmd.Flags().BoolVar(&flagAuth, "auth", false, "Include authentication boilerplate")
	createCmd.Flags().BoolVar(&flagAI, "ai", false, "Include AI integration boilerplate")
	createCmd.Flags().BoolVar(&flagNonInteractive, "non-interactive", false, "Run without interactive prompts")
}

func runCreate(cmd *cobra.Command, args []string) error {
	var cfg *config.ProjectConfig

	if flagNonInteractive {
		// Non-interactive mode: use flags directly
		cfg = buildConfigFromFlags()
	} else {
		// Interactive mode: run Bubbletea TUI
		var err error
		cfg, err = runInteractivePrompt()
		if err != nil {
			return fmt.Errorf("prompt error: %w", err)
		}
		if cfg == nil {
			// User cancelled
			fmt.Println("Project creation cancelled.")
			return nil
		}
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Set output directory
	if cfg.OutputDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		cfg.OutputDir = filepath.Join(cwd, cfg.Name)
	}

	// Generate project
	gen := generator.New(cfg)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	printSuccess(cfg)
	return nil
}

func buildConfigFromFlags() *config.ProjectConfig {
	cfg := config.NewDefaultConfig()
	cfg.Name = flagName
	cfg.ModulePath = flagModule
	cfg.Database = config.Database(flagDatabase)
	cfg.IncludeGRPC = flagGRPC
	cfg.IncludeAuth = flagAuth
	cfg.IncludeAI = flagAI
	return cfg
}

func runInteractivePrompt() (*config.ProjectConfig, error) {
	m := prompt.NewModel()
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(prompt.Model)
	if result.Cancelled {
		return nil, nil
	}

	return result.Config, nil
}

func printSuccess(cfg *config.ProjectConfig) {
	fmt.Println()
	fmt.Printf("✓ Project '%s' created successfully!\n", cfg.Name)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", cfg.Name)
	fmt.Println("  make setup      # Install dependencies")
	fmt.Println("  make dev        # Start development server")
	fmt.Println()
	fmt.Printf("  Open http://localhost:8080 in your browser\n")
	fmt.Println()
}
