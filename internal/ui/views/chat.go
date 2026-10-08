package views

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/wodorfells/vefier-cli/internal/api"
	"github.com/wodorfells/vefier-cli/internal/keys"
	"github.com/wodorfells/vefier-cli/internal/providers"
	"github.com/wodorfells/vefier-cli/internal/sessions"
	"github.com/wodorfells/vefier-cli/internal/tools"
	"github.com/wodorfells/vefier-cli/internal/ui/styles"
)

type StreamMsg struct {
	Chunk api.StreamChunk
	ch    <-chan api.StreamChunk
}

type ChatView struct {
	provider    providers.Provider
	keyStore    *keys.Store
	session     *sessions.Session
	sessionStore *sessions.Store
	tools       *tools.Registry
	
	viewport    viewport.Model
	textarea    textarea.Model
	messages    []string
	isStreaming bool
	cancel      context.CancelFunc
}

func NewChatView(p providers.Provider, ks *keys.Store) *ChatView {
	ta := textarea.New()
	ta.Placeholder = "Сообщение (Ctrl+S отправка, Ctrl+P права, Esc возврат)..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 10000
	ta.SetWidth(80)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(80, 20)
	vp.SetContent("Чат инициирован.\n")

	sStore := sessions.NewStore()
	sess := sStore.CreateSession(p.ID, "gpt-4")

	return &ChatView{
		provider:     p,
		keyStore:     ks,
		sessionStore: sStore,
		session:      sess,
		tools:        tools.NewRegistry(),
		textarea:     ta,
		viewport:     vp,
		messages:     []string{},
	}
}

func (v *ChatView) Init() tea.Cmd {
	return textarea.Blink
}

func (v *ChatView) IsStreaming() bool {
	return v.isStreaming
}

func listenStream(ch <-chan api.StreamChunk) tea.Cmd {
	return func() tea.Msg {
		chunk, ok := <-ch
		if !ok {
			return StreamMsg{Chunk: api.StreamChunk{Done: true}, ch: ch}
		}
		return StreamMsg{Chunk: chunk, ch: ch}
	}
}

func (v *ChatView) Update(msg tea.Msg) (*ChatView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if v.isStreaming && msg.String() == "esc" {
			if v.cancel != nil {
				v.cancel()
				v.isStreaming = false
			}
			return v, nil
		}

		switch msg.String() {
		case "ctrl+s":
			if v.isStreaming {
				break
			}
			text := strings.TrimSpace(v.textarea.Value())
			if text == "" {
				break
			}
			
			// Сохраняем в историю сессии
			v.sessionStore.AddMessage(v.session.ID, "user", text, len(text)/4, 0)
			
			v.messages = append(v.messages, lipgloss.NewStyle().Foreground(styles.ThemeCrushAcc).Render("Вы: ")+text)
			v.messages = append(v.messages, styles.ListSelected.Render(v.provider.Name+": "))
			v.viewport.SetContent(strings.Join(v.messages, "\n\n"))
			v.viewport.GotoBottom()
			v.textarea.Reset()
			
			v.isStreaming = true
			key := v.keyStore.GetWorkingKey(v.provider.ID)
			token := "mock-token"
			if key != nil {
				token = key.Value
			}
			
			var ctx context.Context
			ctx, v.cancel = context.WithCancel(context.Background())
			ch := make(chan api.StreamChunk)
			
			go api.StreamCompletion(ctx, v.provider, token, v.session.Model, text, ch)
			cmds = append(cmds, listenStream(ch))
			
		case "ctrl+p":
			// Переключаем режим прав для инструментов
			v.tools.ToggleMode()

		case "ctrl+c":
			return v, tea.Quit
		}
		
	case StreamMsg:
		if msg.Chunk.Error != nil {
			v.messages[len(v.messages)-1] += fmt.Sprintf("\n[Ошибка: %v]", msg.Chunk.Error)
			v.isStreaming = false
		} else if msg.Chunk.Done {
			v.isStreaming = false
			// Мокаем добавление токенов/стоимости (в реальном API приходит в [DONE] или usage)
			v.sessionStore.AddMessage(v.session.ID, "assistant", "generated", 150, 0.002)
		} else {
			v.messages[len(v.messages)-1] += msg.Chunk.Content
			cmds = append(cmds, listenStream(msg.ch))
		}
		v.viewport.SetContent(strings.Join(v.messages, "\n\n"))
		v.viewport.GotoBottom()
	}

	var cmd tea.Cmd
	v.textarea, cmd = v.textarea.Update(msg)
	cmds = append(cmds, cmd)

	v.viewport, cmd = v.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return v, tea.Batch(cmds...)
}

func (v *ChatView) View() string {
	var b strings.Builder
	
	status := "Ожидание"
	if v.isStreaming {
		status = "Генерация..."
	}
	
	modeStr := "Спрашивать"
	if v.tools.Mode == tools.ModeReadOnly { modeStr = "Чтение" }
	if v.tools.Mode == tools.ModeAllowSession { modeStr = "Разрешить" }
	
	// Верхний статус бар (модель, сессия, токены)
	header := fmt.Sprintf("[%s] %s | 💰 $%.3f (%d tok) | 🛡️ Права: %s | %s", 
		v.provider.Name, v.session.Model, v.session.Cost, v.session.Tokens, modeStr, status)
		
	b.WriteString(styles.Title.Render(header))
	b.WriteString("\n\n")
	b.WriteString(v.viewport.View())
	b.WriteString("\n\n")
	b.WriteString(v.textarea.View())
	b.WriteString("\n")
	b.WriteString(styles.ListNormal.Render("Вернуться: Esc | Отправить: Ctrl+S | Права: Ctrl+P"))

	return b.String()
}

