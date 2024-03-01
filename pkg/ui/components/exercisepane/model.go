package exercisepane

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maxneuvians/copilot-french-tutor/pkg/ui/consts"
	proxy "github.com/maxneuvians/go-copilot-proxy/pkg"
)

type exerciseResponse string

type Model struct {
	exercises      []map[string]string
	exerciseIdx    int
	height         int
	loading        bool
	feedback       string
	selectedPrompt prompt
	sessionToken   string
	spinner        spinner.Model
	table          table.Model
	ti             textinput.Model
	width          int
}

func New() Model {
	m := Model{}

	rows := []table.Row{}

	for i, p := range prompts {
		rows = append(rows, table.Row{fmt.Sprintf("%d", i+1), p.name})
	}

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Exercise name", Width: 20},
	}

	m.table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	m.table.SetStyles(s)

	m.ti = textinput.New()
	m.ti.CharLimit = 280

	m.spinner = spinner.New()
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m.spinner.Spinner = spinner.Points

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

	case consts.SessionTokenUpdateMsg:
		m.sessionToken = string(msg)

	case exerciseResponse:
		m.loading = false
		m.feedback = ""
		// Parse json string into exercises
		err := json.Unmarshal([]byte(msg), &m.exercises)
		if err != nil {
			return m, nil
		}
		m.exerciseIdx = 0

	case tea.KeyMsg:
		if msg.String() == "enter" {
			if len(m.exercises) > 0 {
				if m.feedback != "" {
					m.feedback = ""
					m.exerciseIdx++
					if m.exerciseIdx >= len(m.exercises) {
						m.exercises = nil
					}
					m.ti.SetValue("")
				} else {
					want := m.exercises[m.exerciseIdx]["answer"]
					got := m.ti.Value()
					if want == got {
						m.feedback = "Correct!"
					} else {
						m.feedback = fmt.Sprintf("Incorrect. The correct answer is: %s, you wrote: %s", want, got)
					}
				}
				return m, nil
			}
			m.loading = true
			selectionID, _ := strconv.Atoi(m.table.SelectedRow()[0])
			selectionID--
			m.selectedPrompt = prompts[selectionID]
			m.ti.Focus()
			cmds = append(cmds, m.sendExerciseRequest(), m.spinner.Tick)
		}

	case tea.WindowSizeMsg:
		columns := []table.Column{
			{Title: "ID", Width: 4},
			{Title: "Exercise name", Width: msg.Width - 10},
		}
		m.table.SetColumns(columns)
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(msg.Height - 6)

		m.width = msg.Width
		m.height = msg.Height
	}

	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)

		m.ti, cmd = m.ti.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.loading {
		return m.spinner.View()
	}

	if len(m.exercises) > 0 {
		return m.renderExercise(m.exerciseIdx)
	}
	return m.table.View()
}

func (m *Model) renderExercise(index int) string {

	q := fmt.Sprintf(`
		%s

		Question: %s
		
		Your answer: %s

		Translation: %s`,
		renderFeedbackStyle(m.feedback),
		m.exercises[index]["question"],
		m.ti.View(),
		m.exercises[index]["translation"])

	return lipgloss.Place(m.width, m.height/2,
		lipgloss.Left, lipgloss.Center,
		lipgloss.NewStyle().Render(q))

}

func renderFeedbackStyle(feedback string) string {
	if strings.HasPrefix(feedback, "Correct") {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render(feedback)
	} else {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(feedback)
	}
}

func (m *Model) sendExerciseRequest() tea.Cmd {
	return func() tea.Msg {
		var chatMessages []proxy.Message

		chatMessages = append(chatMessages, proxy.Message{
			Role: "system",
			Content: `You are a system that generates simple quiz questions for french grammar in the context of working in the Canadian government. 
				For example you could be asked, 'Generate five questions that test the students understanding of the verb aller in the past tense'
				You provide your question in a json format. 
				If you are returning multiple questions, please use a JSON array. 
				Your response should be valid JSON that is minized. 
				Please validate that the answers form gramatically correct sentences when inserted into the sentences.
				`,
		})

		chatMessages = append(chatMessages, proxy.Message{
			Role:    "user",
			Content: m.selectedPrompt.prompt,
		})

		var resp string
		err := proxy.Chat(m.sessionToken, chatMessages, false, func(cr proxy.CompletionResponse) error {
			resp = cr.Choices[0].Message.Content
			return nil
		})

		if err != nil {
			return exerciseResponse("Error sending chat messages.")
		}

		return exerciseResponse(resp)
	}
}
