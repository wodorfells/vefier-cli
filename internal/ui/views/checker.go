package views

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vefier/vefier-cli/internal/api"
	"github.com/vefier/vefier-cli/internal/keys"
	"github.com/vefier/vefier-cli/internal/providers"
	"github.com/vefier/vefier-cli/internal/ui/styles"
)

type CheckResultMsg struct {
	Index  int
	Status string
}

type CheckerView struct {
	keysStore *keys.Store
	catalog   []providers.Provider
	cursor    int
	checking  bool
}

func NewCheckerView(ks *keys.Store) *CheckerView {
	return &CheckerView{
		keysStore: ks,
		catalog:   providers.GetCatalog(),
	}
}

func (v *CheckerView) Init() tea.Cmd {
	return nil
}

func checkKeyCmd(idx int, k keys.Key, p providers.Provider) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(500 * time.Millisecond) // Визуальная симуляция процесса
		status := api.CheckKey(context.Background(), p, k)
		return CheckResultMsg{Index: idx, Status: status}
	}
}

func (v *CheckerView) StartCheck() tea.Cmd {
	v.checking = true
	var cmds []tea.Cmd
	
	v.keysStore.Mu().Lock()
	for i, k := range v.keysStore.Keys {
		v.keysStore.Keys[i].Status = "checking..."
		
		var prov providers.Provider
		for _, cat := range v.catalog {
			if cat.ID == k.Provider {
				prov = cat
				break
			}
		}
		cmds = append(cmds, checkKeyCmd(i, k, prov))
	}
	v.keysStore.Mu().Unlock()
	
	return tea.Batch(cmds...)
}

func (v *CheckerView) Update(msg tea.Msg) (*CheckerView, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if v.cursor > 0 { v.cursor-- }
		case "down", "j":
			v.keysStore.Mu().RLock()
			kLen := len(v.keysStore.Keys)
			v.keysStore.Mu().RUnlock()
			if v.cursor < kLen-1 { v.cursor++ }
		case "r":
			if !v.checking {
				return v, v.StartCheck()
			}
		case "i":
			v.keysStore.AddKeys(keys.ParseKeys("Anthropic Demo = sk-ant-api03\nTako Key = tako_sk_1234\nBadKey=sk-invalid"))
		}
	
	case CheckResultMsg:
		v.keysStore.Mu().Lock()
		if msg.Index < len(v.keysStore.Keys) {
			v.keysStore.Keys[msg.Index].Status = msg.Status
		}
		
		allDone := true
		for _, k := range v.keysStore.Keys {
			if k.Status == "checking..." {
				allDone = false
			}
		}
		v.checking = !allDone
		v.keysStore.Mu().Unlock()
	}
	return v, nil
}

func (v *CheckerView) View() string {
	var b strings.Builder
	b.WriteString(styles.Title.Render("Массовая проверка ключей"))
	b.WriteString("\n\n")

	v.keysStore.Mu().RLock()
	for i, k := range v.keysStore.Keys {
		cursor := "  "
		style := styles.ListNormal
		if i == v.cursor {
			cursor = "> "
			style = styles.ListSelected
		}
		
		var statusColor string
		switch k.Status {
		case "active": statusColor = "🟢 Рабочий"
		case "invalid": statusColor = "🔴 Ошибка"
		case "untested": statusColor = "⚪ Не проверен"
		case "checking...": statusColor = "🔄 Проверка..."
		default: statusColor = "🟡 " + k.Status
		}
		
		row := fmt.Sprintf("%s %-20s | %-12s | %s", cursor, k.Name, k.Provider, statusColor)
		b.WriteString(style.Render(row))
		b.WriteString("\n")
	}
	v.keysStore.Mu().RUnlock()

	b.WriteString("\n")
	b.WriteString(styles.ListNormal.Render("[r] Проверить все  [i] Демо-импорт  [c] Каталог  [q] Выход"))
	return b.String()
}