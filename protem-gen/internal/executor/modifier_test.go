package executor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModifier(t *testing.T) {
	m := NewModifier("/tmp/test-project")

	require.NotNil(t, m)
	assert.Equal(t, "/tmp/test-project", m.projectDir)
}

func TestModifyAirConfig(t *testing.T) {
	t.Run("modifies air config correctly", func(t *testing.T) {
		tmpDir := t.TempDir()
		m := NewModifier(tmpDir)

		// Create initial .air.toml (simulating air init output)
		initialConfig := `root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  post_cmd = []
  pre_cmd = []
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[proxy]
  app_port = 0
  enabled = false
  proxy_port = 0

[screen]
  clear_on_rebuild = false
  keep_scroll = true
`
		airPath := filepath.Join(tmpDir, ".air.toml")
		err := os.WriteFile(airPath, []byte(initialConfig), 0644)
		require.NoError(t, err)

		// Run modification
		err = m.ModifyAirConfig()
		require.NoError(t, err)

		// Read modified content
		modified, err := os.ReadFile(airPath)
		require.NoError(t, err)

		content := string(modified)

		// Verify modifications
		assert.Contains(t, content, `cmd = "go build -o ./tmp/main ./cmd/server"`,
			"build cmd should be modified")
		assert.Contains(t, content, `pre_cmd = ["templ generate"]`,
			"pre_cmd should include templ generate")
		assert.Contains(t, content, `include_ext = ["go", "tpl", "tmpl", "html", "templ"]`,
			"include_ext should include templ")
		assert.Contains(t, content, `exclude_regex = ["_test.go", "_templ.go"`,
			"exclude_regex should include _templ.go")
	})

	t.Run("returns error when file not found", func(t *testing.T) {
		m := NewModifier("/nonexistent/path")

		err := m.ModifyAirConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read .air.toml")
	})
}

func TestModifyPackageJSON(t *testing.T) {
	t.Run("adds scripts to package.json", func(t *testing.T) {
		tmpDir := t.TempDir()
		m := NewModifier(tmpDir)

		// Create initial package.json (simulating npm init -y output)
		initialPkg := map[string]interface{}{
			"name":        "test-project",
			"version":     "1.0.0",
			"description": "",
			"main":        "index.js",
			"scripts": map[string]interface{}{
				"test": "echo \"Error: no test specified\" && exit 1",
			},
			"keywords": []string{},
			"author":   "",
			"license":  "ISC",
		}

		content, err := json.MarshalIndent(initialPkg, "", "  ")
		require.NoError(t, err)

		pkgPath := filepath.Join(tmpDir, "package.json")
		err = os.WriteFile(pkgPath, content, 0644)
		require.NoError(t, err)

		// Run modification
		err = m.ModifyPackageJSON()
		require.NoError(t, err)

		// Read modified content
		modified, err := os.ReadFile(pkgPath)
		require.NoError(t, err)

		var pkg map[string]interface{}
		err = json.Unmarshal(modified, &pkg)
		require.NoError(t, err)

		scripts := pkg["scripts"].(map[string]interface{})

		// Verify scripts are added
		assert.Contains(t, scripts["build:css"], "@tailwindcss/cli")
		assert.Contains(t, scripts["watch:css"], "@tailwindcss/cli")
		assert.Contains(t, scripts["watch:css"], "--watch")
	})

	t.Run("creates scripts if not present", func(t *testing.T) {
		tmpDir := t.TempDir()
		m := NewModifier(tmpDir)

		// Create package.json without scripts
		initialPkg := map[string]interface{}{
			"name":    "test-project",
			"version": "1.0.0",
		}

		content, err := json.MarshalIndent(initialPkg, "", "  ")
		require.NoError(t, err)

		pkgPath := filepath.Join(tmpDir, "package.json")
		err = os.WriteFile(pkgPath, content, 0644)
		require.NoError(t, err)

		// Run modification
		err = m.ModifyPackageJSON()
		require.NoError(t, err)

		// Read modified content
		modified, err := os.ReadFile(pkgPath)
		require.NoError(t, err)

		var pkg map[string]interface{}
		err = json.Unmarshal(modified, &pkg)
		require.NoError(t, err)

		scripts := pkg["scripts"].(map[string]interface{})
		assert.NotNil(t, scripts["build:css"])
		assert.NotNil(t, scripts["watch:css"])
	})

	t.Run("returns error when file not found", func(t *testing.T) {
		m := NewModifier("/nonexistent/path")

		err := m.ModifyPackageJSON()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read package.json")
	})

	t.Run("returns error for invalid json", func(t *testing.T) {
		tmpDir := t.TempDir()
		m := NewModifier(tmpDir)

		pkgPath := filepath.Join(tmpDir, "package.json")
		err := os.WriteFile(pkgPath, []byte("invalid json"), 0644)
		require.NoError(t, err)

		err = m.ModifyPackageJSON()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse package.json")
	})
}
