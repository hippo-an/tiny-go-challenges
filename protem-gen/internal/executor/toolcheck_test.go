package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequiredTools(t *testing.T) {
	tools := RequiredTools()

	require.Len(t, tools, 4, "should have 4 required tools")

	// Verify expected tools are present
	toolNames := make([]string, len(tools))
	for i, tool := range tools {
		toolNames[i] = tool.Name
	}

	assert.Contains(t, toolNames, "go")
	assert.Contains(t, toolNames, "npm")
	assert.Contains(t, toolNames, "air")
	assert.Contains(t, toolNames, "templ")

	// Verify all tools have required fields
	for _, tool := range tools {
		assert.NotEmpty(t, tool.Name, "tool name should not be empty")
		assert.NotEmpty(t, tool.CheckCmd, "check command should not be empty")
		assert.NotEmpty(t, tool.InstallHelp, "install help should not be empty")
		assert.True(t, tool.Required, "all tools should be required")
	}
}

func TestToolInfo(t *testing.T) {
	tool := ToolInfo{
		Name:        "test-tool",
		CheckCmd:    "test",
		CheckArgs:   []string{"--version"},
		InstallHelp: "Install test-tool",
		Required:    true,
	}

	assert.Equal(t, "test-tool", tool.Name)
	assert.Equal(t, "test", tool.CheckCmd)
	assert.Equal(t, []string{"--version"}, tool.CheckArgs)
	assert.Equal(t, "Install test-tool", tool.InstallHelp)
	assert.True(t, tool.Required)
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "go version",
			input:    "go version go1.21.0 darwin/arm64",
			expected: "go version go1.21.0 darwin/arm64",
		},
		{
			name:     "npm version",
			input:    "10.2.0",
			expected: "10.2.0",
		},
		{
			name:     "multi-line output",
			input:    "v1.0.0\nsome other info\nmore info",
			expected: "v1.0.0",
		},
		{
			name:     "long version truncated",
			input:    "This is a very long version string that exceeds fifty characters in total length",
			expected: "This is a very long version string that exceeds fi...",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace trimmed",
			input:    "  v1.0.0  \n",
			expected: "v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractVersion(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToolCheckResult(t *testing.T) {
	tool := ToolInfo{
		Name:     "test",
		Required: true,
	}

	result := ToolCheckResult{
		Tool:      tool,
		Available: true,
		Version:   "1.0.0",
		Error:     nil,
	}

	assert.Equal(t, "test", result.Tool.Name)
	assert.True(t, result.Available)
	assert.Equal(t, "1.0.0", result.Version)
	assert.NoError(t, result.Error)
}

func TestMissingToolsError(t *testing.T) {
	results := []ToolCheckResult{
		{
			Tool: ToolInfo{
				Name:        "go",
				InstallHelp: "Install Go",
				Required:    true,
			},
			Available: true,
			Version:   "1.21.0",
		},
		{
			Tool: ToolInfo{
				Name:        "missing-tool",
				InstallHelp: "Install missing-tool from example.com",
				Required:    true,
			},
			Available: false,
		},
	}

	err := &MissingToolsError{
		Tools:   []string{"missing-tool"},
		Results: results,
	}

	errorMsg := err.Error()
	assert.Contains(t, errorMsg, "Missing required tools")
	assert.Contains(t, errorMsg, "missing-tool")
	assert.Contains(t, errorMsg, "not found")
	assert.Contains(t, errorMsg, "Install missing-tool from example.com")
	assert.NotContains(t, errorMsg, "go") // Available tool should not be in error
}
