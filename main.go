package main

import (
	tea "github.com/charmbracelet/bubbletea"
	discord "go-clippy/bot"
)

type model struct {
	Bot *discord.BOT
}

func main() {
	m := model{
		Bot: discord.GetBot(),
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

func (m model) Init() tea.Cmd {
	m.Bot.Run()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m model) View() string {
	return ""
}
