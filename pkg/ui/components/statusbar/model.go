package statusbar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxneuvians/copilot-french-tutor/pkg/ui/consts"
)

var (
	buttonNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	buttonBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	buttonStyle = lipgloss.NewStyle().
			Inherit(buttonBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	buttonText = lipgloss.NewStyle().Inherit(buttonBarStyle)

	exitButtonStyle = buttonNugget.Copy().
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	sessionStateStyle = buttonNugget.Copy().Background(lipgloss.Color("#6124DF"))
)

type Model struct {
	sessionState string
	width        int
}

func New() Model {
	return Model{
		sessionState: consts.LoggedOut,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetSessionState(sessionState string) {
	m.sessionState = sessionState
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case consts.SessionUpdateMsg:
		m.SetSessionState(msg.State)

	}

	return m, nil
}

func (m Model) View() string {

	w := lipgloss.Width

	chatButton := buttonStyle.Render("F2: Chat")

	exitButton := exitButtonStyle.Render("F9: Quit")
	sessionState := sessionStateStyle.Render("🍥 " + m.sessionState)

	filler := buttonText.Copy().
		Width(m.width - w(chatButton) - w(exitButton) - w(sessionState)).
		Render("")

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		chatButton,
		filler,
		exitButton,
		sessionState,
	)

	return buttonBarStyle.Width(m.width).Render(bar)
}
