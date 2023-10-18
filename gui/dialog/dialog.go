// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package dialog

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"winsleepd/cmd/tui/service"
	"winsleepd/cmd/tui/table"
)

var (
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0)

	normal = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#888B7E")).
		Padding(0, 3).
		MarginTop(1).
		MarginRight(2)

	activated = normal.Copy().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#F25D94")).
			MarginRight(2).
			Underline(true)

	hover = activated.Copy().
		Background(lipgloss.Color("#FFC09C"))
)

type DialogModel struct {
	Id     string
	Height int
	width  int

	Active   string
	hover    string
	Question string
}

func (m DialogModel) Init() tea.Cmd {
	return nil
}

func (m DialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case table.Query:
		if service.Get().IsInstalled() {
			m.Active = "confirm"
		} else {
			m.Active = "cancel"
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.MouseMsg:
		switch {
		case zone.Get(m.Id + "confirm").InBounds(msg):
			m.hover = "confirm"
		case zone.Get(m.Id + "cancel").InBounds(msg):
			m.hover = "cancel"
		default:
			m.hover = ""
		}
		if msg.Type == tea.MouseLeft {
			if zone.Get(m.Id + "confirm").InBounds(msg) {
				service.Get().Install(false)
				m.Active = "confirm"
			} else if zone.Get(m.Id + "cancel").InBounds(msg) {
				service.Get().Uninstall()
				m.Active = "cancel"
			}
		}
	}
	return m, nil
}

func (m DialogModel) View() string {
	var okButton, cancelButton string

	okButton = normal.Render("Yes")
	cancelButton = normal.Render("No")

	switch m.Active {
	case "confirm":
		okButton = activated.Render("✨")
	case "cancel":
		cancelButton = activated.Render("🌱")
	}

	switch m.hover {
	case "confirm":
		if m.Active == "confirm" {
			okButton = hover.Render("✨")
		} else {
			okButton = hover.Render("Yes")
		}
	case "cancel":
		if m.Active == "cancel" {
			cancelButton = hover.Render("🌱")
		} else {
			cancelButton = hover.Render("No")
		}
	}

	question := lipgloss.NewStyle().Width(27).Align(lipgloss.Center).Render("Install winsleepd service?")
	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
		zone.Mark(m.Id+"confirm", okButton),
		zone.Mark(m.Id+"cancel", cancelButton),
	)
	return dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, question, buttons))
}
