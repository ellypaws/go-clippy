package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A simple example that shows how to send messages to a Bubble Tea program
// from outside the program using Program.Send(Msg).

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	dotStyle      = helpStyle.Copy().UnsetMargins()
	durationStyle = dotStyle.Copy()
	appStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type ResultMsg struct {
	time    time.Time
	message string
}

func MessageF(format string, a ...any) ResultMsg {
	return Message(fmt.Sprintf(format, a...))
}

func Message(msg string, t ...time.Time) ResultMsg {
	if len(t) == 0 {
		t = []time.Time{time.Now()}
	}
	return ResultMsg{time: t[0], message: msg}
}

func (r ResultMsg) String() string {
	if r.message == "" {
		return dotStyle.Render(strings.Repeat(".", 30))
	}
	return fmt.Sprintf("[%s] %s",
		// show as hour i.e. 15:23
		durationStyle.Render(r.time.Format("15:04")),
		r.message,
	)
}

type Model struct {
	spinner     spinner.Model
	results     []ResultMsg
	quitting    bool
	LastResults int
}

func NewSender() *Model {
	const numLastResults = 15
	s := spinner.New()
	s.Style = spinnerStyle
	return &Model{
		spinner:     s,
		results:     make([]ResultMsg, numLastResults),
		LastResults: numLastResults,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case ResultMsg:
		m.results = append(m.results[1:], msg)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m *Model) View() string {
	var s string

	if m.quitting {
		s += "That’s all for today!"
	} else {
		s += m.spinner.View() + " Eating food..."
	}

	s += "\n\n"

	for _, res := range m.results {
		s += res.String() + "\n"
	}

	if !m.quitting {
		s += helpStyle.Render("Press any key to exit")
	}

	if m.quitting {
		s += "\n"
	}

	return appStyle.Render(s)
}
