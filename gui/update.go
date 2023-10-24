package gui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"go-clippy/gui/input"
	"go-clippy/gui/load"
	"go-clippy/gui/log"
	"go-clippy/gui/table"
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
		switch {
		case key.Matches(msg, m.help.Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.help.Keys.Settings):
			if m.input.Enabled() {
				return m.propagate(msg)
			}
			m.table = m.table.Toggle()
			if m.table.Visible() {
				m.logger.LastResults = 5
			} else {
				m.logger.LastResults = 15
			}
			tbl, cmd := m.table.Update(tea.Msg(table.Leaderboard{}))
			m.table = &tbl
			return m, cmd
		}
	case progress.FrameMsg, load.Progress, load.Goal:
		*m.progress, cmd = m.progress.Update(msg)
		return m, cmd
	case table.PromptUser:
		//m.table.Toggle()
		m.input.ReceivePoints(msg.Points, msg.Snowflake, msg.Private)
		return m, nil
	case input.Points:
		*m.table = m.table.SetPoints(msg.Points, msg.Snowflake, msg.Private)
	}
	return m.propagate(msg)
	//return m, nil
}

func (m Model) propagate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd
	var cmd tea.Cmd
	*m.table, cmd = m.table.Update(msg)
	commands = append(commands, cmd)
	m.help, cmd = m.help.Update(msg)
	commands = append(commands, cmd)
	*m.progress, cmd = m.progress.Update(msg)
	commands = append(commands, cmd)
	*m.input, cmd = m.input.Update(msg)
	commands = append(commands, cmd)
	return m, tea.Batch(commands...)
}
