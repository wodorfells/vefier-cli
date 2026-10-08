package views

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vefier/vefier-cli/internal/ui/styles"
)

type Item struct {
	title, desc, action string
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title + " " + i.desc }
func (i Item) Action() string      { return i.action }

type PaletteView struct {
	Model  list.Model
	Action string
}

func NewPaletteView() *PaletteView {
	items := []list.Item{
		Item{title: "Отменить и вернуться", desc: "Скрыть палитру", action: "esc"},
		Item{title: "Переключить тему (Ctrl+T)", desc: "Сменить визуальное оформление", action: "theme"},
		Item{title: "Каталог провайдеров", desc: "Главный экран", action: "providers"},
		Item{title: "Чекер ключей", desc: "Проверка баланса и валидности", action: "checker"},
		Item{title: "Сравнение моделей", desc: "Параллельный запрос к 2 провайдерам", action: "compare"},
		Item{title: "Выход", desc: "Закрыть VeFier", action: "quit"},
	}

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = styles.ListSelected
	d.Styles.SelectedDesc = styles.ListSelected.Copy().Faint(true)
	
	m := list.New(items, d, 60, 15)
	m.Title = "Командная палитра (Ctrl+K)"
	m.SetShowStatusBar(false)
	m.SetFilteringEnabled(true)

	return &PaletteView{Model: m}
}

func (v *PaletteView) Init() tea.Cmd {
	return nil
}

func (v *PaletteView) Update(msg tea.Msg) (*PaletteView, tea.Cmd) {
	v.Action = ""
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if i, ok := v.Model.SelectedItem().(Item); ok {
				v.Action = i.Action()
				return v, nil
			}
		}
	}

	var cmd tea.Cmd
	v.Model, cmd = v.Model.Update(msg)
	return v, cmd
}

func (v *PaletteView) View() string {
	return "\n  " + v.Model.View()
}