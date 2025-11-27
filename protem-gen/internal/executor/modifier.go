package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Modifier handles post-CLI file modifications
type Modifier struct {
	projectDir string
}

// NewModifier creates a new Modifier
func NewModifier(projectDir string) *Modifier {
	return &Modifier{
		projectDir: projectDir,
	}
}

// ModifyAirConfig modifies the air.toml configuration
func (m *Modifier) ModifyAirConfig() error {
	airPath := filepath.Join(m.projectDir, ".air.toml")

	content, err := os.ReadFile(airPath)
	if err != nil {
		return fmt.Errorf("failed to read .air.toml: %w", err)
	}

	config := string(content)

	// Modify build command
	config = strings.Replace(config,
		`cmd = "go build -o ./tmp/main ."`,
		`cmd = "go build -o ./tmp/main ./cmd/server"`,
		1)

	// Add templ files to include_ext if not present
	if !strings.Contains(config, "templ") {
		config = strings.Replace(config,
			`include_ext = ["go", "tpl", "tmpl", "html"]`,
			`include_ext = ["go", "tpl", "tmpl", "html", "templ"]`,
			1)
	}

	// Add _templ.go to exclude_regex
	if strings.Contains(config, `exclude_regex = ["_test.go"`) {
		config = strings.Replace(config,
			`exclude_regex = ["_test.go"`,
			`exclude_regex = ["_test.go", "_templ.go"`,
			1)
	}

	// Add pre_cmd for templ generate
	if !strings.Contains(config, "templ generate") {
		// Find [build] section and add pre_cmd after cmd line
		lines := strings.Split(config, "\n")
		var newLines []string
		for i, line := range lines {
			newLines = append(newLines, line)
			if strings.HasPrefix(strings.TrimSpace(line), "cmd = ") && i > 0 {
				// Add pre_cmd after cmd
				newLines = append(newLines, `  pre_cmd = ["templ generate"]`)
			}
		}
		config = strings.Join(newLines, "\n")
	}

	return os.WriteFile(airPath, []byte(config), 0644)
}

// ModifyTailwindConfig modifies the tailwind.config.js configuration
func (m *Modifier) ModifyTailwindConfig() error {
	tailwindPath := filepath.Join(m.projectDir, "tailwind.config.js")

	content, err := os.ReadFile(tailwindPath)
	if err != nil {
		return fmt.Errorf("failed to read tailwind.config.js: %w", err)
	}

	config := string(content)

	// Add content paths
	config = strings.Replace(config,
		`content: []`,
		`content: [
    "./web/templates/**/*.templ",
    "./web/templates/**/*_templ.go",
  ]`,
		1)

	// Add plugins
	config = strings.Replace(config,
		`plugins: []`,
		`plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]`,
		1)

	return os.WriteFile(tailwindPath, []byte(config), 0644)
}

// ModifyPackageJSON modifies the package.json to add scripts
func (m *Modifier) ModifyPackageJSON() error {
	pkgPath := filepath.Join(m.projectDir, "package.json")

	content, err := os.ReadFile(pkgPath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return fmt.Errorf("failed to parse package.json: %w", err)
	}

	// Get or create scripts section
	scripts, ok := pkg["scripts"].(map[string]interface{})
	if !ok {
		scripts = make(map[string]interface{})
	}

	// Add CSS build scripts
	scripts["build:css"] = "tailwindcss -i ./web/tailwind/input.css -o ./web/static/css/output.css --minify"
	scripts["watch:css"] = "tailwindcss -i ./web/tailwind/input.css -o ./web/static/css/output.css --watch"
	pkg["scripts"] = scripts

	// Write back with indentation
	output, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal package.json: %w", err)
	}

	return os.WriteFile(pkgPath, output, 0644)
}

// MoveTailwindConfig moves tailwind.config.js to web/tailwind/
func (m *Modifier) MoveTailwindConfig() error {
	srcPath := filepath.Join(m.projectDir, "tailwind.config.js")
	dstPath := filepath.Join(m.projectDir, "web/tailwind/tailwind.config.js")

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	// Read source file
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// Write to destination
	if err := os.WriteFile(dstPath, content, 0644); err != nil {
		return err
	}

	// Remove source file
	return os.Remove(srcPath)
}
