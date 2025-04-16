package input

import "github.com/charmbracelet/bubbles/key"

var Keymap = struct {
	Send key.Binding
	Quit key.Binding
}{
	Send: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "send message"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc/ctrl+c", "quit"),
	),
}
