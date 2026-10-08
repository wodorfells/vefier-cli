package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vefier/vefier-cli/internal/config"
	"github.com/vefier/vefier-cli/internal/ui"
)

func main() {
	cfg := config.LoadConfig()
	
	app := ui.NewApp(cfg)
	p := tea.NewProgram(app, tea.WithAltScreen())
	
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ошибка запуска VeFier: %v\n", err)
		os.Exit(1)
	}
}