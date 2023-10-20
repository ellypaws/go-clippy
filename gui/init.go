package gui

import (
	"github.com/charmbracelet/bubbletea"
	"go-clippy/gui/help"
	"go-clippy/gui/log"
	"go-clippy/gui/table"
)

type Model struct {
	logger *logger.Model
	table  *table.Model
	help   help.Model
	height int
	width  int
}

func (m Model) Init() tea.Cmd {
	return m.logger.Init()
}

func NewModel() Model {
	return Model{
		logger: logger.NewSender(),
		table:  table.New(),
		help:   help.New(),
	}
}
