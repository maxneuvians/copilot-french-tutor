package pkg

import "github.com/charmbracelet/lipgloss"

func (m Model) renderStatusBar() string {
	buttonNugget := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Padding(0, 1)

	buttonBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	buttonStyle := lipgloss.NewStyle().
		Inherit(buttonBarStyle).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#FF5F87")).
		Padding(0, 1).
		MarginRight(1)

	exitButtonStyle := buttonNugget.Copy().
		Background(lipgloss.Color("#A550DF")).
		Align(lipgloss.Right)

	buttonText := lipgloss.NewStyle().Inherit(buttonBarStyle)

	sessionStateStyle := buttonNugget.Copy().Background(lipgloss.Color("#6124DF"))

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
