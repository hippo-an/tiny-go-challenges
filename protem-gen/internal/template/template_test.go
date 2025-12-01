package template

import (
	"testing"

	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "CamelCase to snake_case",
			input:    "MyProject",
			expected: "my_project",
		},
		{
			name:     "hyphenated to snake_case",
			input:    "my-app",
			expected: "my_app",
		},
		{
			name:     "acronym handling - each uppercase gets underscore",
			input:    "APIServer",
			expected: "a_p_i_server",
		},
		{
			name:     "already snake_case",
			input:    "my_project",
			expected: "my_project",
		},
		{
			name:     "lowercase",
			input:    "myapp",
			expected: "myapp",
		},
		{
			name:     "single word uppercase - each char gets underscore",
			input:    "APP",
			expected: "a_p_p",
		},
		{
			name:     "mixed hyphens and camel",
			input:    "My-APIServer",
			expected: "my__a_p_i_server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "hyphenated to CamelCase",
			input:    "my-app",
			expected: "MyApp",
		},
		{
			name:     "snake_case to CamelCase",
			input:    "hello_world",
			expected: "HelloWorld",
		},
		{
			name:     "lowercase to CamelCase",
			input:    "myapp",
			expected: "Myapp",
		},
		{
			name:     "space separated",
			input:    "hello world",
			expected: "HelloWorld",
		},
		{
			name:     "mixed separators",
			input:    "hello-world_test",
			expected: "HelloWorldTest",
		},
		{
			name:     "uppercase input",
			input:    "HELLO",
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toCamelCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDatabaseImport(t *testing.T) {
	tests := []struct {
		name     string
		database config.Database
		expected string
	}{
		{
			name:     "postgres import",
			database: config.DatabasePostgres,
			expected: "github.com/jackc/pgx/v5",
		},
		{
			name:     "sqlite import",
			database: config.DatabaseSQLite,
			expected: "modernc.org/sqlite",
		},
		{
			name:     "none returns empty",
			database: config.DatabaseNone,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := databaseImport(tt.database)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewEngine(t *testing.T) {
	engine := NewEngine()

	require.NotNil(t, engine)
	require.NotNil(t, engine.funcMap)

	// Verify expected functions are registered
	expectedFuncs := []string{
		"title", "lower", "upper", "replace", "contains",
		"ToSnakeCase", "ToCamelCase", "hasGRPC", "hasAuth", "hasAI", "hasDB",
		"databaseImport",
	}

	for _, fn := range expectedFuncs {
		assert.NotNil(t, engine.funcMap[fn], "function %s should be registered", fn)
	}
}

func TestEngine_Render(t *testing.T) {
	engine := NewEngine()

	t.Run("render existing template", func(t *testing.T) {
		cfg := &config.ProjectConfig{
			Name:       "test-project",
			ModulePath: "github.com/user/test-project",
			Database:   config.DatabasePostgres,
		}

		// Test with a known template
		result, err := engine.Render("base/.gitignore.tmpl", cfg)
		require.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Contains(t, result, "node_modules")
		assert.Contains(t, result, "tmp/")
	})

	t.Run("non-existent template returns error", func(t *testing.T) {
		cfg := &config.ProjectConfig{}

		_, err := engine.Render("non-existent.tmpl", cfg)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "template not found")
	})
}

func TestHelperFunctions(t *testing.T) {
	engine := NewEngine()

	t.Run("hasGRPC function", func(t *testing.T) {
		fn := engine.funcMap["hasGRPC"].(func(*config.ProjectConfig) bool)

		cfg := &config.ProjectConfig{IncludeGRPC: true}
		assert.True(t, fn(cfg))

		cfg.IncludeGRPC = false
		assert.False(t, fn(cfg))
	})

	t.Run("hasAuth function", func(t *testing.T) {
		fn := engine.funcMap["hasAuth"].(func(*config.ProjectConfig) bool)

		cfg := &config.ProjectConfig{IncludeAuth: true}
		assert.True(t, fn(cfg))

		cfg.IncludeAuth = false
		assert.False(t, fn(cfg))
	})

	t.Run("hasAI function", func(t *testing.T) {
		fn := engine.funcMap["hasAI"].(func(*config.ProjectConfig) bool)

		cfg := &config.ProjectConfig{IncludeAI: true}
		assert.True(t, fn(cfg))

		cfg.IncludeAI = false
		assert.False(t, fn(cfg))
	})

	t.Run("hasDB function", func(t *testing.T) {
		fn := engine.funcMap["hasDB"].(func(*config.ProjectConfig) bool)

		cfg := &config.ProjectConfig{Database: config.DatabasePostgres}
		assert.True(t, fn(cfg))

		cfg.Database = config.DatabaseNone
		assert.False(t, fn(cfg))
	})
}
