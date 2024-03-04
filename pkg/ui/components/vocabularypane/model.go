package vocabularypane

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gocarina/gocsv"
)

type word struct {
	English string `csv:"english"`
	French  string `csv:"french"`
	German  string `csv:"german"`
}

type Model struct {
	drillList []word
	feedback  string
	height    int
	redoList  []word
	ti        textinput.Model
	width     int
	wordIdx   int
	words     []word
}

func New() Model {
	m := Model{
		wordIdx: 0,
		words:   make([]word, 0),
	}

	m.ti = textinput.New()
	m.ti.Placeholder = "Answer"
	m.ti.CharLimit = 100
	m.ti.Focus()

	// Get all the CSV file names in the ./vocabulay directory
	files, err := os.ReadDir("./vocabulary")

	if err != nil {
		fmt.Println(err)
	}

	// Read the CSV files and store the words in the words map
	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("./vocabulary/%v", file.Name()))

		if err != nil {
			fmt.Println(err)
		}

		defer f.Close()

		var words []word
		if err := gocsv.UnmarshalFile(f, &words); err != nil {
			fmt.Println(err)
		}

		m.words = append(m.words, words...)
	}

	// Randomize the order of the words
	rand.Shuffle(len(m.words), func(i, j int) {
		m.words[i], m.words[j] = m.words[j], m.words[i]
	})

	m.drillList = m.words[:10]

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
		if msg.String() == "enter" {
			if m.feedback != "" {
				m.feedback = ""
				m.wordIdx++
				if m.wordIdx >= len(m.drillList) {
					m.drillList = m.redoList
					m.redoList = nil
					m.wordIdx = 0
				}
				m.ti.SetValue("")
			} else {
				want := strings.TrimSpace(m.drillList[m.wordIdx].French)
				got := m.ti.Value()
				if want == got {
					m.feedback = "Correct!"
				} else {
					m.redoList = append(m.redoList, m.drillList[m.wordIdx])
					m.feedback = fmt.Sprintf("Incorrect. The correct answer is: %s, you wrote: %s", want, got)
				}
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.ti, cmd = m.ti.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.renderExercise()
}

func (m *Model) renderExercise() string {
	q := fmt.Sprintf(`
		%s

		Word: %s / %s
		
		Your answer: %s`,
		renderFeedbackStyle(m.feedback),
		m.drillList[m.wordIdx].English,
		m.drillList[m.wordIdx].German,
		m.ti.View(),
	)

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
