package pkg

import tea "github.com/charmbracelet/bubbletea"

func keyHandler(m *Model, msg tea.KeyMsg) tea.Cmd {

	key := msg.String()

	// Handle pane switching
	switch key {
	case "f2":
		m.activePane = ChatPane
		return nil

	case "f12":
		if m.activePane == LogsPane {
			m.activePane = m.lastPane
		} else {
			m.lastPane = m.activePane
			m.activePane = LogsPane
		}
		return nil
	}

	switch m.activePane {
	case LoginPane:
		return loginKeyHandler(m, key)
	}

	return nil
}
