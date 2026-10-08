package api

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/wodorfells/vefier-cli/internal/config"
)

type StreamChunk struct {
	Content string
	Error   error
	Done    bool
}

type Message struct {
	Role    string
	Content string
}

func StreamChat(ctx context.Context, p config.Provider, history []Message, ch chan<- StreamChunk) {
	defer close(ch)
	
	if p.Key == "" {
		ch <- StreamChunk{Error: errors.New("API-ключ не настроен. Нажмите Ctrl+O для настроек.")}
		return
	}

	cfg := openai.DefaultConfig(p.Key)
	cfg.BaseURL = p.BaseURL
	
	// Попытка исправить совместимость с OpenRouter/другими (они могут падать на кастомных заголовках openai)
	if strings.Contains(p.BaseURL, "openrouter") {
		cfg.BaseURL = "https://openrouter.ai/api/v1"
	}

	client := openai.NewClientWithConfig(cfg)

	var msgs []openai.ChatCompletionMessage
	for _, m := range history {
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	req := openai.ChatCompletionRequest{
		Model:    p.Model,
		Messages: msgs,
		Stream:   true,
	}

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		ch <- StreamChunk{Error: err}
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			ch <- StreamChunk{Done: true}
			return
		}
		if err != nil {
			ch <- StreamChunk{Error: err}
			return
		}
		if len(response.Choices) > 0 {
			ch <- StreamChunk{Content: response.Choices[0].Delta.Content}
		}
	}
}