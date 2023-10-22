package gui

import (
	"github.com/charmbracelet/bubbletea"
	"go-clippy/gui/help"
	"go-clippy/gui/input"
	"go-clippy/gui/load"
	"go-clippy/gui/log"
	"go-clippy/gui/table"
)

type Model struct {
	logger   *logger.Model
	table    *table.Model
	progress *load.Model
	input    *input.Model
	help     help.Model
	height   int
	width    int
}

func (m Model) Init() tea.Cmd {
	var commands []tea.Cmd
	commands = append(commands, m.logger.Init())
	commands = append(commands, m.table.Init())
	commands = append(commands, m.progress.Init())
	return tea.Batch(commands...)
}

func NewModel() Model {
	return Model{
		logger:   logger.NewSender(),
		table:    table.New(),
		progress: load.New(),
		input:    input.New(),
		help:     help.New(),
	}
}
