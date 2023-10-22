package gui

import (
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"strings"
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}
	out := strings.Builder{}
	out.WriteString(
		lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinHorizontal(lipgloss.Center,
				lipgloss.JoinVertical(lipgloss.Center,
					m.progress.View(),
					m.logger.View(),
				),
				lipgloss.JoinVertical(lipgloss.Center,
					m.table.View(),
					m.input.View(),
				),
			)))
	out.WriteString(m.help.View())
	return zone.Scan(out.String())
}
