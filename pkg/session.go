package pkg

import (
	"bufio"
	"os"

	"github.com/maxneuvians/go-copilot-proxy/pkg"
)

func setSessionState(m *Model) {
	//Check if .github_copilot_token file exists
	_, err := os.Stat(TokenFile)

	if err != nil {
		m.sessionState = LoggedOut
		m.addLogMessage("No token file found.")
		return
	}

	// Get the authentication token
	file, err := os.Open(TokenFile)

	if err != nil {
		m.addLogMessage("Error reading token file.")
		m.sessionState = LoggedOut
		return
	}

	// If the file exists, read the first line
	r := bufio.NewReader(file)
	buffer, _, err := r.ReadLine()

	if err != nil {
		m.addLogMessage("Error reading first line of token file.")
		m.sessionState = LoggedOut
		return
	}

	token := string(buffer)

	// Check if the token is valid by starting with "ghu_"
	if len(token) < 4 || token[:4] != "ghu_" {
		m.addLogMessage("Invalid token found in token file.")
		m.sessionState = LoggedOut
		return
	}

	m.accessToken = token
	m.sessionState = LoggedIn
	getSessionToken(m)
	m.addLogMessage("Token found in token file.")
}

func getSessionToken(m *Model) {
	// Check if an access token is set
	if m.accessToken == "" {
		m.addLogMessage("No access token found.")
		return
	}

	sessionResponse, err := pkg.GetSessionToken(m.accessToken)

	if err != nil {
		m.addLogMessage("Error getting session token.")
		return
	}

	m.addLogMessage("Session token received.")
	m.sessionToken = sessionResponse.Token
}
