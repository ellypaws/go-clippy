package gui

import (
	"github.com/charmbracelet/bubbles/key"
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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
		switch {
		case key.Matches(msg, m.help.Keys.Settings):
			m.table = m.table.Toggle()
			if m.table.Visible() {
				m.logger.LastResults = 5
			} else {
				m.logger.LastResults = 15
			}
			return m, nil
		}
	}
	return m.propagate(msg)
	//return m, nil
}

func (m Model) propagate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd
	model, cmd := m.table.Update(msg)
	commands = append(commands, cmd)
	m.table = &model
	m.help, cmd = m.help.Update(msg)
	return m, cmd
}
