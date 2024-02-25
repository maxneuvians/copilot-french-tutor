package pkg

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	proxy "github.com/maxneuvians/go-copilot-proxy/pkg"
)

type logMessage struct {
	content string
	time    time.Time
}

type Model struct {
	accessToken   string
	activePane    string
	height        int
	lastPane      string
	loginResponse proxy.LoginResponse
	loginTimer    timer.Model
	logs          []logMessage
	pane          pane
	sessionState  string
	sessionToken  string
	width         int
}

func InitialModel() Model {
	m := Model{
		activePane: LoginPane,
		loginTimer: timer.NewWithInterval(2*time.Minute, time.Second),
	}

	m.addLogMessage("Initializing...")
	setSessionState(&m)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
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
		cmds = append(cmds, keyHandler(&m, msg))

	case timer.TickMsg:
		cmds = append(cmds, handleTimerTick(&m, msg))

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	renderPane(&m, msg)
	m.pane.viewport, cmd = m.pane.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s", m.pane.viewport.View(), m.renderStatusBar())
}

func (m *Model) addLogMessage(content string) {
	// Prepend the new log message
	m.logs = append([]logMessage{{content, time.Now()}}, m.logs...)

}
