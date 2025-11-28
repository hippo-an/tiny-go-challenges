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

	// Add CSS build scripts (Tailwind CSS v4 uses @tailwindcss/cli)
	scripts["build:css"] = "npx @tailwindcss/cli -i ./web/tailwind/input.css -o ./web/static/css/output.css --minify"
	scripts["watch:css"] = "npx @tailwindcss/cli -i ./web/tailwind/input.css -o ./web/static/css/output.css --watch"
	pkg["scripts"] = scripts

	// Write back with indentation
	output, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal package.json: %w", err)
	}

	return os.WriteFile(pkgPath, output, 0644)
}

