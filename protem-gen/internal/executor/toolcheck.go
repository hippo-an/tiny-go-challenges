package executor

import (
	"context"
	"fmt"
	"strings"
)

// ToolInfo represents a required CLI tool
type ToolInfo struct {
	Name        string
	CheckCmd    string
	CheckArgs   []string
	InstallHelp string
	Required    bool
}

// RequiredTools returns the list of tools needed for project generation
func RequiredTools() []ToolInfo {
	return []ToolInfo{
		{
			Name:        "go",
			CheckCmd:    "go",
			CheckArgs:   []string{"version"},
			InstallHelp: "Install Go from https://go.dev/dl/",
			Required:    true,
		},
		{
			Name:        "npm",
			CheckCmd:    "npm",
			CheckArgs:   []string{"--version"},
			InstallHelp: "Install Node.js from https://nodejs.org/",
			Required:    true,
		},
		{
			Name:        "air",
			CheckCmd:    "air",
			CheckArgs:   []string{"-v"},
			InstallHelp: "Install with: go install github.com/air-verse/air@latest",
			Required:    true,
		},
		{
			Name:        "templ",
			CheckCmd:    "templ",
			CheckArgs:   []string{"version"},
			InstallHelp: "Install with: go install github.com/a-h/templ/cmd/templ@latest",
			Required:    true,
		},
	}
}

// ToolCheckResult represents the result of checking a tool
type ToolCheckResult struct {
	Tool      ToolInfo
	Available bool
	Version   string
	Error     error
}

// CheckTools verifies all required tools are available
func (e *Executor) CheckTools(ctx context.Context) ([]ToolCheckResult, error) {
	tools := RequiredTools()
	results := make([]ToolCheckResult, len(tools))
	var missingTools []string

	for i, tool := range tools {
		result := e.checkTool(ctx, tool)
		results[i] = result

		if !result.Available && tool.Required {
			missingTools = append(missingTools, tool.Name)
		}
	}

	if len(missingTools) > 0 {
		return results, &MissingToolsError{
			Tools:   missingTools,
			Results: results,
		}
	}

	return results, nil
}

func (e *Executor) checkTool(ctx context.Context, tool ToolInfo) ToolCheckResult {
	result := ToolCheckResult{Tool: tool}

	cmdResult := e.Run(ctx, tool.CheckCmd, tool.CheckArgs...)
	if cmdResult.Error != nil {
		result.Available = false
		result.Error = cmdResult.Error
		return result
	}

	result.Available = true
	// Extract version from output
	output := strings.TrimSpace(cmdResult.Stdout)
	if output == "" {
		output = strings.TrimSpace(cmdResult.Stderr)
	}
	result.Version = extractVersion(output)

	return result
}

func extractVersion(output string) string {
	// Try to extract version from common formats
	lines := strings.Split(output, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		// Limit length
		if len(firstLine) > 50 {
			firstLine = firstLine[:50] + "..."
		}
		return firstLine
	}
	return ""
}

// MissingToolsError represents an error when required tools are missing
type MissingToolsError struct {
	Tools   []string
	Results []ToolCheckResult
}

func (e *MissingToolsError) Error() string {
	var sb strings.Builder
	sb.WriteString("Missing required tools:\n\n")

	for _, result := range e.Results {
		if !result.Available && result.Tool.Required {
			sb.WriteString(fmt.Sprintf("  ✗ %s - not found\n", result.Tool.Name))
			sb.WriteString(fmt.Sprintf("    %s\n\n", result.Tool.InstallHelp))
		}
	}

	return sb.String()
}

// PrintToolStatus prints the status of all tools
func PrintToolStatus(results []ToolCheckResult) {
	fmt.Println("Checking required tools...")
	fmt.Println()

	for _, result := range results {
		if result.Available {
			fmt.Printf("  ✓ %s (%s)\n", result.Tool.Name, result.Version)
		} else {
			fmt.Printf("  ✗ %s - not found\n", result.Tool.Name)
		}
	}
	fmt.Println()
}
