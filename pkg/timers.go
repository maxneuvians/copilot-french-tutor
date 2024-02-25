package pkg

import (
	"os"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maxneuvians/go-copilot-proxy/pkg"
)

func handleTimerTick(m *Model, msg timer.TickMsg) tea.Cmd {
	var cmd tea.Cmd

	if m.loginTimer.ID() == msg.ID {
		if int(m.loginTimer.Timeout.Seconds())%6 == 0 {

			authResponse, err := pkg.Authenticate(m.loginResponse)

			if err != nil {
				m.addLogMessage("Error authenticating.")
			}

			if authResponse.AccessToken != "" {
				m.addLogMessage("Token received.")
				file, err := os.Create(TokenFile)

				if err != nil {
					m.addLogMessage("Error creating token file.")
				}

				_, err = file.WriteString(authResponse.AccessToken)

				if err != nil {
					m.addLogMessage("Error writing token to file.")
				}

				m.accessToken = authResponse.AccessToken
				m.sessionState = LoggedIn
				getSessionToken(m)
				return m.loginTimer.Stop()
			}
		}
		m.loginTimer, cmd = m.loginTimer.Update(msg)
	}

	return cmd
}
