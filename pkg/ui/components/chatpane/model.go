package chatpane

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxneuvians/copilot-french-tutor/pkg/ui/consts"
	proxy "github.com/maxneuvians/go-copilot-proxy/pkg"
)

type chatMsg struct {
	content   string
	role      string
	timestamp time.Time
}

type chatResponse string

type Model struct {
	loading      bool
	msgs         []chatMsg
	sessionToken string
	spinner      spinner.Model
	ta           textarea.Model
	vp           viewport.Model
}

func New() Model {
	m := Model{
		loading: false,
		spinner: spinner.New(),
		ta:      textarea.New(),
		vp:      viewport.New(0, 0),
	}

	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m.spinner.Spinner = spinner.Meter

	m.ta.Placeholder = "Ask your questions about french grammar here. Press enter to send."
	m.ta.Focus()
	m.ta.Cursor.Blink = false
	m.ta.Prompt = "| "
	m.ta.CharLimit = 280
	m.ta.SetHeight(2)
	m.ta.KeyMap.InsertNewline.SetEnabled(false)
	m.ta.ShowLineNumbers = false

	// Set system prompt
	m.msgs = append(m.msgs, chatMsg{
		content:   "You are a french tutor who is helping a student learn french. Please assist them in learning grammar.",
		role:      "system",
		timestamp: time.Now(),
	})

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case chatResponse:
		m.msgs = append(m.msgs, chatMsg{
			content:   string(msg),
			role:      "assistant",
			timestamp: time.Now(),
		})
		m.loading = false
		m.updateConversation()

	case consts.SessionTokenUpdateMsg:
		m.sessionToken = string(msg)

	case tea.KeyMsg:
		if msg.String() == "enter" {
			if !m.loading {
				m.msgs = append(m.msgs, chatMsg{
					content:   m.ta.Value(),
					role:      "user",
					timestamp: time.Now(),
				})
				m.ta.SetValue("")
				m.loading = true
				m.updateConversation()
				return m, tea.Batch(m.spinner.Tick, m.sendChatMessages())
			}
		}

	case tea.WindowSizeMsg:
		m.ta.SetWidth(msg.Width - 2)
		m.vp.Width = msg.Width - 2
		m.vp.Height = msg.Height - 6
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.vp, cmd = m.vp.Update(msg)
		cmds = append(cmds, cmd)

		m.ta, cmd = m.ta.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	if m.loading {
		return m.vp.View() + "\n" + m.spinner.View()
	}

	return m.vp.View() + "\n" + m.ta.View()
}

func renderChatLine(msg *chatMsg) string {
	userStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	assistantStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#A550DF"))

	switch msg.role {
	case "user":
		return userStyle.Render("["+msg.timestamp.Format("02 Jan 06 15:04:05 MST")+"] You: ") + msg.content + "\n\n"
	case "system":
		return ""
	case "assistant":
		return assistantStyle.Render("["+msg.timestamp.Format("02 Jan 06 15:04:05 MST")+"] Assistant: ") + msg.content + "\n\n"
	}
	return ""
}

func (m *Model) sendChatMessages() tea.Cmd {
	return func() tea.Msg {
		var chatMessages []proxy.Message

		for _, msg := range m.msgs {
			chatMessages = append(chatMessages, proxy.Message{Content: msg.content, Role: msg.role})
		}

		var resp string
		err := proxy.Chat(m.sessionToken, chatMessages, false, func(cr proxy.CompletionResponse) error {
			resp = cr.Choices[0].Message.Content
			return nil
		})

		if err != nil {
			return chatResponse("Error sending chat messages.")
		}

		return chatResponse(resp)
	}
}

func (m *Model) updateConversation() {
	var conversation string

	for _, msg := range m.msgs {
		conversation += renderChatLine(&msg)
	}

	m.vp.SetContent(conversation)
}
