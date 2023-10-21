package main

import (
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	discord "go-clippy/bot"
	"go-clippy/gui"
	"log"
)

func main() {
	zone.NewGlobal()
	p := tea.NewProgram(gui.NewModel(),
		//tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithMouseAllMotion(),
	)
	go discord.GetBot().Run(p)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}
