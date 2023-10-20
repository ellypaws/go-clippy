package gui

import (
	zone "github.com/lrstanley/bubblezone"
	"strings"
)

func (m Model) View() string {
	out := strings.Builder{}
	out.WriteString(m.logger.View())
	if m.table.Visible() {
		out.WriteString(m.table.View())
	}
	out.WriteString(m.help.View())
	return zone.Scan(out.String())
}
