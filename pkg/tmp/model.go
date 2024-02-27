package pkg

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/bubbles/viewport"
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
	chatMessages  []proxy.Message
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
		pane: pane{
			viewport: viewport.Model{},
		},
	}

	// Disable the mouse
	m.pane.viewport.MouseWheelEnabled = false
	m.pane.viewport.HighPerformanceRendering = false

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

	//m.addLogMessage(fmt.Sprintf("Received message: %T", msg))

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

		// Set the dimensions
		m.pane.viewport.Width = msg.Width
		m.pane.viewport.Height = msg.Height

		m.pane.viewport, cmd = m.pane.viewport.Update(msg)
		cmds = append(cmds, cmd)

	case chatMsg:
		m.chatMessages = append(m.chatMessages, proxy.Message{Content: string(msg), Role: "assistant"})
		m.pane.loading = false
		m.pane.textarea.Reset()
	}

	cmd = renderPane(&m, msg)
	cmds = append(cmds, cmd)

	if m.activePane == ChatPane {
		m.pane.textarea, cmd = m.pane.textarea.Update(msg)
		cmds = append(cmds, cmd)

		if m.pane.loading {
			m.pane.spinner, cmd = m.pane.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

		m.pane.display, cmd = m.pane.display.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s", m.renderStatusBar(), m.pane.viewport.View())
}

func (m *Model) addLogMessage(content string) {
	// Prepend the new log message
	m.logs = append([]logMessage{{content, time.Now()}}, m.logs...)

}
