package executor

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	exec := New("/tmp/test")

	require.NotNil(t, exec)
	assert.Equal(t, "/tmp/test", exec.workDir)
	assert.False(t, exec.verbose)
}

func TestSetVerbose(t *testing.T) {
	exec := New("/tmp")

	assert.False(t, exec.verbose)

	exec.SetVerbose(true)
	assert.True(t, exec.verbose)

	exec.SetVerbose(false)
	assert.False(t, exec.verbose)
}

func TestRun(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exec := New(t.TempDir())

	t.Run("successful command", func(t *testing.T) {
		var result *Result
		if runtime.GOOS == "windows" {
			result = exec.Run(ctx, "cmd", "/c", "echo", "hello")
		} else {
			result = exec.Run(ctx, "echo", "hello")
		}

		require.NoError(t, result.Error)
		assert.Equal(t, 0, result.ExitCode)
		assert.Contains(t, result.Stdout, "hello")
	})

	t.Run("command with args recorded", func(t *testing.T) {
		result := exec.Run(ctx, "echo", "arg1", "arg2")

		assert.Equal(t, "echo", result.Command)
		assert.Equal(t, []string{"arg1", "arg2"}, result.Args)
	})

	t.Run("failing command", func(t *testing.T) {
		result := exec.Run(ctx, "false")

		if runtime.GOOS != "windows" {
			assert.Error(t, result.Error)
			assert.NotEqual(t, 0, result.ExitCode)
		}
	})

	t.Run("non-existent command", func(t *testing.T) {
		result := exec.Run(ctx, "nonexistent_command_xyz123")

		assert.Error(t, result.Error)
	})
}

func TestRunWithEnv(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exec := New(t.TempDir())

	t.Run("with custom environment", func(t *testing.T) {
		env := []string{"TEST_VAR=test_value"}
		result := exec.RunWithEnv(ctx, env, "sh", "-c", "echo $TEST_VAR")

		require.NoError(t, result.Error)
		assert.Contains(t, result.Stdout, "test_value")
	})

	t.Run("without custom environment", func(t *testing.T) {
		result := exec.RunWithEnv(ctx, nil, "echo", "test")

		require.NoError(t, result.Error)
		assert.Contains(t, result.Stdout, "test")
	})
}

func TestResult(t *testing.T) {
	result := &Result{
		Command:  "test",
		Args:     []string{"arg1", "arg2"},
		Stdout:   "output",
		Stderr:   "error",
		ExitCode: 1,
	}

	assert.Equal(t, "test", result.Command)
	assert.Equal(t, []string{"arg1", "arg2"}, result.Args)
	assert.Equal(t, "output", result.Stdout)
	assert.Equal(t, "error", result.Stderr)
	assert.Equal(t, 1, result.ExitCode)
}
