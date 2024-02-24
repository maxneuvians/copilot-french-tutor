package pkg

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

func renderVP(width, height int) viewport.Model {
	vp := viewport.Model{}
	vp.Width = width
	vp.Height = height
	vp.SetContent("width: " + fmt.Sprint(width) + " height: " + fmt.Sprint(height))
	vp.Style = lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)
	return vp
}
