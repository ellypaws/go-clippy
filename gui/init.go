package gui

import (
	"github.com/charmbracelet/bubbletea"
	"go-clippy/gui/log"
)

type Model struct {
	logger *logger.Model
	height int
	width  int
}

func (m Model) Init() tea.Cmd {
	return m.logger.Init()
}

func NewModel() Model {
	return Model{
		logger: logger.NewSender(),
	}
}
