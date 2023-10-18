package gui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"go-clippy/gui/log"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		msg.Height -= 2
		msg.Width -= 4
		return m, nil
	case spinner.TickMsg, logger.ResultMsg:
		m.logger, cmd = m.logger.Update(msg)
		return m, cmd
	}
	return m, nil
}
