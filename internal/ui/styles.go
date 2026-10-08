package ui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name   string
	Accent string
	Text   string
	Bg     string
	Border string
}

var Themes = []Theme{
	{Name: "crush", Accent: "#c86bff", Text: "#e6e6f2", Border: "#6b50ff"},
	{Name: "ocean", Accent: "#4fc3f7", Text: "#e3f2fd", Border: "#2962ff"},
	{Name: "matrix", Accent: "#69f0ae", Text: "#d7ffd9", Border: "#00c853"},
	{Name: "dracula", Accent: "#ff79c6", Text: "#f8f8f2", Border: "#bd93f9"},
}

func getTheme(name string) Theme {
	for _, t := range Themes {
		if t.Name == name {
			return t
		}
	}
	return Themes[0]
}

func NextTheme(current string) string {
	for i, t := range Themes {
		if t.Name == current {
			next := (i + 1) % len(Themes)
			return Themes[next].Name
		}
	}
	return Themes[0].Name
}

// GenerateStyles returns fresh styles based on the theme
func GenerateStyles(themeName string) map[string]lipgloss.Style {
	t := getTheme(themeName)
	
	s := make(map[string]lipgloss.Style)
	s["title"] = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Bold(true)
	s["text"] = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text))
	s["border"] = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(t.Border))
	s["user"] = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text)).Bold(true).Background(lipgloss.Color(t.Border)).Padding(0, 1).MarginBottom(1)
	s["ai"] = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Bold(true).MarginBottom(1)
	s["help"] = lipgloss.NewStyle().Foreground(lipgloss.Color("#707070"))
	
	return s
}