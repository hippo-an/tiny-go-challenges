package prompt

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hippo-an/tiny-go-challenges/protem-gen/internal/config"
)

// Step represents the current step in the wizard
type Step int

const (
	StepName Step = iota
	StepModule
	StepDatabase
	StepFeatures
	StepConfirm
	StepDone
)

// Model is the Bubbletea model for the interactive prompt
type Model struct {
	step      Step
	Config    *config.ProjectConfig
	Cancelled bool

	// Text inputs
	nameInput   textinput.Model
	moduleInput textinput.Model

	// Selection states
	databaseIdx int
	featureIdx  int
	features    []bool // [grpc, auth, ai]

	// UI state
	width  int
	height int
	err    error
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	checkboxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))
)

// NewModel creates a new prompt model
func NewModel() Model {
	nameInput := textinput.New()
	nameInput.Placeholder = "my-app"
	nameInput.Focus()
	nameInput.CharLimit = 64
	nameInput.Width = 40

	moduleInput := textinput.New()
	moduleInput.Placeholder = "github.com/username/my-app"
	moduleInput.CharLimit = 128
	moduleInput.Width = 50

	return Model{
		step:        StepName,
		Config:      config.NewDefaultConfig(),
		nameInput:   nameInput,
		moduleInput: moduleInput,
		features:    []bool{false, false, false}, // grpc, auth, ai
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.Cancelled = true
			return m, tea.Quit

		case "enter":
			return m.handleEnter()

		case "up", "k":
			return m.handleUp()

		case "down", "j":
			return m.handleDown()

		case "tab":
			return m.handleTab()

		case " ":
			return m.handleSpace()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Handle text input updates
	var cmd tea.Cmd
	switch m.step {
	case StepName:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case StepModule:
		m.moduleInput, cmd = m.moduleInput.Update(msg)
	}

	return m, cmd
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case StepName:
		if m.nameInput.Value() == "" {
			m.err = fmt.Errorf("project name is required")
			return m, nil
		}
		m.Config.Name = m.nameInput.Value()
		m.err = nil

		// Pre-fill module path suggestion
		if m.moduleInput.Value() == "" {
			m.moduleInput.SetValue(fmt.Sprintf("github.com/username/%s", m.Config.Name))
		}

		m.step = StepModule
		m.nameInput.Blur()
		m.moduleInput.Focus()
		return m, textinput.Blink

	case StepModule:
		if m.moduleInput.Value() == "" {
			m.err = fmt.Errorf("module path is required")
			return m, nil
		}
		m.Config.ModulePath = m.moduleInput.Value()
		m.err = nil
		m.step = StepDatabase
		m.moduleInput.Blur()
		return m, nil

	case StepDatabase:
		databases := config.DatabaseOptions()
		m.Config.Database = databases[m.databaseIdx]
		m.step = StepFeatures
		return m, nil

	case StepFeatures:
		m.Config.IncludeGRPC = m.features[0]
		m.Config.IncludeAuth = m.features[1]
		m.Config.IncludeAI = m.features[2]
		m.step = StepConfirm
		return m, nil

	case StepConfirm:
		m.step = StepDone
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) handleUp() (tea.Model, tea.Cmd) {
	switch m.step {
	case StepDatabase:
		if m.databaseIdx > 0 {
			m.databaseIdx--
		}
	case StepFeatures:
		if m.featureIdx > 0 {
			m.featureIdx--
		}
	}
	return m, nil
}

func (m Model) handleDown() (tea.Model, tea.Cmd) {
	switch m.step {
	case StepDatabase:
		if m.databaseIdx < len(config.DatabaseOptions())-1 {
			m.databaseIdx++
		}
	case StepFeatures:
		if m.featureIdx < 2 {
			m.featureIdx++
		}
	}
	return m, nil
}

func (m Model) handleTab() (tea.Model, tea.Cmd) {
	return m.handleDown()
}

func (m Model) handleSpace() (tea.Model, tea.Cmd) {
	if m.step == StepFeatures {
		m.features[m.featureIdx] = !m.features[m.featureIdx]
	}
	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("protem-gen - Go Web Application Generator"))
	b.WriteString("\n\n")

	switch m.step {
	case StepName:
		b.WriteString(m.viewNameStep())
	case StepModule:
		b.WriteString(m.viewModuleStep())
	case StepDatabase:
		b.WriteString(m.viewDatabaseStep())
	case StepFeatures:
		b.WriteString(m.viewFeaturesStep())
	case StepConfirm:
		b.WriteString(m.viewConfirmStep())
	}

	// Help
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press Enter to continue • Esc to cancel"))

	return b.String()
}

func (m Model) viewNameStep() string {
	var b strings.Builder
	b.WriteString("Project name:\n")
	b.WriteString(m.nameInput.View())
	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.err.Error()))
	}
	return b.String()
}

func (m Model) viewModuleStep() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Project: %s\n\n", selectedStyle.Render(m.Config.Name)))
	b.WriteString("Go module path:\n")
	b.WriteString(m.moduleInput.View())
	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.err.Error()))
	}
	return b.String()
}

func (m Model) viewDatabaseStep() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Project: %s\n", selectedStyle.Render(m.Config.Name)))
	b.WriteString(fmt.Sprintf("Module:  %s\n\n", selectedStyle.Render(m.Config.ModulePath)))
	b.WriteString("Select database:\n\n")

	databases := []struct {
		name string
		desc string
	}{
		{"postgres", "PostgreSQL with pgx driver"},
		{"mysql", "MySQL with go-sql-driver"},
		{"sqlite", "SQLite (file-based, no server)"},
		{"none", "No database integration"},
	}

	for i, db := range databases {
		cursor := "  "
		style := blurredStyle
		if i == m.databaseIdx {
			cursor = "> "
			style = selectedStyle
		}
		b.WriteString(fmt.Sprintf("%s%s - %s\n", cursor, style.Render(db.name), db.desc))
	}

	return b.String()
}

func (m Model) viewFeaturesStep() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Project:  %s\n", selectedStyle.Render(m.Config.Name)))
	b.WriteString(fmt.Sprintf("Module:   %s\n", selectedStyle.Render(m.Config.ModulePath)))
	b.WriteString(fmt.Sprintf("Database: %s\n\n", selectedStyle.Render(string(m.Config.Database))))
	b.WriteString("Select optional features (Space to toggle):\n\n")

	features := []struct {
		name string
		desc string
	}{
		{"gRPC", "Include gRPC server support"},
		{"Auth", "Include authentication boilerplate"},
		{"AI", "Include AI/LLM integration boilerplate"},
	}

	for i, feat := range features {
		cursor := "  "
		style := blurredStyle
		if i == m.featureIdx {
			cursor = "> "
			style = focusedStyle
		}

		checkbox := "[ ]"
		if m.features[i] {
			checkbox = checkboxStyle.Render("[x]")
		}

		b.WriteString(fmt.Sprintf("%s%s %s - %s\n", cursor, checkbox, style.Render(feat.name), feat.desc))
	}

	return b.String()
}

func (m Model) viewConfirmStep() string {
	var b strings.Builder
	b.WriteString("Configuration Summary:\n\n")
	b.WriteString(fmt.Sprintf("  Project:  %s\n", selectedStyle.Render(m.Config.Name)))
	b.WriteString(fmt.Sprintf("  Module:   %s\n", selectedStyle.Render(m.Config.ModulePath)))
	b.WriteString(fmt.Sprintf("  Database: %s\n", selectedStyle.Render(string(m.Config.Database))))

	var features []string
	if m.features[0] {
		features = append(features, "gRPC")
	}
	if m.features[1] {
		features = append(features, "Auth")
	}
	if m.features[2] {
		features = append(features, "AI")
	}
	if len(features) == 0 {
		features = append(features, "none")
	}
	b.WriteString(fmt.Sprintf("  Features:  %s\n", selectedStyle.Render(strings.Join(features, ", "))))

	b.WriteString("\nPress Enter to create project")

	return b.String()
}
