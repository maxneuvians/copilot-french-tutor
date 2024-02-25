package pkg

import (
	"github.com/charmbracelet/lipgloss"
)

func renderLogsPane(m *Model) {
	logs := ""

	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	contentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

	for _, log := range m.logs {
		logs += timeStyle.Render(log.time.Format("02 Jan 06 15:04:05 MST")) + " " + contentStyle.Render(log.content) + "\n"
	}

	m.pane.Content = logs

}
