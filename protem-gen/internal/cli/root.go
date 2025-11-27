package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "protem-gen",
	Short: "Go web application template generator",
	Long: `protem-gen is a CLI tool that generates templated Go web application projects
with hot-reloading development environment and modern architecture patterns.

It scaffolds new projects with:
  - HTTP framework (Gin, Echo, Chi, or Fiber)
  - Database integration (PostgreSQL, MySQL, SQLite) with sqlc
  - Frontend stack (Tailwind CSS, Alpine.js, htmx, templ)
  - Hot-reloading development environment (air, templ watch, tailwind watch)
  - Clean Architecture project structure`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(createCmd)
}
