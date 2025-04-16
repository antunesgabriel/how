package confirm

import (
	"github.com/antunesgabriel/how-ai/presetation/models"
	"github.com/antunesgabriel/how-ai/presetation/theme"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyMap defines the keybindings for the confirmation dialog
type KeyMap struct {
	Yes key.Binding
	No  key.Binding
}

// DefaultKeyMap returns the default keybindings
var DefaultKeyMap = KeyMap{
	Yes: key.NewBinding(
		key.WithKeys("y", "Y"),
		key.WithHelp("y", "yes"),
	),
	No: key.NewBinding(
		key.WithKeys("n", "N"),
		key.WithHelp("n", "no"),
	),
}

// Model represents the confirmation dialog
type Model struct {
	Command      models.Command
	Width        int
	Height       int
	KeyMap       KeyMap
	BorderColor  lipgloss.Color
	SelectedItem int // 0 for No, 1 for Yes
}

// New creates a new confirmation dialog
func New() Model {
	return Model{
		KeyMap:       DefaultKeyMap,
		BorderColor:  theme.Colors.Green,
		SelectedItem: 0, // Default to "No" for safety
	}
}

// SetCommand sets the command to be confirmed
func (m *Model) SetCommand(cmd models.Command) {
	m.Command = cmd
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and user input
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Yes):
			return m, func() tea.Msg {
				return models.CommandConfirmationMsg{
					Command:  m.Command,
					Approved: true,
				}
			}
		case key.Matches(msg, m.KeyMap.No):
			return m, func() tea.Msg {
				return models.CommandConfirmationMsg{
					Command:  m.Command,
					Approved: false,
				}
			}
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	return m, nil
}

// View renders the confirmation dialog
func (m Model) View() string {
	if m.Command.Raw == "" {
		return ""
	}

	// Adjust dialog width based on terminal size
	dialogWidth := min(m.Width-4, 80)

	// Title with warning icon
	title := lipgloss.NewStyle().
		Foreground(theme.Colors.Yellow).
		Bold(true).
		Render("⚠️  Execute Command?")

	// Format command with box and syntax highlighting
	commandStyle := lipgloss.NewStyle().
		Foreground(theme.Colors.Cyan).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.Colors.Gray).
		Padding(0, 1).
		Width(dialogWidth - 8)

	commandBox := commandStyle.Render(m.Command.Raw)

	// Description/context if available
	description := ""
	if m.Command.Description != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(theme.Colors.White).
			Width(dialogWidth - 8)

		description = "\n\n" + descStyle.Render("Context: "+m.Command.Description)
	}

	// Action buttons with highlighting for selected option
	noBtn := "[ No ]"
	yesBtn := "[ Yes ]"

	if m.SelectedItem == 0 {
		// No is selected (default for safety)
		noBtn = lipgloss.NewStyle().
			Background(theme.Colors.Green).
			Foreground(theme.Colors.Black).
			Bold(true).
			Padding(0, 1).
			Render("[ No ]")
	} else {
		yesBtn = lipgloss.NewStyle().
			Background(theme.Colors.Green).
			Foreground(theme.Colors.Black).
			Bold(true).
			Padding(0, 1).
			Render("[ Yes ]")
	}

	// Space the buttons properly
	actions := lipgloss.JoinHorizontal(
		lipgloss.Center,
		noBtn,
		"     ",
		yesBtn,
	)

	// Help text at the bottom
	help := lipgloss.NewStyle().
		Foreground(theme.Colors.Gray).
		Italic(true).
		Render("Press 'y' to execute or 'n' to cancel")

	// Warning about command execution
	warning := lipgloss.NewStyle().
		Foreground(theme.Colors.Yellow).
		Width(dialogWidth - 8).
		Align(lipgloss.Center).
		Render("This command will be executed on your system")

	// Combine all parts with proper spacing
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		warning,
		"",
		commandBox,
		description,
		"",
		actions,
		"",
		help,
	)

	// Dialog box with border
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.BorderColor).
		Padding(1, 2).
		Width(dialogWidth).
		Align(lipgloss.Center)

	dialog := dialogStyle.Render(content)

	// Calculate dialog height for proper placement
	dialogHeight := 10
	if m.Command.Description != "" {
		dialogHeight += 2
	}

	// Center in terminal window
	return lipgloss.Place(
		m.Width,
		dialogHeight,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

// min returns the smaller of a or b
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
