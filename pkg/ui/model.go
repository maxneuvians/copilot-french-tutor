package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxneuvians/copilot-french-tutor/pkg/ui/components/chatpane"
	"github.com/maxneuvians/copilot-french-tutor/pkg/ui/components/loginpane"
	"github.com/maxneuvians/copilot-french-tutor/pkg/ui/components/statusbar"
	"github.com/maxneuvians/copilot-french-tutor/pkg/ui/consts"
)

type Model struct {
	accessToken string
	activePane  string
	panes       map[string]tea.Model
	statusBar   statusbar.Model
	view        lipgloss.Style
}

func New() Model {
	m := Model{
		panes:     make(map[string]tea.Model),
		statusBar: statusbar.New(),
	}

	m.panes[consts.ChatPane] = chatpane.New()
	m.panes[consts.LoginPane] = loginpane.New()

	m.activePane = consts.LoginPane

	m.view = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	return m
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return consts.InitializingMsg(true)
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "f9" {
			return m, tea.Quit
		}

		if msg.String() == "f2" && m.panes[consts.LoginPane].(loginpane.Model).GetSessionState() == consts.LoggedIn {
			m.activePane = consts.ChatPane
		}

	case tea.WindowSizeMsg:
		m.view.Width(msg.Width - 2)
		m.view.Height(msg.Height - 4)
	}

	//m.panes[m.activePane], cmd = m.panes[m.activePane].Update(msg)
	//cmds = append(cmds, cmd)

	for k, v := range m.panes {
		// Do not send key messages to inactive panes
		if _, ok := msg.(tea.KeyMsg); ok {
			if k != m.activePane {
				continue
			}
		}

		m.panes[k], cmd = v.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	content := m.panes[m.activePane].View()
	return m.view.Render(content) + "\n" + m.statusBar.View()
}
