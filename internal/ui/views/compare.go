package views

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/wodorfells/vefier-cli/internal/providers"
	"github.com/wodorfells/vefier-cli/internal/ui/styles"
)

type CompareView struct {
	providers []providers.Provider
	results   []string
	loading   bool
}

func NewCompareView() *CompareView {
	cat := providers.GetCatalog()
	prots := cat
	if len(cat) > 2 {
		prots = cat[:2]
	}
	return &CompareView{
		providers: prots,
		results:   []string{"Нажмите Enter для запуска", "Ожидание..."},
	}
}

func (v *CompareView) Init() tea.Cmd {
	return nil
}

type CompareResultMsg struct {
	Index int
	Text  string
}

func fakeWait(idx int, delay time.Duration, res string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(delay)
		return CompareResultMsg{Index: idx, Text: res}
	}
}

func (v *CompareView) Update(msg tea.Msg) (*CompareView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" && !v.loading {
			v.loading = true
			v.results[0] = "Стриминг: токены идут..."
			v.results[1] = "Стриминг: токены идут..."
			return v, tea.Batch(
				fakeWait(0, 1*time.Second, "Ответ от модели сгенерирован.\n\n$0.02 | 400 tok | 0.8s"),
				fakeWait(1, 2*time.Second, "Ответ от модели сгенерирован.\n\n$0.04 | 400 tok | 1.8s"),
			)
		}
	case CompareResultMsg:
		v.results[msg.Index] = msg.Text
		if msg.Index == 1 {
			v.loading = false
		}
	}
	return v, nil
}

func (v *CompareView) View() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ThemeCrushAcc).
		Width(35).Height(15).Padding(1)

	var boxes []string
	for i, p := range v.providers {
		title := styles.ListSelected.Render(p.Name)
		content := title + "\n\n" + v.results[i]
		boxes = append(boxes, boxStyle.Render(content))
	}

	layout := lipgloss.JoinHorizontal(lipgloss.Top, boxes...)
	header := styles.Title.Render("Сравнение моделей (A/B Test)")
	
	footer := "Enter = Запуск | Esc = Каталог | Ctrl+K = Меню"
	if v.loading {
		footer = "Идет генерация..."
	}
	
	return header + "\n\n" + layout + "\n\n" + styles.ListNormal.Render(footer)
}
