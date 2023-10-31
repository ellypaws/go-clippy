package table

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go-clippy/database/clippy"
	"go-clippy/gui/input"
	"strconv"
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
			if m.table.Focused() {
				m.table.Blur()
				return m, promptPoints(m.table.SelectedRow()[3], m.table.SelectedRow()[2], m.table.SelectedRow()[4])
			}
		case "r":
			go clippy.GetCache().SynchronizeAllPoints()
			return m, nil
		}
	case Leaderboard, clippy.Sync:
		m.table = m.UpdateLeaderboard(15)
		return m, nil
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) SetPoints(points string, snowflake string, private string) Model {
	cur := m.table.Rows()
	i := m.table.Cursor()
	//cur[i][2] = snowflake
	cur[i][3] = points
	cur[i][4] = private
	m.table.SetRows(cur)
	m.table.Focus()
	return m
}

func promptPoints(points string, snowflake string, private string) tea.Cmd {
	return func() tea.Msg {
		return PromptUser{
			Points:    points,
			Snowflake: snowflake,
			Private:   private,
		}
	}
}

type PromptUser struct {
	Points    string
	Snowflake string
	Private   string
}

func (m *Model) Toggle() *Model {
	m.visible = !m.visible
	if m.visible {
		m.table.Focus()
	} else {
		m.table.Blur()
	}
	m.table.UpdateViewport()
	return m
}

func (m Model) Visible() bool {
	return m.visible
}

func (m Model) View() string {
	if !m.visible {
		return baseStyle.Render("Press 's' to adjust user points.")
	}
	return baseStyle.Render(m.table.View()) + "\n"
}

type Leaderboard struct{}

func (m Model) UpdateLeaderboard(max int) table.Model {
	m.table.SetRows(clippy.GetCache().LeaderboardTable(max, clippy.Request{}))
	return m.table
}

func (m Model) UpdateUser(points input.Points) table.Model {
	user, ok := clippy.GetCache().GetConfig(points.Snowflake)
	if !ok {
		return m.table
	}
	pointsInt, _ := strconv.Atoi(points.Points)
	user.Points = pointsInt
	user.Record()
	return m.UpdateLeaderboard(len(m.table.Rows()))
}

func New() *Model {
	columns := []table.Column{
		{Title: "#", Width: 2},
		{Title: "Username", Width: 12},
		{Title: "Snowflake", Width: 12},
		{Title: "Points", Width: 8},
		{Title: "Private", Width: 8},
	}

	rows := []table.Row{
		{"1", "Username", "000000", "00", "false"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithKeyMap(table.DefaultKeyMap()),
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
