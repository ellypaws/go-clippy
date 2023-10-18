package gui

import (
	zone "github.com/lrstanley/bubblezone"
	"strings"
)

func (m Model) View() string {
	out := strings.Builder{}
	out.WriteString(m.logger.View())
	return zone.Scan(out.String())
}
