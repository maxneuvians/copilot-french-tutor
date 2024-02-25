package pkg

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pane struct {
	Content  string
	viewport viewport.Model
}

func renderPane(m *Model, msg tea.Msg) {
	// Create the viewport
	m.pane.viewport = viewport.Model{}

	// Set the dimensions
	m.pane.viewport.Width = m.width
	m.pane.viewport.Height = m.height

	// Set the content
	switch m.activePane {

	case ChatPane:
		renderChatPane(m)

	case LoginPane:
		renderLoginPane(m)

	case LogsPane:
		renderLogsPane(m)

	}

	m.pane.viewport.SetContent(m.pane.Content)

	// Apply styles
	m.pane.viewport.Style = lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.height - 4).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

}
