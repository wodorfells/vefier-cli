package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wodorfells/vefier-cli/internal/config"
	"github.com/wodorfells/vefier-cli/internal/keys"
	"github.com/wodorfells/vefier-cli/internal/ui/styles"
	"github.com/wodorfells/vefier-cli/internal/ui/views"
)

type App struct {
	cfg           *config.Config
	keyStore      *keys.Store
	
	state         string // "providers", "chat", "checker", "palette", "compare"
	previousState string
	providersView *views.ProvidersView
	chatView      *views.ChatView
	checkerView   *views.CheckerView
	paletteView   *views.PaletteView
	compareView   *views.CompareView
}

func NewApp(cfg *config.Config) *App {
	ks := keys.NewStore()
	return &App{
		cfg:           cfg,
		keyStore:      ks,
		state:         "providers",
		providersView: views.NewProvidersView(),
		checkerView:   views.NewCheckerView(ks),
		paletteView:   views.NewPaletteView(),
		compareView:   views.NewCompareView(),
	}
}

func (a *App) Init() tea.Cmd {
	return a.providersView.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		}
		
		if msg.String() == "ctrl+k" {
			a.previousState = a.state
			a.state = "palette"
			return a, a.paletteView.Init()
		}
		
		if msg.String() == "ctrl+t" {
			styles.CycleTheme()
			// Forces re-render across the app with new styles
			return a, nil
		}

		if a.state == "palette" {
			var cmd tea.Cmd
			tempView, cmd := a.paletteView.Update(msg)
			if tempView != nil {
				a.paletteView = tempView
			}
			
			if a.paletteView.Action != "" {
				action := a.paletteView.Action
				a.paletteView.Action = ""
				
				switch action {
				case "theme":
					styles.CycleTheme()
					a.state = a.previousState
				case "esc":
					a.state = a.previousState
				case "quit":
					return a, tea.Quit
				case "providers":
					a.state = "providers"
				case "checker":
					a.state = "checker"
				case "compare":
					a.state = "compare"
				}
			}
			return a, cmd
		}

		if a.state == "providers" {
			if msg.String() == "enter" {
				p := a.providersView.GetSelected()
				a.chatView = views.NewChatView(p, a.keyStore)
				a.state = "chat"
				return a, a.chatView.Init()
			}
			if msg.String() == "c" {
				a.state = "checker"
				return a, a.checkerView.Init()
			}
		} else if a.state == "chat" {
			if msg.String() == "esc" {
				if !a.chatView.IsStreaming() {
					a.state = "providers"
					return a, nil
				}
			}
		} else if a.state == "checker" {
			if msg.String() == "c" || msg.String() == "esc" {
				a.state = "providers"
				return a, nil
			}
		} else if a.state == "compare" {
			if msg.String() == "esc" {
				a.state = "providers"
				return a, nil
			}
		}
	}

	if a.state == "providers" {
		var cmd tea.Cmd
		a.providersView, cmd = a.providersView.Update(msg)
		cmds = append(cmds, cmd)
	} else if a.state == "chat" {
		var cmd tea.Cmd
		tempModel, cmd := a.chatView.Update(msg)
		if tempModel != nil {
			a.chatView = tempModel
		}
		cmds = append(cmds, cmd)
	} else if a.state == "checker" {
		var cmd tea.Cmd
		tempModel, cmd := a.checkerView.Update(msg)
		if tempModel != nil {
			a.checkerView = tempModel
		}
		cmds = append(cmds, cmd)
	} else if a.state == "compare" {
		var cmd tea.Cmd
		tempModel, cmd := a.compareView.Update(msg)
		if tempModel != nil {
			a.compareView = tempModel
		}
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a *App) View() string {
	if a.state == "palette" && a.paletteView != nil {
		return a.paletteView.View()
	}
	if a.state == "providers" {
		return a.providersView.View() + "\n[c] Менеджер ключей | [Ctrl+K] Палитра\n"
	}
	if a.state == "chat" && a.chatView != nil {
		return a.chatView.View()
	}
	if a.state == "checker" && a.checkerView != nil {
		return a.checkerView.View()
	}
	if a.state == "compare" && a.compareView != nil {
		return a.compareView.View()
	}
	return "Загрузка..."
}

