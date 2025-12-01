package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cfg := &config.ProjectConfig{
		Name:       "test-project",
		ModulePath: "github.com/user/test-project",
		Database:   config.DatabasePostgres,
	}

	gen := New(cfg)

	require.NotNil(t, gen)
	assert.Equal(t, cfg, gen.config)
	assert.NotNil(t, gen.template)
	assert.Nil(t, gen.executor, "executor should be nil before Generate()")
	assert.Nil(t, gen.modifier, "modifier should be nil before Generate()")
}

func TestGenerator_createDirectories(t *testing.T) {
	t.Run("creates base directories", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputDir := filepath.Join(tmpDir, "test-project")

		cfg := &config.ProjectConfig{
			Name:       "test-project",
			ModulePath: "github.com/user/test-project",
			OutputDir:  outputDir,
			Database:   config.DatabasePostgres,
		}

		gen := New(cfg)
		err := gen.createDirectories()
		require.NoError(t, err)

		// Verify base directories exist
		expectedDirs := []string{
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

		for _, dir := range expectedDirs {
			path := filepath.Join(outputDir, dir)
			info, err := os.Stat(path)
			require.NoError(t, err, "directory %s should exist", dir)
			assert.True(t, info.IsDir(), "%s should be a directory", dir)
		}
	})

	t.Run("creates gRPC directories when enabled", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputDir := filepath.Join(tmpDir, "test-project-grpc")

		cfg := &config.ProjectConfig{
			Name:        "test-project",
			ModulePath:  "github.com/user/test-project",
			OutputDir:   outputDir,
			Database:    config.DatabasePostgres,
			IncludeGRPC: true,
		}

		gen := New(cfg)
		err := gen.createDirectories()
		require.NoError(t, err)

		// Verify gRPC directories exist
		grpcDirs := []string{
			"internal/interfaces/grpc",
			"proto",
		}

		for _, dir := range grpcDirs {
			path := filepath.Join(outputDir, dir)
			info, err := os.Stat(path)
			require.NoError(t, err, "gRPC directory %s should exist", dir)
			assert.True(t, info.IsDir())
		}
	})

	t.Run("creates AI directories when enabled", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputDir := filepath.Join(tmpDir, "test-project-ai")

		cfg := &config.ProjectConfig{
			Name:       "test-project",
			ModulePath: "github.com/user/test-project",
			OutputDir:  outputDir,
			Database:   config.DatabasePostgres,
			IncludeAI:  true,
		}

		gen := New(cfg)
		err := gen.createDirectories()
		require.NoError(t, err)

		// Verify AI directories exist
		aiDirs := []string{
			"internal/infrastructure/llm",
			"internal/infrastructure/prompt",
			"internal/infrastructure/stream",
		}

		for _, dir := range aiDirs {
			path := filepath.Join(outputDir, dir)
			info, err := os.Stat(path)
			require.NoError(t, err, "AI directory %s should exist", dir)
			assert.True(t, info.IsDir())
		}
	})

	t.Run("creates Auth directories when enabled", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputDir := filepath.Join(tmpDir, "test-project-auth")

		cfg := &config.ProjectConfig{
			Name:        "test-project",
			ModulePath:  "github.com/user/test-project",
			OutputDir:   outputDir,
			Database:    config.DatabasePostgres,
			IncludeAuth: true,
		}

		gen := New(cfg)
		err := gen.createDirectories()
		require.NoError(t, err)

		// Verify Auth directory exists
		authPath := filepath.Join(outputDir, "internal/infrastructure/auth")
		info, err := os.Stat(authPath)
		require.NoError(t, err, "auth directory should exist")
		assert.True(t, info.IsDir())
	})

	t.Run("creates all optional directories when all options enabled", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputDir := filepath.Join(tmpDir, "test-project-all")

		cfg := &config.ProjectConfig{
			Name:        "test-project",
			ModulePath:  "github.com/user/test-project",
			OutputDir:   outputDir,
			Database:    config.DatabasePostgres,
			IncludeGRPC: true,
			IncludeAI:   true,
			IncludeAuth: true,
		}

		gen := New(cfg)
		err := gen.createDirectories()
		require.NoError(t, err)

		// Verify all optional directories exist
		allDirs := []string{
			"internal/interfaces/grpc",
			"proto",
			"internal/infrastructure/llm",
			"internal/infrastructure/prompt",
			"internal/infrastructure/stream",
			"internal/infrastructure/auth",
		}

		for _, dir := range allDirs {
			path := filepath.Join(outputDir, dir)
			_, err := os.Stat(path)
			require.NoError(t, err, "directory %s should exist", dir)
		}
	})
}

func TestGenerator_generateFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "test-project")
	err := os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)

	cfg := &config.ProjectConfig{
		Name:       "test-project",
		ModulePath: "github.com/user/test-project",
		OutputDir:  outputDir,
		Database:   config.DatabasePostgres,
	}

	gen := New(cfg)

	t.Run("generates file from template", func(t *testing.T) {
		err := gen.generateFile("base/.gitignore.tmpl", ".gitignore")
		require.NoError(t, err)

		// Verify file was created
		content, err := os.ReadFile(filepath.Join(outputDir, ".gitignore"))
		require.NoError(t, err)
		assert.Contains(t, string(content), "node_modules")
	})

	t.Run("creates parent directories if needed", func(t *testing.T) {
		err := gen.generateFile("base/.gitignore.tmpl", "nested/deep/.gitignore")
		require.NoError(t, err)

		// Verify file was created in nested directory
		_, err = os.Stat(filepath.Join(outputDir, "nested/deep/.gitignore"))
		require.NoError(t, err)
	})

	t.Run("returns error for non-existent template", func(t *testing.T) {
		err := gen.generateFile("nonexistent.tmpl", "output.txt")
		assert.Error(t, err)
	})
}

func TestGenerator_cleanup(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "test-project")

	// Create a directory to cleanup
	err := os.MkdirAll(filepath.Join(outputDir, "some/nested/dir"), 0755)
	require.NoError(t, err)

	cfg := &config.ProjectConfig{
		Name:      "test-project",
		OutputDir: outputDir,
	}

	gen := New(cfg)
	gen.cleanup()

	// Verify directory was removed
	_, err = os.Stat(outputDir)
	assert.True(t, os.IsNotExist(err), "output directory should be removed")
}
