package loginpane

import (
	"bufio"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/timer"
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
	loginTimer   timer.Model
	sessionState string
	userCode     string
	width        int
}

func (m Model) GetSessionState() string {
	return m.sessionState
}

func New() Model {
	m := Model{
		loginTimer:   timer.NewWithInterval(2*time.Minute, time.Second),
		sessionState: consts.LoggedOut,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case consts.InitializingMsg:
		return m, m.updateSessionState()

	case consts.SessionUpdateMsg:
		m.accessToken = msg.AccessToken
		m.sessionState = msg.State

		if m.sessionState == consts.LoggedIn {
			return m, m.getSessionToken()
		}

	case tea.KeyMsg:
		if msg.String() == "enter" {
			if m.sessionState == consts.LoggedOut {
				loginResponse, err := proxy.Login()

				if err != nil {
					return m, nil
				}

				m.deviceCode = loginResponse.DeviceCode
				m.userCode = loginResponse.UserCode

				return m, tea.Batch(
					m.loginTimer.Init(),
					func() tea.Msg {
						return consts.SessionUpdateMsg{
							State: consts.LoggingIn,
						}
					})
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case timer.TickMsg:
		m.loginTimer, cmd = m.loginTimer.Update(msg)

		if int(m.loginTimer.Timeout.Seconds())%6 == 0 {
			resp, err := m.checkForAuthentication()

			if err != nil {
				return m, nil
			}

			if resp.AccessToken != "" {
				file, err := os.Create(consts.TokenFile)

				if err != nil {
					return m, nil
				}

				_, err = file.WriteString(resp.AccessToken)

				if err != nil {
					return m, nil
				}

				m.accessToken = resp.AccessToken
				return m, tea.Batch(m.getSessionToken(), m.loginTimer.Stop(), m.updateSessionState())
			}
		}
	}

	return m, cmd
}

func (m Model) View() string {
	var content string

	switch m.sessionState {

	case consts.LoggedOut:
		content = lipgloss.
			NewStyle().
			Width(50).
			Align(lipgloss.Center).
			Render("You are not logged in. Would you like to log in? \n\n Press Enter to log in or Ctrl+C to quit.")

	case consts.LoggingIn:
		remaingTime := int(m.loginTimer.Timeout.Seconds()) % 6
		content = lipgloss.
			NewStyle().
			Width(50).
			Align(lipgloss.Center).
			Render("Please go to https://github.com/login/device and enter the following code: \n\n" + m.userCode + "\n\n Checking in " + strconv.Itoa(remaingTime) + " seconds.")

	case consts.LoggedIn:
		content = lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("You are logged in. \n\n Choose one of the options below using the F2-F4 keys.")
	}

	return lipgloss.Place(m.width, m.height/2,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(content))
}

func (m Model) checkForAuthentication() (proxy.AuthenticationResponse, error) {
	payload := proxy.LoginResponse{
		DeviceCode: m.deviceCode,
		UserCode:   m.userCode,
	}

	return proxy.Authenticate(payload)
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
