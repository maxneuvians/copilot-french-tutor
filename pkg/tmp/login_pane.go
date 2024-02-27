package pkg

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	proxy "github.com/maxneuvians/go-copilot-proxy/pkg"
)

func loginKeyHandler(m *Model, key string) tea.Cmd {
	switch key {
	case "enter":
		if m.sessionState == LoggedOut {
			m.sessionState = LoggingIn
			return m.loginTimer.Init()
		}
	}
	return nil
}

func getLoginText() string {
	return lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("You are not logged in. Would you like to log in? \n\n Press Enter to log in or Ctrl+C to quit.")
}

func getDeviceCodeText(code string, seconds string) string {
	return lipgloss.
		NewStyle().
		Width(50).
		Align(lipgloss.Center).
		Render("Please go to https://github.com/login/device and enter the following code: \n\n" + code + "\n\n Checking in " + seconds + " seconds.")
}

func renderLoginPane(m *Model) {
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 1).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	var content string

	switch m.sessionState {

	case LoggedOut:
		content = getLoginText()

	case LoggingIn:
		if m.loginResponse.UserCode == "" {
			loginResponse, err := proxy.Login()
			if err != nil {
				m.addLogMessage("Error getting device code.")
				content = getLoginText()
			} else {
				m.loginResponse = loginResponse
				m.addLogMessage("Device code: " + m.loginResponse.UserCode)
				content = getDeviceCodeText(m.loginResponse.UserCode, "5")

			}
		} else {
			remaingTime := int(m.loginTimer.Timeout.Seconds()) % 6
			content = getDeviceCodeText(m.loginResponse.UserCode, strconv.Itoa(remaingTime))
		}

	case LoggedIn:
		content = lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("You are logged in. \n\n Choose one of the options below using the F2-F4 keys.")
	}
	dialog := lipgloss.Place(m.width, m.height/2,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(content))

	m.pane.Content = dialog
}
