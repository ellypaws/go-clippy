package gui

import (
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"strings"
)

func (m Model) View() string {
	out := strings.Builder{}
	if m.table.Visible() {
		out.WriteString(
			lipgloss.JoinHorizontal(lipgloss.Center,
				m.logger.View(),
				m.table.View(),
			))
	} else {
		out.WriteString(m.logger.View())
	}
	if m.progress.Visible() {
		out.WriteString(m.progress.View())
	}
	out.WriteString(m.help.View())
	return zone.Scan(out.String())
}
