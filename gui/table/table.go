package table

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	table   table.Model
	visible bool
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *Model) Toggle() *Model {
	m.visible = !m.visible
	m.table.Blur()
	m.table.UpdateViewport()
	return m
}

func (m Model) Visible() bool {
	return m.visible
}

func (m Model) View() string {
	if !m.visible {
		return baseStyle.Render("")
	}
	return baseStyle.Render(m.table.View()) + "\n"
}

func New() *Model {
	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Username", Width: 12},
		{Title: "Snowflake", Width: 12},
		{Title: "Points", Width: 5},
	}

	rows := []table.Row{
		{"1", "Username", "000000", "15"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := Model{table: t}
	return &m
}
