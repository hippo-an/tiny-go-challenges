package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Executor handles CLI command execution
type Executor struct {
	workDir string
	verbose bool
}

// Result represents the outcome of a command execution
type Result struct {
	Command  string
	Args     []string
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
}

// New creates a new Executor
func New(workDir string) *Executor {
	return &Executor{
		workDir: workDir,
		verbose: false,
	}
}

// SetVerbose enables verbose output
func (e *Executor) SetVerbose(v bool) {
	e.verbose = v
}

// Run executes a command and returns the result
func (e *Executor) Run(ctx context.Context, cmd string, args ...string) *Result {
	return e.RunWithEnv(ctx, nil, cmd, args...)
}

// RunWithEnv executes with custom environment variables
func (e *Executor) RunWithEnv(ctx context.Context, env []string, cmd string, args ...string) *Result {
	result := &Result{
		Command: cmd,
		Args:    args,
	}

	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		fullArgs := append([]string{"/c", cmd}, args...)
		command = exec.CommandContext(ctx, "cmd", fullArgs...)
	} else {
		command = exec.CommandContext(ctx, cmd, args...)
	}

	command.Dir = e.workDir

	// Set environment
	command.Env = os.Environ()
	if env != nil {
		command.Env = append(command.Env, env...)
	}

	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if e.verbose {
		fmt.Printf("  Running: %s %s\n", cmd, strings.Join(args, " "))
	}

	err := command.Run()
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
		result.Error = err
	}

	return result
}

// GoModInit runs `go mod init <module>`
func (e *Executor) GoModInit(ctx context.Context, modulePath string) error {
	result := e.Run(ctx, "go", "mod", "init", modulePath)
	if result.Error != nil {
		return fmt.Errorf("go mod init failed: %s", result.Stderr)
	}
	return nil
}

// GoGet runs `go get <packages...>`
func (e *Executor) GoGet(ctx context.Context, packages []string) error {
	for _, pkg := range packages {
		result := e.Run(ctx, "go", "get", pkg)
		if result.Error != nil {
			return fmt.Errorf("go get %s failed: %s", pkg, result.Stderr)
		}
	}
	return nil
}

// GoModTidy runs `go mod tidy`
func (e *Executor) GoModTidy(ctx context.Context) error {
	result := e.Run(ctx, "go", "mod", "tidy")
	if result.Error != nil {
		return fmt.Errorf("go mod tidy failed: %s", result.Stderr)
	}
	return nil
}

// NpmInit runs `npm init -y`
func (e *Executor) NpmInit(ctx context.Context) error {
	result := e.Run(ctx, "npm", "init", "-y")
	if result.Error != nil {
		return fmt.Errorf("npm init failed: %s", result.Stderr)
	}
	return nil
}

// NpmInstall runs `npm install <packages...>`
func (e *Executor) NpmInstall(ctx context.Context, packages []string, dev bool) error {
	args := []string{"install"}
	args = append(args, packages...)
	if dev {
		args = append(args, "--save-dev")
	}

	result := e.Run(ctx, "npm", args...)
	if result.Error != nil {
		return fmt.Errorf("npm install failed: %s", result.Stderr)
	}
	return nil
}

// AirInit runs `air init`
func (e *Executor) AirInit(ctx context.Context) error {
	result := e.Run(ctx, "air", "init")
	if result.Error != nil {
		return fmt.Errorf("air init failed: %s", result.Stderr)
	}
	return nil
}

// TailwindInit runs `npx tailwindcss init`
func (e *Executor) TailwindInit(ctx context.Context) error {
	npx := "npx"
	if runtime.GOOS == "windows" {
		npx = "npx.cmd"
	}

	result := e.Run(ctx, npx, "tailwindcss", "init")
	if result.Error != nil {
		return fmt.Errorf("tailwindcss init failed: %s", result.Stderr)
	}
	return nil
}
