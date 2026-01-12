package main

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keybindings for gritt
type KeyMap struct {
	// Actions
	Execute     key.Binding
	ToggleDebug key.Binding
	CyclePane   key.Binding
	ClosePane   key.Binding
	Quit        key.Binding
	Help        key.Binding

	// Navigation
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Home  key.Binding
	End   key.Binding
	PgUp  key.Binding
	PgDn  key.Binding

	// Editing
	Backspace key.Binding
	Delete    key.Binding
}

// DefaultKeyMap provides the default keybindings
var DefaultKeyMap = KeyMap{
	Execute: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "execute"),
	),
	ToggleDebug: key.NewBinding(
		key.WithKeys("f12"),
		key.WithHelp("F12", "debug"),
	),
	CyclePane: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "cycle pane"),
	),
	ClosePane: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close pane"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("C-c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),

	// Navigation
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "right"),
	),
	Home: key.NewBinding(
		key.WithKeys("home"),
		key.WithHelp("home", "line start"),
	),
	End: key.NewBinding(
		key.WithKeys("end"),
		key.WithHelp("end", "line end"),
	),
	PgUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "page up"),
	),
	PgDn: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdn", "page down"),
	),

	// Editing
	Backspace: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("bksp", "delete back"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete"),
		key.WithHelp("del", "delete forward"),
	),
}

// ShortHelp returns keybindings for the short help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Execute, k.ToggleDebug, k.Help, k.Quit}
}

// FullHelp returns keybindings for the full help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Execute, k.ToggleDebug, k.CyclePane, k.ClosePane},
		{k.Up, k.Down, k.Left, k.Right},
		{k.Home, k.End, k.PgUp, k.PgDn},
		{k.Help, k.Quit},
	}
}
