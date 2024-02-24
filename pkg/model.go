package pkg

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	vp viewport.Model
}

func initialModel() model {

	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:

		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:

		m.vp = renderVP(msg.Width, msg.Height)
	}

	m.vp, _ = m.vp.Update(msg)

	return m, nil
}

func (m model) View() string {
	return m.vp.View()
}
