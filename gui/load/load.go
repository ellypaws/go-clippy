package load

// A simple example that shows how to render an animated progress bar. In this
// example we bump the progress by 25% every two seconds, animating our
// progress bar to its new target state.
//
// It's also possible to render a progress bar in a more static fashion without
// transitions. For details on that approach see the progress-static example.

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func New() *Model {
	m := Model{
		progress: progress.New(progress.WithDefaultGradient()),
	}
	return &m
}

type Model struct {
	progress progress.Model
	visible  bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case Progress:
		if m.visible != msg.Show {
			m.visible = msg.Show
		}
		cmd := m.progress.SetPercent(msg.Percent)
		return m, cmd
	case Goal:
		if msg.Current == msg.Total {
			if m.visible {
				m.visible = false
			}
			return m, nil
		}
		if m.visible != msg.Show {
			m.visible = msg.Show
		}
		cmd := m.progress.SetPercent(float64(msg.Current) / float64(msg.Total))
		return m, cmd
	default:
		return m, nil
	}
}

func (m *Model) Toggle() *Model {
	m.visible = !m.visible
	return m
}

func (m Model) Visible() bool {
	return m.visible
}

func (m Model) View() string {
	if !m.visible {
		return "Loaded"
	}
	return m.progress.View()
}

type Progress struct {
	Percent float64
	Show    bool
}

type Goal struct {
	Current int
	Total   int
	Show    bool
}
