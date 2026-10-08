package styles

import "github.com/charmbracelet/lipgloss"

var (
	ThemeCrushText lipgloss.TerminalColor
	ThemeCrushAcc  lipgloss.TerminalColor
	
	Title        lipgloss.Style
	ListSelected lipgloss.Style
	ListNormal   lipgloss.Style
	AppBanner    lipgloss.Style
	ErrorText    lipgloss.Style
)

func init() {
	ApplyTheme("crush")
}

func ApplyTheme(name string) {
	var p Palette
	for _, t := range Themes {
		if t.Name == name {
			p = t
			break
		}
	}
	if p.Name == "" {
		p = Themes[0]
	}

	ThemeCrushText = lipgloss.Color(p.Text)
	ThemeCrushAcc = lipgloss.Color(p.Accent)

	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(ThemeCrushAcc).
		MarginBottom(1)
		
	ListSelected = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Ok)).
		Bold(true)
		
	ListNormal = lipgloss.NewStyle().
		Foreground(ThemeCrushText)
		
	ErrorText = lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.Bad)).
		Bold(true)
		
	AppBanner = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(p.Text)).
		Background(lipgloss.Color(p.GradB)).
		Padding(0, 2).
		MarginBottom(1)
}

