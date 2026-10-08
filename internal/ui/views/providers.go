package views

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wodorfells/vefier-cli/internal/providers"
	"github.com/wodorfells/vefier-cli/internal/ui/styles"
)

type ProvidersView struct {
	providers []providers.Provider
	cursor    int
}

func NewProvidersView() *ProvidersView {
	return &ProvidersView{
		providers: providers.GetCatalog(),
	}
}

func (v *ProvidersView) Init() tea.Cmd {
	return nil
}

func (v *ProvidersView) Update(msg tea.Msg) (*ProvidersView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.cursor > 0 {
				v.cursor--
			}
		case "down", "j":
			if v.cursor < len(v.providers)-1 {
				v.cursor++
			}
		}
	}
	return v, nil
}

func (v *ProvidersView) View() string {
	var b strings.Builder
	
	ascii := `
 ___      ___    _______  ________  ___  _______   ________     
|"  \    /"  |  /"     "|/"       )|"  |/"     "| /"       )    
 \   \  //   | (: ______)(:   \___/ ||  (: ______)(:   \___/    
  \\  \/.    |  \/    |   \___  \   |:  |\/    |   \___  \      
   \.    //  |  // ___)_   __/  \\  |.  |// ___)_   __/  \\     
    \\   /   | (:      "| /" \   :) /\  |(:      "|/" \   :)    
     \__/    |  \_______)(_______/ (__\_)\_______)(_______/     
`
	b.WriteString(styles.AppBanner.Render("VeFier CLI " + string(rune(128640)))) 
	b.WriteString(styles.Title.Render(ascii))
	b.WriteString("\nКаталог провайдеров\n\n")

	for i, p := range v.providers {
		cursor := "  "
		style := styles.ListNormal
		if i == v.cursor {
			cursor = "> "
			style = styles.ListSelected
		}
		
		row := fmt.Sprintf("%s%s (%s)", cursor, p.Name, p.BaseURL)
		b.WriteString(style.Render(row))
		b.WriteString("\n")
	}
	
	b.WriteString("\nНажмите 'Enter' для запуска чата, 'q' для выхода.\n")
	return b.String()
}

func (v *ProvidersView) GetSelected() providers.Provider {
    return v.providers[v.cursor]
}
