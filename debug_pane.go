package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// DebugPane displays the debug log using a viewport
type DebugPane struct {
	viewport viewport.Model
	log      *[]string // Pointer to Model's debugLog
	lastLen  int       // Track log length for auto-scroll
}

// NewDebugPane creates a debug pane backed by the given log slice
func NewDebugPane(log *[]string) *DebugPane {
	vp := viewport.New(0, 0)
	vp.MouseWheelEnabled = true
	return &DebugPane{
		viewport: vp,
		log:      log,
		lastLen:  0,
	}
}

func (d *DebugPane) Title() string {
	return "debug"
}

func (d *DebugPane) Render(w, h int) string {
	// Update viewport dimensions
	d.viewport.Width = w
	d.viewport.Height = h

	// Update content
	content := strings.Join(*d.log, "\n")
	d.viewport.SetContent(content)

	// Auto-scroll to bottom if new content was added
	if len(*d.log) > d.lastLen {
		d.viewport.GotoBottom()
		d.lastLen = len(*d.log)
	}

	return d.viewport.View()
}

func (d *DebugPane) HandleKey(msg tea.KeyMsg) bool {
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return cmd != nil
}

func (d *DebugPane) HandleMouse(x, y int, msg tea.MouseMsg) bool {
	var cmd tea.Cmd
	d.viewport, cmd = d.viewport.Update(msg)
	return cmd != nil
}
