package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/wodorfells/vefier-cli/internal/api"
	"github.com/wodorfells/vefier-cli/internal/config"
)

type StreamChunkMsg struct {
	Chunk api.StreamChunk
	ch    <-chan api.StreamChunk
}

type App struct {
	cfg      *config.Config
	styles   map[string]lipgloss.Style
	renderer *glamour.TermRenderer

	textarea textarea.Model
	viewport viewport.Model

	history []api.Message
	current []string
	
	isStreaming bool
	cancel      context.CancelFunc
}

func NewApp(c *config.Config) *App {
	ta := textarea.New()
	ta.Placeholder = "Сообщение..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 20000
	ta.SetHeight(3)
	ta.ShowLineNumbers = false

	vp := viewport.New(100, 20) 

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)

	app := &App{
		cfg:      c,
		styles:   GenerateStyles(c.Theme),
		renderer: r,
		textarea: ta,
		viewport: vp,
		history:  []api.Message{{Role: "system", Content: "Вы — полезный AI-ассистент."}},
	}
	app.refreshViewport()
	return app
}

func (a *App) Init() tea.Cmd {
	return textarea.Blink
}

func listenStream(ch <-chan api.StreamChunk) tea.Cmd {
	return func() tea.Msg {
		chunk, ok := <-ch
		if !ok {
			return StreamChunkMsg{Chunk: api.StreamChunk{Done: true}, ch: ch}
		}
		return StreamChunkMsg{Chunk: chunk, ch: ch}
	}
}

func (a *App) refreshViewport() {
	var b strings.Builder
	provider := a.cfg.GetActive()

	if len(a.history) == 1 {
		welcome := "Привет! Я готов к работе."
		if a.cfg.Language == "en" {
			welcome = "Hello! I am ready to help."
		}
		b.WriteString(a.styles["ai"].Render(provider.Name + ":"))
		b.WriteString("\n" + a.styles["text"].Render(welcome) + "\n\n")
	}

	for _, m := range a.history {
		if m.Role == "system" {
			continue
		}
		if m.Role == "user" {
			b.WriteString(a.styles["user"].Render("User"))
			b.WriteString("\n" + a.styles["text"].Render(m.Content) + "\n\n")
		} else {
			b.WriteString(a.styles["ai"].Render(provider.Name))
			md, _ := a.renderer.Render(m.Content)
			b.WriteString("\n" + md + "\n")
		}
	}

	if a.isStreaming && len(a.current) > 0 {
		b.WriteString(a.styles["ai"].Render(provider.Name))
		md, _ := a.renderer.Render(strings.Join(a.current, ""))
		b.WriteString("\n" + md + "\n")
	}

	a.viewport.SetContent(b.String())
	a.viewport.GotoBottom()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.viewport.Width = msg.Width
		a.viewport.Height = msg.Height - 6
		a.textarea.SetWidth(msg.Width)
		a.refreshViewport()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if a.isStreaming && a.cancel != nil {
				a.cancel()
				a.isStreaming = false
			} else {
				return a, tea.Quit
			}
		case "ctrl+t":
			a.cfg.Theme = NextTheme(a.cfg.Theme)
			a.styles = GenerateStyles(a.cfg.Theme)
			a.cfg.Save()
			a.refreshViewport()
		case "ctrl+l":
			if a.cfg.Language == "ru" {
				a.cfg.Language = "en"
			} else {
				a.cfg.Language = "ru"
			}
			a.cfg.Save()
			a.refreshViewport()
		case "ctrl+o":
			if a.cfg.ActiveProvider == "openai" {
				a.cfg.ActiveProvider = "openrouter"
			} else {
				a.cfg.ActiveProvider = "openai"
			}
			a.cfg.Save()
			a.refreshViewport()
		case "enter":
			if !a.isStreaming && strings.TrimSpace(a.textarea.Value()) != "" {
				text := a.textarea.Value()
				a.textarea.Reset()
				
				a.history = append(a.history, api.Message{Role: "user", Content: text})
				a.current = []string{}
				a.isStreaming = true
				
				a.refreshViewport()

				var ctx context.Context
				ctx, a.cancel = context.WithCancel(context.Background())
				ch := make(chan api.StreamChunk)
				
				go api.StreamChat(ctx, a.cfg.GetActive(), a.history, ch)
				cmds = append(cmds, listenStream(ch))
			}
		}

	case StreamChunkMsg:
		if msg.Chunk.Error != nil {
			a.current = append(a.current, "\n**Ошибка:** "+msg.Chunk.Error.Error())
			a.isStreaming = false
		} else if msg.Chunk.Done {
			a.history = append(a.history, api.Message{Role: "assistant", Content: strings.Join(a.current, "")})
			a.current = []string{}
			a.isStreaming = false
		} else {
			a.current = append(a.current, msg.Chunk.Content)
			cmds = append(cmds, listenStream(msg.ch))
		}
		a.refreshViewport()
	}

	var cmd tea.Cmd
	a.textarea, cmd = a.textarea.Update(msg)
	cmds = append(cmds, cmd)

	a.viewport, cmd = a.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *App) View() string {
	provider := a.cfg.GetActive()
	
	headerText := fmt.Sprintf(" VeFier CLI | Провайдер: %s | Язык: %s | Тема: %s ",
		provider.Name, a.cfg.Language, a.cfg.Theme)
		
	header := a.styles["title"].Render(headerText)

	helpStr := "Enter: Отправить | Ctrl+C: Отмена/Выход | Ctrl+T: Тема | Ctrl+L: Язык | Ctrl+O: Провайдер"
	if a.cfg.Language == "en" {
		helpStr = "Enter: Send | Ctrl+C: Cancel/Quit | Ctrl+T: Theme | Ctrl+L: Language | Ctrl+O: Provider"
	}
	help := a.styles["help"].Render(helpStr)

	return fmt.Sprintf("%s\n\n%s\n%s\n%s", header, a.viewport.View(), a.textarea.View(), help)
}