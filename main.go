package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"vefier-cli/config" // Если имя модуля в go.mod другое, укажите ваше (например: "github.com/wodorfells/vefier-cli/config")
	"vefier-cli/ui"     // Аналогично укажите правильный путь к пакету ui
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
