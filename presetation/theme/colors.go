package theme

import "github.com/charmbracelet/lipgloss"

// Colors used throughout the application
var Colors = struct {
	Green  lipgloss.Color
	Blue   lipgloss.Color
	Red    lipgloss.Color
	Cyan   lipgloss.Color
	Gray   lipgloss.Color
	White  lipgloss.Color
	Black  lipgloss.Color
	Yellow lipgloss.Color
}{
	Green:  lipgloss.Color("#8bd5ca"),
	Blue:   lipgloss.Color("#89b4fa"),
	Red:    lipgloss.Color("#f38ba8"),
	Cyan:   lipgloss.Color("#74c7ec"),
	Gray:   lipgloss.Color("#6c7086"),
	White:  lipgloss.Color("#cdd6f4"),
	Black:  lipgloss.Color("#11111b"),
	Yellow: lipgloss.Color("#f9e2af"),
}

// Catppuccin Mocha color scheme
var (
	Rosewater = lipgloss.Color("#f5e0dc")
	Flamingo  = lipgloss.Color("#f2cdcd")
	Pink      = lipgloss.Color("#f5c2e7")
	Mauve     = lipgloss.Color("#cba6f7")
	Red       = lipgloss.Color("#f38ba8")
	Maroon    = lipgloss.Color("#eba0ac")
	Peach     = lipgloss.Color("#fab387")
	Yellow    = lipgloss.Color("#f9e2af")
	Green     = lipgloss.Color("#a6e3a1")
	Teal      = lipgloss.Color("#94e2d5")
	Sky       = lipgloss.Color("#89dceb")
	Sapphire  = lipgloss.Color("#74c7ec")
	Blue      = lipgloss.Color("#89b4fa")
	Lavender  = lipgloss.Color("#b4befe")
	Text      = lipgloss.Color("#cdd6f4")
	Subtext1  = lipgloss.Color("#bac2de")
	Subtext0  = lipgloss.Color("#a6adc8")
	Overlay2  = lipgloss.Color("#9399b2")
	Overlay1  = lipgloss.Color("#7f849c")
	Overlay0  = lipgloss.Color("#6c7086")
	Surface2  = lipgloss.Color("#585b70")
	Surface1  = lipgloss.Color("#45475a")
	Surface0  = lipgloss.Color("#313244")
	Base      = lipgloss.Color("#1e1e2e")
	Mantle    = lipgloss.Color("#181825")
	Crust     = lipgloss.Color("#11111b")
)
