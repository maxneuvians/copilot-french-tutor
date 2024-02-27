package loginpane

import (
	"bufio"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxneuvians/copilot-french-tutor/pkg/ui/consts"
	proxy "github.com/maxneuvians/go-copilot-proxy/pkg"
)

var (
	dialogBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 1).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
)

type Model struct {
	accessToken  string
	deviceCode   string
	height       int
	sessionState string
	userCode     string
	width        int
}

func New() Model {
	m := Model{
		sessionState: consts.LoggedOut,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case consts.InitializingMsg:
		return m, m.updateSessionState()

	case consts.SessionUpdateMsg:
		m.accessToken = msg.AccessToken
		m.sessionState = msg.State

		if m.sessionState == consts.LoggedIn {
			return m, m.getSessionToken()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Model) View() string {

	content := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("You are logged in. \n\n Choose one of the options below using the F2-F4 keys.")

	return lipgloss.Place(m.width, m.height/2,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(content))
}

func (m Model) getSessionToken() tea.Cmd {
	// Check if an access token is set
	if m.accessToken == "" {
		return nil
	}

	return func() tea.Msg {
		sessionResponse, err := proxy.GetSessionToken(m.accessToken)

		if err != nil {
			return consts.SessionUpdateMsg{
				State: consts.LoggedOut,
			}
		}

		return consts.SessionTokenUpdateMsg(sessionResponse.Token)

	}
}

func (m *Model) updateSessionState() tea.Cmd {
	return func() tea.Msg {

		//Check if .github_copilot_token file exists
		_, err := os.Stat(consts.TokenFile)

		if err != nil {
			return consts.SessionUpdateMsg{
				State: consts.LoggedOut,
			}
		}

		// Get the authentication token
		file, err := os.Open(consts.TokenFile)

		if err != nil {
			return consts.SessionUpdateMsg{
				State: consts.LoggedOut,
			}
		}

		defer file.Close()

		r := bufio.NewReader(file)
		buffer, _, err := r.ReadLine()

		if err != nil {
			return consts.SessionUpdateMsg{
				State: consts.LoggedOut,
			}
		}

		token := string(buffer)

		// Check if the token is valid by starting with "ghu_"
		if len(token) < 4 || token[:4] != "ghu_" {
			return consts.SessionUpdateMsg{
				State: consts.LoggedOut,
			}
		}

		return consts.SessionUpdateMsg{
			AccessToken: token,
			State:       consts.LoggedIn,
		}
	}
}
