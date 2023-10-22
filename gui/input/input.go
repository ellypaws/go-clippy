package input

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

const (
	TokenInput = iota
	PointsInput
	PasswordInput
	Submit
)

type Model struct {
	active, focus int
	inputs        []textinput.Model
	cursorMode    cursor.Mode
	setPoints     Points
	enabled       bool
}

func New() *Model {
	m := Model{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case TokenInput:
			t.Placeholder = "Token"
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case PointsInput:
			t.Placeholder = "Points"
			t.CharLimit = 4
		case PasswordInput:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return &m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.active == len(m.inputs) {
				return m.newPoints()
			}
			//
			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.active = min(m.focus, Submit)
			} else {
				m.active = Submit
			}

			//if m.active > len(m.inputs) {
			//	m.active = 0
			//} else if m.active < 0 {
			//	m.active = len(m.inputs)
			//}

			if m.focus < 0 {
				m.focus = 0
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focus {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *Model) ReceivePoints(points string, snowflake string) tea.Cmd {
	m.enabled = true
	m.active = PointsInput
	m.focus = PointsInput
	m.inputs[PointsInput].Placeholder = points
	m.inputs[PointsInput].Reset()
	m.inputs[PointsInput].Focus()
	m.setPoints = Points{
		Points:    points,
		Snowflake: snowflake,
	}
	return nil
}

func (m *Model) newPoints() (Model, tea.Cmd) {
	m.enabled = false
	if m.inputs[PointsInput].Value() == "" {
		return *m, nil
	}
	m.setPoints.Points = m.inputs[PointsInput].Value()
	m.inputs[PointsInput].Reset()
	m.inputs[PointsInput].Blur()
	return *m, func() tea.Msg {
		return m.setPoints
	}
}

type Points struct {
	Points    string
	Snowflake string
}

func (m Model) View() string {
	if !m.enabled {
		return ""
	}
	var b strings.Builder

	b.WriteString(m.inputs[min(m.focus, len(m.inputs)-1)].View())

	button := &blurredButton
	if m.active == len(m.inputs) {
		button = &focusedButton
	}
	//b.WriteString(fmt.Sprintf("%v:%v", m.focus, m.active))
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	//b.WriteString(helpStyle.Render("cursor mode is "))
	//b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	//b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}